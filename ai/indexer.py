"""Embedding pipeline: read articles from yarr DB, chunk, embed, store in ChromaDB."""

import asyncio
import hashlib
import logging
from datetime import datetime

import numpy as np

from .article_extractor import html_to_text
from .chunker import chunk_text
from .store import add_chunks, article_exists, update_article_tags
from .yarr_db import get_items_by_ids, list_items, iter_all_items, get_feed_folder_map, open_db

log = logging.getLogger(__name__)

MIN_ARTICLE_LENGTH = 100
MAX_CHUNK_CHARS = 6000  # ~4000 tokens, safe for nomic-embed-text (8192 token window)
REINDEX_BATCH_SIZE = 200  # articles per embed+upsert batch


def generate_article_id(url: str) -> str:
    return hashlib.md5(url.encode()).hexdigest()


def _get_existing_ids(collection) -> tuple[set[str], set[str]]:
    """Pre-fetch all indexed URLs and content hashes in a single ChromaDB call."""
    data = collection.get(include=["metadatas"])
    urls: set[str] = set()
    hashes: set[str] = set()
    for m in data["metadatas"]:
        if url := m.get("url", ""):
            urls.add(url)
        if h := m.get("content_hash", ""):
            hashes.add(h)
    return urls, hashes


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


def index_items(config, collection, item_ids: list[int]) -> tuple[int, list[str]]:
    """Index specific items by ID (webhook handler).

    Returns (count_indexed, list_of_new_urls).
    """
    conn = open_db(config.yarr_db)
    try:
        items = get_items_by_ids(conn, item_ids)
    finally:
        conn.close()

    if not items:
        return 0, []

    existing_urls, existing_hashes = _get_existing_ids(collection)

    count = 0
    new_urls: list[str] = []
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
        new_urls.append(result["url"])
        count += 1

    log.info("Indexed %d/%d items", count, len(items))
    return count, new_urls


def _article_mean_embedding(collection, url: str) -> "np.ndarray | None":
    """Return the mean embedding vector for all chunks of a given article URL.

    Returns None if no chunks are found or embeddings are unavailable.
    """
    try:
        data = collection.get(
            where={"url": {"$eq": url}},
            include=["embeddings"],
        )
    except Exception as e:
        log.warning("ChromaDB get embeddings failed for %s: %s", url, e)
        return None

    embeddings = data.get("embeddings")
    if embeddings is None:
        return None

    try:
        vecs = np.array(embeddings, dtype=np.float64)
        if vecs.size == 0:
            return None
        return vecs.mean(axis=0)
    except Exception as e:
        log.warning("Failed to compute mean embedding for %s: %s", url, e)
        return None


def assign_articles_to_clusters(
    collection,
    urls: list[str],
    cluster_ids: list[int],
    cluster_labels: list[str],
    centroid_matrix: "np.ndarray",
    min_similarity: float = 0.55,
) -> list[dict]:
    """Assign a list of article URLs to the nearest existing cluster.

    For each URL:
    - Fetches the mean embedding from ChromaDB.
    - Computes cosine similarity against every centroid (single dot product).
    - Assigns the best-matching label if similarity >= min_similarity.
    - Updates the ChromaDB tags metadata for the article.

    Returns a list of {"url": ..., "tag": ...} for all successfully assigned articles.
    """
    if not urls or centroid_matrix is None or len(centroid_matrix) == 0:
        return []

    # Normalise centroids once for cosine similarity via dot product
    norms = np.linalg.norm(centroid_matrix, axis=1, keepdims=True)
    norms = np.where(norms == 0, 1.0, norms)
    normed_centroids = centroid_matrix / norms

    assigned: list[dict] = []
    skipped_threshold = 0
    skipped_no_emb = 0

    for url in urls:
        vec = _article_mean_embedding(collection, url)
        if vec is None:
            skipped_no_emb += 1
            continue

        vec_norm = np.linalg.norm(vec)
        if vec_norm == 0:
            skipped_no_emb += 1
            continue
        normed_vec = vec / vec_norm

        sims = normed_centroids.dot(normed_vec)
        best_idx = int(np.argmax(sims))
        best_sim = float(sims[best_idx])

        if best_sim < min_similarity:
            skipped_threshold += 1
            log.debug("Incremental assign: %s skipped (best_sim=%.3f < %.2f)", url, best_sim, min_similarity)
            continue

        label = cluster_labels[best_idx]
        update_article_tags(collection, url, label)
        assigned.append({"url": url, "tag": label})

    log.info(
        "Incremental assign: %d/%d new articles tagged (threshold=%.2f), "
        "skipped %d (below threshold), %d (no embedding)",
        len(assigned), len(urls), min_similarity, skipped_threshold, skipped_no_emb,
    )
    return assigned


def index_and_assign_items(config, collection, item_ids: list[int]) -> tuple[int, list[str]]:
    """Index items and immediately assign them to existing clusters (if centroids exist).

    This is the main entry point for the webhook handlers. It:
    1. Embeds new articles via index_items().
    2. Loads persisted centroids from the Go DB.
    3. Assigns each new article to its nearest cluster.
    4. POSTs the {url, tag} pairs to /api/ai/articles/append (append-only).

    Returns (count_indexed, new_urls) — same shape as index_items().
    """
    from .cluster import fetch_previous_centroids, load_centroid_matrix, post_article_tags_append

    count, new_urls = index_items(config, collection, item_ids)

    if count == 0 or not new_urls:
        return count, new_urls

    if not config.yarr_api_url:
        log.info("Incremental assign: YARR_API_URL not set — skipping topic assignment")
        return count, new_urls

    previous_centroids = fetch_previous_centroids(config.yarr_api_url)
    if not previous_centroids:
        log.info("Incremental assign: no prior centroids — skipping topic assignment")
        return count, new_urls

    centroid_data = load_centroid_matrix(previous_centroids)
    if centroid_data is None:
        log.warning("Incremental assign: centroid matrix could not be loaded — skipping")
        return count, new_urls

    cluster_ids, cluster_labels, centroid_matrix = centroid_data

    article_tags = assign_articles_to_clusters(
        collection,
        new_urls,
        cluster_ids,
        cluster_labels,
        centroid_matrix,
        min_similarity=config.assign_min_similarity,
    )

    if article_tags:
        post_article_tags_append(config.yarr_api_url, article_tags)

    return count, new_urls


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

    # Pre-fetch existing URLs and hashes in a single ChromaDB call
    existing_urls, existing_hashes = _get_existing_ids(collection)
    already = len(existing_urls)
    progress(f"Checking already indexed ({already} found), scanning new articles...")

    # --- Phase 1: stream articles from DB in batches, process into chunks (no API calls) ---
    # Articles are streamed 500 at a time so full HTML content is never all in memory at once.
    pending: list[dict] = []
    skipped = 0
    scan_errors = 0
    scanned = 0

    conn = open_db(config.yarr_db)
    try:
        for item in iter_all_items(conn, batch_size=500):
            scanned += 1
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

            if scanned % 500 == 0:
                progress(f"Scanning: {scanned} articles ({len(pending)} to index, {skipped} skipped)")
    finally:
        conn.close()

    new_total = len(pending)
    log.info("Scan complete: %d new articles to index (%d scanned, %d skipped, %d errors)", new_total, scanned, skipped, scan_errors)
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

        batch_num = batch_start // REINDEX_BATCH_SIZE + 1
        total_batches = (new_total + REINDEX_BATCH_SIZE - 1) // REINDEX_BATCH_SIZE
        progress(f"Embedding: {batch_start}/{new_total} articles (batch {batch_num}/{total_batches})")

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
        log.info("Indexed batch %d/%d (%d articles, %d total)", batch_num, total_batches, len(batch), count)
        progress(f"Indexing: {count}/{new_total} articles (batch {batch_num}/{total_batches})")

    log.info("Reindex complete: %d new articles indexed (%d errors)", count, embed_errors)
    progress(f"Complete: {count} new articles indexed" + (f" ({embed_errors} errors)" if embed_errors else ""))
    return count
