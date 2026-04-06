"""Article deduplication via embedding cosine similarity."""

import logging

import numpy as np

from .cluster import get_article_embeddings

log = logging.getLogger(__name__)


def find_dedup_groups(
    collection,
    threshold: float = 0.92,
    since_ts: float | None = None,
) -> list[dict]:
    """Find groups of near-duplicate articles by embedding similarity.

    Args:
        collection: ChromaDB collection
        threshold: cosine similarity threshold (0-1). Articles above this are dupes.
        since_ts: only consider articles published after this timestamp.

    Returns:
        List of dedup groups, each with representative + grouped articles.
    """
    articles, embeddings = get_article_embeddings(collection)
    if len(articles) < 2:
        return []

    # Filter by time if requested
    if since_ts is not None:
        from datetime import datetime

        keep = []
        for i, a in enumerate(articles):
            pub = a.get("published", "")
            try:
                ts = datetime.fromisoformat(pub).timestamp()
            except (ValueError, TypeError):
                ts = 0.0
            if ts >= since_ts:
                keep.append(i)

        if len(keep) < 2:
            return []

        articles = [articles[i] for i in keep]
        embeddings = embeddings[keep]

    # Compute pairwise cosine similarity
    norms = np.linalg.norm(embeddings, axis=1, keepdims=True)
    norms = np.where(norms == 0, 1, norms)
    normalized = embeddings / norms
    sim_matrix = np.dot(normalized, normalized.T)

    # Find connected components above threshold (greedy union-find)
    n = len(articles)
    parent = list(range(n))

    def find(x):
        while parent[x] != x:
            parent[x] = parent[parent[x]]
            x = parent[x]
        return x

    def union(a, b):
        ra, rb = find(a), find(b)
        if ra != rb:
            parent[ra] = rb

    for i in range(n):
        for j in range(i + 1, n):
            if sim_matrix[i, j] >= threshold:
                union(i, j)

    # Group by root
    groups_map: dict[int, list[int]] = {}
    for i in range(n):
        root = find(i)
        if root not in groups_map:
            groups_map[root] = []
        groups_map[root].append(i)

    # Only keep groups with 2+ articles
    dedup_groups = []
    for indices in groups_map.values():
        if len(indices) < 2:
            continue

        group_articles = [articles[i] for i in indices]
        # Representative = first published (earliest)
        group_articles.sort(key=lambda a: a.get("published", ""))
        representative = group_articles[0]

        dedup_groups.append({
            "representative": representative,
            "articles": group_articles,
            "source_count": len(group_articles),
        })

    dedup_groups.sort(key=lambda g: g["source_count"], reverse=True)
    return dedup_groups
