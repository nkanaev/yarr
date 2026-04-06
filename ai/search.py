"""Hybrid search: BM25 keyword + vector + Reciprocal Rank Fusion."""

import logging

import chromadb
from rank_bm25 import BM25Okapi

log = logging.getLogger(__name__)


def _to_float(value: object) -> float:
    try:
        return float(value)
    except (TypeError, ValueError):
        return 0.0


def build_bm25_index(
    collection: chromadb.Collection,
) -> tuple[BM25Okapi | None, list[dict]]:
    """Build BM25 index from all documents in ChromaDB."""
    data = collection.get(include=["documents", "metadatas"])
    if not data["documents"]:
        return None, []

    corpus = [doc.lower().split() for doc in data["documents"]]
    doc_metas = []
    for i, doc in enumerate(data["documents"]):
        meta = data["metadatas"][i] if data["metadatas"] else {}
        doc_metas.append({
            "text": doc,
            "url": meta.get("url", ""),
            "title": meta.get("title", ""),
            "published": meta.get("published", ""),
            "published_ts": _to_float(meta.get("published_ts", 0)),
            "folder": meta.get("folder", ""),
            "feed_name": meta.get("feed_name", ""),
            "chunk_index": meta.get("chunk_index", "0"),
            "tags": meta.get("tags", ""),
        })

    try:
        index = BM25Okapi(corpus)
    except Exception as e:
        log.error("Failed to build BM25 index: %s", e)
        return None, []

    return index, doc_metas


def bm25_search(
    bm25_index: BM25Okapi,
    doc_metas: list[dict],
    query_text: str,
    n_results: int = 10,
    folder: str | None = None,
    tag: str | None = None,
    since_ts: float | None = None,
) -> list[dict]:
    """Search using BM25 with optional filters."""
    tokens = query_text.lower().split()
    if not tokens:
        return []

    scores = bm25_index.get_scores(tokens)
    # Over-fetch 3x to compensate for post-filtering
    top_n = min(n_results * 3, len(scores))
    top_indices = sorted(range(len(scores)), key=lambda i: scores[i], reverse=True)[:top_n]

    results = []
    for idx in top_indices:
        if scores[idx] <= 0:
            continue

        meta = doc_metas[idx]

        # Apply filters
        if folder and meta.get("folder") != folder:
            continue
        if tag and tag not in meta.get("tags", ""):
            continue
        if since_ts is not None and meta.get("published_ts", 0) < since_ts:
            continue

        results.append({**meta, "bm25_score": float(scores[idx])})

        if len(results) >= n_results:
            break

    return results


def reciprocal_rank_fusion(
    vector_results: list[dict],
    bm25_results: list[dict],
    k: int = 60,
) -> list[dict]:
    """Merge results from vector and BM25 search using RRF."""
    scored: dict[tuple, dict] = {}

    for rank, r in enumerate(vector_results):
        key = (r.get("url", ""), r.get("chunk_index", "0"))
        if key not in scored:
            scored[key] = {**r, "rrf_score": 0.0}
        scored[key]["rrf_score"] += 1.0 / (k + rank + 1)

    for rank, r in enumerate(bm25_results):
        key = (r.get("url", ""), r.get("chunk_index", "0"))
        if key not in scored:
            scored[key] = {**r, "rrf_score": 0.0}
        scored[key]["rrf_score"] += 1.0 / (k + rank + 1)

    return sorted(scored.values(), key=lambda x: x["rrf_score"], reverse=True)
