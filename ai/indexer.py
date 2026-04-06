"""Embedding pipeline: read articles from yarr DB, chunk, embed, store in ChromaDB."""

import asyncio
import hashlib
import logging

from .article_extractor import html_to_text
from .chunker import chunk_text
from .store import add_chunks, article_exists
from .yarr_db import get_items_by_ids, list_items, list_all_items, get_feed_folder_map, open_db

log = logging.getLogger(__name__)

MIN_ARTICLE_LENGTH = 100
MAX_CHUNK_CHARS = 6000  # ~4000 tokens, safe for nomic-embed-text (8192 token window)


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


def reindex_all(config, collection, on_progress=None) -> int:
    """Full reindex: fetch all items from yarr DB, index new ones."""
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
    progress(f"Found {total} articles ({already} already indexed), processing...")

    count = 0
    errors = 0
    skipped = 0
    processed = 0
    for item in all_items:
        processed += 1
        try:
            result = _process_item(item, config, existing_urls, existing_hashes)
            if result is None:
                skipped += 1
                if processed % 100 == 0:
                    progress(f"Scanning: {processed}/{total} articles ({count} new, {skipped} skipped)")
                continue

            # Filter out empty chunks
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
            errors += 1
            log.warning("Failed to index item %s: %s", item.get("link", "?"), e)
            continue
        existing_urls.add(result["url"])
        existing_hashes.add(result["content_hash"])
        count += 1

        if count % 5 == 0:
            log.info("Indexed %d articles so far...", count)
            progress(f"Indexing: {count} new ({processed}/{total} scanned)")

    log.info("Reindex complete: %d new articles indexed (%d errors)", count, errors)
    progress(f"Complete: {count} new articles indexed ({errors} errors)")
    return count
