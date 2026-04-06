"""ChromaDB wrapper for article chunk storage and retrieval."""

import logging
from datetime import datetime
from pathlib import Path

import chromadb

from .providers import EmbedProvider

log = logging.getLogger(__name__)


class ChromaEmbedAdapter(chromadb.EmbeddingFunction):
    """Wraps an EmbedProvider for ChromaDB's EmbeddingFunction interface."""

    def __init__(self, provider: EmbedProvider):
        self.provider = provider

    def __call__(self, input: list[str]) -> list[list[float]]:
        return self.provider.embed(input)


def init_store(chroma_path: str, embed_provider: EmbedProvider) -> tuple[chromadb.Collection, EmbedProvider]:
    """Initialize ChromaDB with an embedding provider.

    Returns (collection, embed_provider).
    """
    Path(chroma_path).mkdir(parents=True, exist_ok=True)
    client = chromadb.PersistentClient(path=chroma_path)
    adapter = ChromaEmbedAdapter(embed_provider)
    collection = client.get_or_create_collection(
        name="articles",
        embedding_function=adapter,
        metadata={"hnsw:space": "cosine"},
    )
    return collection, embed_provider


def article_exists(collection: chromadb.Collection, url: str) -> bool:
    result = collection.get(where={"url": {"$eq": url}})
    return len(result["ids"]) > 0


def add_chunks(
    collection: chromadb.Collection,
    chunks: list[str],
    article_id: str,
    title: str,
    url: str,
    published: str | None = None,
    folder: str | None = None,
    feed_name: str | None = None,
    tags: str | None = None,
    content_hash: str | None = None,
) -> None:
    """Upsert article chunks with metadata."""
    ids = [f"{article_id}_chunk_{i}" for i in range(len(chunks))]

    published_ts = 0.0
    if published:
        try:
            published_ts = datetime.fromisoformat(published).timestamp()
        except (ValueError, TypeError):
            pass

    metadatas = [
        {
            "url": url,
            "title": title,
            "chunk_index": str(i),
            "published": published or "",
            "published_ts": published_ts,
            "folder": folder or "uncategorized",
            "feed_name": feed_name or "unknown",
            "tags": tags or "",
            "content_hash": content_hash or "",
        }
        for i in range(len(chunks))
    ]

    collection.upsert(ids=ids, documents=chunks, metadatas=metadatas)


def query(
    collection: chromadb.Collection,
    question: str,
    n_results: int = 5,
    where: dict | None = None,
    distance_threshold: float | None = None,
    since_ts: float | None = None,
    embed_fn: EmbedProvider | None = None,
) -> list[dict]:
    """Query ChromaDB with optional filtering and distance threshold."""
    filters = []
    if where:
        filters.append(where)
    if since_ts is not None:
        filters.append({"published_ts": {"$gte": since_ts}})

    combined_where = None
    if len(filters) == 1:
        combined_where = filters[0]
    elif len(filters) > 1:
        combined_where = {"$and": filters}

    try:
        if embed_fn is not None:
            embeddings = embed_fn.embed([question])
            result = collection.query(
                query_embeddings=embeddings,
                n_results=n_results,
                include=["documents", "metadatas", "distances"],
                where=combined_where,
            )
        else:
            result = collection.query(
                query_texts=[question],
                n_results=n_results,
                include=["documents", "metadatas", "distances"],
                where=combined_where,
            )
    except Exception as e:
        log.error("ChromaDB query failed: %s", e)
        return []

    if not result["ids"] or not result["ids"][0]:
        return []

    results = []
    best = None

    for i, doc_id in enumerate(result["ids"][0]):
        dist = result["distances"][0][i] if result["distances"] else 0.0
        meta = result["metadatas"][0][i] if result["metadatas"] else {}
        text = result["documents"][0][i] if result["documents"] else ""

        entry = {
            "text": text,
            "url": meta.get("url", ""),
            "title": meta.get("title", ""),
            "published": meta.get("published", ""),
            "folder": meta.get("folder", ""),
            "feed_name": meta.get("feed_name", ""),
            "chunk_index": meta.get("chunk_index", "0"),
            "distance": dist,
            "below_threshold": True,
        }

        if best is None or dist < best["distance"]:
            best = entry

        if distance_threshold is not None and dist > distance_threshold:
            continue

        results.append(entry)

    if not results and best:
        best["below_threshold"] = False
        results.append(best)

    return results


def get_collection_stats(collection: chromadb.Collection) -> dict:
    return {"total_documents": collection.count()}


def list_articles(
    collection: chromadb.Collection,
    folder: str | None = None,
    tag: str | None = None,
) -> list[dict]:
    """List unique articles, optionally filtered."""
    filters = []
    if folder:
        filters.append({"folder": {"$eq": folder}})
    if tag:
        filters.append({"tags": {"$eq": tag}})

    where = None
    if len(filters) == 1:
        where = filters[0]
    elif len(filters) > 1:
        where = {"$and": filters}

    data = collection.get(include=["metadatas"], where=where)

    seen: dict[str, dict] = {}
    for meta in data["metadatas"]:
        url = meta.get("url", "")
        if url not in seen:
            seen[url] = {
                "url": url,
                "title": meta.get("title", ""),
                "published": meta.get("published", ""),
                "published_ts": float(meta.get("published_ts", 0)),
                "folder": meta.get("folder", ""),
                "feed_name": meta.get("feed_name", ""),
                "tags": meta.get("tags", ""),
                "chunk_count": 1,
            }
        else:
            seen[url]["chunk_count"] += 1

    return sorted(seen.values(), key=lambda x: x["published_ts"], reverse=True)


def get_article_chunks(collection: chromadb.Collection, url: str) -> list[dict]:
    """Get all chunks for a specific article URL."""
    data = collection.get(
        where={"url": {"$eq": url}},
        include=["documents", "metadatas"],
    )
    chunks = []
    for i, doc in enumerate(data["documents"]):
        meta = data["metadatas"][i]
        chunks.append({
            "chunk_index": int(meta.get("chunk_index", 0)),
            "text": doc,
        })
    return sorted(chunks, key=lambda x: x["chunk_index"])


def list_topics(collection: chromadb.Collection) -> list[dict]:
    """Group articles by folder with counts."""
    data = collection.get(include=["metadatas"])

    folders: dict[str, set] = {}
    for meta in data["metadatas"]:
        folder = meta.get("folder", "uncategorized")
        url = meta.get("url", "")
        if folder not in folders:
            folders[folder] = set()
        folders[folder].add(url)

    topics = [
        {"folder": f, "article_count": len(urls)}
        for f, urls in folders.items()
    ]
    return sorted(topics, key=lambda x: x["article_count"], reverse=True)


def list_tags(collection: chromadb.Collection) -> list[dict]:
    """List all tags with article counts."""
    data = collection.get(include=["metadatas"])

    tag_articles: dict[str, set] = {}
    for meta in data["metadatas"]:
        tags_str = meta.get("tags", "")
        url = meta.get("url", "")
        for tag in tags_str.split(","):
            tag = tag.strip()
            if tag:
                if tag not in tag_articles:
                    tag_articles[tag] = set()
                tag_articles[tag].add(url)

    return sorted(
        [{"tag": t, "article_count": len(urls)} for t, urls in tag_articles.items()],
        key=lambda x: x["article_count"],
        reverse=True,
    )


def update_article_tags(collection: chromadb.Collection, url: str, tags: str) -> None:
    """Update tags for all chunks of an article."""
    data = collection.get(where={"url": {"$eq": url}}, include=["metadatas"])
    if not data["ids"]:
        return
    metadatas = [dict(m, tags=tags) for m in data["metadatas"]]
    collection.update(ids=data["ids"], metadatas=metadatas)
