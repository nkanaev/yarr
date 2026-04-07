"""Embedding pipeline: read articles from yarr DB, chunk, embed, store in ChromaDB."""

import asyncio
import hashlib
import logging
from datetime import datetime

from .article_extractor import html_to_text
from .chunker import chunk_text
from .store import add_chunks, article_exists
from .yarr_db import get_items_by_ids, list_items, list_all_items, get_feed_folder_map, open_db

log = logging.getLogger(__name__)

MIN_ARTICLE_LENGTH = 100
MAX_CHUNK_CHARS = 6000  # ~4000 tokens, safe for nomic-embed-text (8192 token window)
REINDEX_BATCH_SIZE = 200  # articles per embed+upsert batch


def generate_article_id(url: str) -> str:
    return hashlib.md5(url.encode()).hexdigest()


def _get_existing_urls(collection) -> set[str]:
    """Pre-fetch all indexed URLs for O(1) dedup."""
    data = collection.get(include=["metadatas"])
    return {m.get("url", "") for m in data["metadatas"]}


def _get_existing_hashes(collection) -> set[str]:
    """Pre-fetch content hashes for cross-URL dedup."""
    data = collection.get(include=["metadatas"])
    return {m.get("content_hash", "") for m in data["metadatas"] if m.get("content_hash")}


def _process_item(item: dict, config, existing_urls: set, existing_hashes: set) -> dict | None:
    """Process a single item: extract text, chunk, return data for storage.

    Returns dict with chunks + metadata, or None if skipped.
    """
    url = item.get("link", "")
    if not url:
        return None
    if url in existing_urls:
        return None

    title = item.get("title", "Untitled")
    html = item.get("content", "")
    if not html:
        return None

    # Extract text
    text = html_to_text(html)
    if not text or len(text) < MIN_ARTICLE_LENGTH:
        return None

    # Content hash dedup
    content_hash = hashlib.sha256(text.encode()).hexdigest()
    if content_hash in existing_hashes:
        return None

    # Chunk with context enrichment
    folder = item.get("folder_name", "uncategorized")
    feed_name = item.get("feed_title", "unknown")

    raw_chunks = chunk_text(text, config.chunk_size, config.chunk_overlap)
    if not raw_chunks:
        return None

    # Prepend context to each chunk, truncate to max chars
    enriched_chunks = [
        f"[{title}]\n[Source: {feed_name} | Folder: {folder}]\n\n{chunk}"[:MAX_CHUNK_CHARS]
        for chunk in raw_chunks
    ]

    return {
        "chunks": enriched_chunks,
        "article_id": generate_article_id(url),
        "title": title,
        "url": url,
        "published": item.get("date", ""),
        "folder": folder,
        "feed_name": feed_name,
        "content_hash": content_hash,
    }


def index_items(config, collection, item_ids: list[int]) -> int:
    """Index specific items by ID (webhook handler)."""
    conn = open_db(config.yarr_db)
    try:
        items = get_items_by_ids(conn, item_ids)
    finally:
        conn.close()

    if not items:
        return 0

    existing_urls = _get_existing_urls(collection)
    existing_hashes = _get_existing_hashes(collection)

    count = 0
    for item in items:
        try:
            result = _process_item(item, config, existing_urls, existing_hashes)
            if result is None:
                continue

            result["chunks"] = [c for c in result["chunks"] if c and c.strip()]
            if not result["chunks"]:
                continue

            add_chunks(
                collection,
                chunks=result["chunks"],
                article_id=result["article_id"],
                title=result["title"],
                url=result["url"],
                published=result["published"],
                folder=result["folder"],
                feed_name=result["feed_name"],
                content_hash=result["content_hash"],
            )
        except Exception as e:
            log.warning("Failed to index item %s: %s", item.get("link", "?"), e)
            continue
        existing_urls.add(result["url"])
        existing_hashes.add(result["content_hash"])
        count += 1

    log.info("Indexed %d/%d items", count, len(items))
    return count


def reindex_all(config, collection, embed_provider=None, on_progress=None) -> int:
    """Full reindex: fetch all items from yarr DB, index new ones.

    Uses batched embed+upsert to minimize Gemini API calls:
    - Phase 1: scan all articles, extract text + chunks (CPU only, no API calls)
    - Phase 2: embed REINDEX_BATCH_SIZE articles at once, upsert to ChromaDB in bulk

    embed_provider: if provided, used directly for batched embedding (bypasses
    ChromaDB's per-document embedding adapter for a ~100x speedup).
    Falls back to the collection's adapter if None.
    """
    def progress(msg):
        if on_progress:
            on_progress(msg)

    conn = open_db(config.yarr_db)
    try:
        all_items = list_all_items(conn)
    finally:
        conn.close()

    total = len(all_items)
    log.info("Found %d total items in yarr DB", total)
    progress(f"Found {total} articles, checking already indexed...")

    existing_urls = _get_existing_urls(collection)
    existing_hashes = _get_existing_hashes(collection)
    already = len(existing_urls)
    progress(f"{total} articles total, {already} already indexed — scanning new ones...")

    # --- Phase 1: scan all articles, process into chunks (no API calls) ---
    pending: list[dict] = []
    skipped = 0
    scan_errors = 0

    for i, item in enumerate(all_items):
        try:
            result = _process_item(item, config, existing_urls, existing_hashes)
            if result is None:
                skipped += 1
            else:
                result["chunks"] = [c for c in result["chunks"] if c and c.strip()]
                if result["chunks"]:
                    pending.append(result)
                    existing_urls.add(result["url"])
                    existing_hashes.add(result["content_hash"])
        except Exception as e:
            scan_errors += 1
            log.warning("Failed to process item %s: %s", item.get("link", "?"), e)

        if (i + 1) % 500 == 0:
            progress(f"Scanning: {i + 1}/{total} articles ({len(pending)} to index, {skipped} skipped)")

    new_total = len(pending)
    log.info("Scan complete: %d new articles to index (%d skipped, %d errors)", new_total, skipped, scan_errors)
    progress(f"Scan complete: {new_total} new articles to embed and index...")

    if not pending:
        progress(f"Nothing new to index ({already} already indexed)")
        return 0

    # --- Phase 2: batch embed + upsert ---
    count = 0
    embed_errors = 0

    for batch_start in range(0, new_total, REINDEX_BATCH_SIZE):
        batch = pending[batch_start : batch_start + REINDEX_BATCH_SIZE]

        # Flatten all chunks from all articles in this batch
        all_chunks: list[str] = []
        offsets: list[tuple[int, int]] = []  # (start_idx, end_idx) per article
        for article in batch:
            start = len(all_chunks)
            all_chunks.extend(article["chunks"])
            offsets.append((start, len(all_chunks)))

        try:
            if embed_provider is not None:
                # Direct batched embed — one call for all chunks in the batch
                all_vectors = embed_provider.embed(all_chunks)
            else:
                # Fallback: let ChromaDB embed (slower, one article at a time via adapter)
                all_vectors = None
        except Exception as e:
            embed_errors += len(batch)
            log.error("Batch embed failed (articles %d-%d): %s", batch_start, batch_start + len(batch), e)
            progress(f"Embed error on batch {batch_start // REINDEX_BATCH_SIZE + 1}, skipping {len(batch)} articles")
            continue

        # Build flat lists for a single bulk upsert
        ids: list[str] = []
        documents: list[str] = []
        pre_embedded: list[list[float]] = []
        metadatas: list[dict] = []

        for article, (start, end) in zip(batch, offsets):
            article_id = article["article_id"]
            published_ts = 0.0
            if article["published"]:
                try:
                    published_ts = datetime.fromisoformat(article["published"]).timestamp()
                except (ValueError, TypeError):
                    pass

            for chunk_i, chunk in enumerate(article["chunks"]):
                ids.append(f"{article_id}_chunk_{chunk_i}")
                documents.append(chunk)
                if all_vectors is not None:
                    pre_embedded.append(all_vectors[start + chunk_i])
                metadatas.append({
                    "url": article["url"],
                    "title": article["title"],
                    "chunk_index": str(chunk_i),
                    "published": article["published"] or "",
                    "published_ts": published_ts,
                    "folder": article["folder"] or "uncategorized",
                    "feed_name": article["feed_name"] or "unknown",
                    "tags": "",
                    "content_hash": article["content_hash"],
                })

        try:
            if pre_embedded:
                collection.upsert(ids=ids, documents=documents, embeddings=pre_embedded, metadatas=metadatas)
            else:
                collection.upsert(ids=ids, documents=documents, metadatas=metadatas)
        except Exception as e:
            embed_errors += len(batch)
            log.error("ChromaDB upsert failed for batch %d: %s", batch_start // REINDEX_BATCH_SIZE + 1, e)
            continue

        for article in batch:
            existing_urls.add(article["url"])
            existing_hashes.add(article["content_hash"])

        count += len(batch)
        batch_num = batch_start // REINDEX_BATCH_SIZE + 1
        total_batches = (new_total + REINDEX_BATCH_SIZE - 1) // REINDEX_BATCH_SIZE
        log.info("Indexed batch %d/%d (%d articles, %d total)", batch_num, total_batches, len(batch), count)
        progress(f"Indexing: {count}/{new_total} articles (batch {batch_num}/{total_batches})")

    log.info("Reindex complete: %d new articles indexed (%d errors)", count, embed_errors)
    progress(f"Complete: {count} new articles indexed" + (f" ({embed_errors} errors)" if embed_errors else ""))
    return count
