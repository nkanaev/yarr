"""Topic clustering via HDBSCAN with LLM labeling and label stability."""

import json
import logging
import re
from collections import defaultdict
from datetime import datetime
from pathlib import Path

import httpx
import numpy as np

log = logging.getLogger(__name__)


def get_article_embeddings(collection) -> tuple[list[dict], np.ndarray]:
    """Extract per-article embeddings by averaging chunk embeddings."""
    data = collection.get(include=["embeddings", "metadatas"])
    if not data["ids"]:
        return [], np.array([])

    # Group chunks by URL
    articles_map: dict[str, dict] = {}
    for i, meta in enumerate(data["metadatas"]):
        url = meta.get("url", "")
        if url not in articles_map:
            articles_map[url] = {
                "url": url,
                "title": meta.get("title", ""),
                "published": meta.get("published", ""),
                "folder": meta.get("folder", ""),
                "feed_name": meta.get("feed_name", ""),
                "embeddings": [],
            }
        articles_map[url]["embeddings"].append(data["embeddings"][i])

    articles = []
    embeddings = []
    for info in articles_map.values():
        articles.append({k: v for k, v in info.items() if k != "embeddings"})
        embeddings.append(np.mean(info["embeddings"], axis=0))

    return articles, np.array(embeddings)


def cluster_hdbscan(embeddings: np.ndarray, min_cluster_size: int = 10) -> np.ndarray:
    """Cluster with UMAP dimensionality reduction + HDBSCAN."""
    import hdbscan
    import umap

    log.info("UMAP: %d articles, n_components=15, min_cluster_size=%d", len(embeddings), min_cluster_size)

    reducer = umap.UMAP(
        n_components=15, metric="cosine", random_state=42, n_jobs=1
    )
    reduced = reducer.fit_transform(embeddings)

    clusterer = hdbscan.HDBSCAN(
        min_cluster_size=min_cluster_size,
        metric="euclidean",
        core_dist_n_jobs=1,
    )
    return clusterer.fit_predict(reduced)


def compute_centroids(
    embeddings: np.ndarray, labels: np.ndarray
) -> dict[int, np.ndarray]:
    """Compute mean embedding per cluster (excluding noise -1)."""
    centroids = {}
    for label in set(labels):
        if label == -1:
            continue
        mask = labels == label
        centroids[int(label)] = np.mean(embeddings[mask], axis=0)
    return centroids


def get_representative_articles(
    embeddings: np.ndarray, labels: np.ndarray, articles: list[dict], n: int = 10
) -> dict[int, list[dict]]:
    """Find representative articles per cluster.

    Returns a mix of closest-to-centroid + random articles for better topic coverage.
    """
    centroids = compute_centroids(embeddings, labels)
    reps: dict[int, list[dict]] = {}

    for label, centroid in centroids.items():
        mask = labels == label
        indices = np.where(mask)[0]
        cluster_embs = embeddings[indices]

        # Cosine similarity
        norms = np.linalg.norm(cluster_embs, axis=1) * np.linalg.norm(centroid)
        norms = np.where(norms == 0, 1, norms)
        similarities = np.dot(cluster_embs, centroid) / norms
        distances = 1 - similarities

        # Take n_close closest to centroid + n_random random for diversity
        n_close = min(n // 2, len(indices))
        n_random = min(n - n_close, len(indices) - n_close)

        sorted_idx = np.argsort(distances)
        close_idx = set(sorted_idx[:n_close].tolist())

        # Random from remaining
        remaining = [i for i in range(len(indices)) if i not in close_idx]
        if remaining and n_random > 0:
            rng = np.random.default_rng(42)
            random_idx = set(rng.choice(remaining, size=min(n_random, len(remaining)), replace=False).tolist())
        else:
            random_idx = set()

        selected = sorted(close_idx | random_idx)
        reps[label] = [
            {**articles[indices[i]], "distance_to_centroid": float(distances[i])}
            for i in selected
        ]

    return reps


def _call_llm_for_label(
    msgs: list[dict],
    llm_provider,
    model: str,
    ollama_url: str,
) -> str:
    """Call LLM for a single label with rate-aware retry for external APIs.

    For external providers (Gemini/Grok): retry up to 5 times on 429 with 10s+ backoff.
    Does NOT fall back to Ollama — better to get "Unlabeled" than a bad llama3.2 label.
    Only falls back to Ollama if the provider is completely unreachable (connection error).
    """
    if llm_provider:
        # Check if the provider is a FallbackLLM — if so, use the primary directly
        # to avoid falling back to Ollama on 429
        from .providers import FallbackLLM, GeminiLLM, GrokLLM
        actual_provider = llm_provider
        if isinstance(llm_provider, FallbackLLM):
            actual_provider = llm_provider.primary

        # For external providers: retry with long backoff on rate limits
        if isinstance(actual_provider, (GeminiLLM, GrokLLM)):
            import asyncio
            import concurrent.futures

            for attempt in range(5):
                try:
                    with concurrent.futures.ThreadPoolExecutor() as pool:
                        content = pool.submit(
                            lambda: asyncio.run(actual_provider.chat(msgs, stream=False, temperature=0.3))
                        ).result()
                    return content
                except Exception as e:
                    err_str = str(e)
                    if "429" in err_str or "Too Many" in err_str:
                        wait = 10 * (2 ** attempt)  # 10, 20, 40, 80, 160 seconds
                        log.warning("Label API rate limited (attempt %d/5), waiting %ds...", attempt + 1, wait)
                        import time
                        time.sleep(wait)
                        continue
                    else:
                        raise
            raise Exception("Rate limited after 5 retries")

        # For Ollama provider: call directly
        import asyncio
        import concurrent.futures
        with concurrent.futures.ThreadPoolExecutor() as pool:
            content = pool.submit(
                lambda: asyncio.run(actual_provider.chat(msgs, stream=False, temperature=0.3))
            ).result()
        return content

    # No provider — direct Ollama HTTP call
    resp = httpx.post(
        f"{ollama_url.rstrip('/')}/api/chat",
        json={"model": model, "messages": msgs, "stream": False},
        timeout=60.0,
    )
    resp.raise_for_status()
    return resp.json()["message"]["content"]


def label_clusters(
    representatives: dict[int, list[dict]],
    llm_provider=None,
    model: str = "",
    ollama_url: str = "",
    on_progress=None,
) -> dict[int, str]:
    """Generate cluster labels using LLM with rate-aware retry."""
    system_prompt = """You are a topic labeler. Your task is to name topic clusters based on article titles.

CRITICAL RULES:
- Respond with ONLY the topic name
- Use 2-4 words maximum
- No explanations, no reasoning, no preamble
- No quotes, no punctuation at the end
- Just the topic name itself
- The name should describe the MAJORITY of articles, not outliers

Examples:

Articles:
1. Kubernetes 1.30 Released
2. Docker Compose Best Practices
3. CNCF Landscape 2024
4. Helm Chart Debugging Tips
5. Service Mesh Performance
Topic name: Cloud Native Infrastructure

Articles:
1. Episode 45: The Vanishing at Cecil Hotel
2. Episode 46: Cold Case Files Update
3. Episode 47: Unsolved Mysteries Revisited
4. Jim Morrison: The Final Days
5. Episode 48: The Zodiac Killer
Topic name: True Crime Podcast"""

    labels: dict[int, str] = {}
    preamble_patterns = [
        r"^The topic name is[:\s]+",
        r"^Topic name[:\s]+",
        r"^Topic[:\s]+",
        r"^The topic is[:\s]+",
        r"^Name[:\s]+",
        r"^Label[:\s]+",
        r"^Here[:\s]+",
        r"^Based on (the )?articles?,?\s*",
        r"^I would suggest[:\s]+",
        r"^The topic name appears to be[:\s]+",
        r"^This topic (is|appears to be|can be named)[:\s]+",
        r"^A suitable topic name (is|would be)[:\s]+",
    ]

    total = len(representatives)
    for idx, (label_id, articles) in enumerate(representatives.items()):
        # Use up to 10 titles for better topic representation
        titles = "\n".join(
            f"{i + 1}. {a['title']}" for i, a in enumerate(articles[:10])
        )
        user_msg = f"Articles:\n{titles}\nTopic name:"
        msgs = [
            {"role": "system", "content": system_prompt},
            {"role": "user", "content": user_msg},
        ]

        # Rate limiting: 8-second delay between calls (stays well under 15 RPM)
        if idx > 0:
            import time
            time.sleep(8)

        try:
            content = _call_llm_for_label(msgs, llm_provider, model, ollama_url)
        except Exception as e:
            log.warning("Cluster labeling failed for %d: %s", label_id, e)
            labels[label_id] = f"Cluster {label_id}"
            if on_progress:
                on_progress(f"Labeling: {idx + 1}/{total} clusters (failed)")
            continue

        log.info("LLM raw response for cluster %d: %r", label_id, content[:200])

        # Clean LLM output
        content = re.sub(r"<think>.*?</think>", "", content, flags=re.DOTALL).strip()
        content = content.split("\n")[0].strip()
        for pattern in preamble_patterns:
            content = re.sub(pattern, "", content, flags=re.IGNORECASE).strip()
        content = content.strip("\"'""''")
        content = content.rstrip(".,;:!?")

        log.info("Cleaned label for cluster %d: %r (len=%d)", label_id, content, len(content))

        if len(content) < 4 or len(content) > 40:
            content = f"Cluster {label_id}"

        labels[label_id] = content
        if on_progress:
            on_progress(f"Labeling: {idx + 1}/{total} clusters ({content})")

    return labels


def assign_noise_to_nearest(
    embeddings: np.ndarray, labels: np.ndarray, centroids: dict[int, np.ndarray]
) -> np.ndarray:
    """Reassign noise articles (-1) to nearest cluster centroid."""
    new_labels = labels.copy()
    noise_mask = labels == -1
    noise_indices = np.where(noise_mask)[0]

    if len(noise_indices) == 0 or not centroids:
        return new_labels

    centroid_ids = list(centroids.keys())
    centroid_matrix = np.array([centroids[cid] for cid in centroid_ids])

    for idx in noise_indices:
        emb = embeddings[idx]
        norms = np.linalg.norm(centroid_matrix, axis=1) * np.linalg.norm(emb)
        norms = np.where(norms == 0, 1, norms)
        sims = np.dot(centroid_matrix, emb) / norms
        best = centroid_ids[int(np.argmax(sims))]
        new_labels[idx] = best

    return new_labels


def stabilize_labels(
    new_centroids: dict[int, np.ndarray],
    previous_clusters: dict | None,
    threshold: float = 0.85,
) -> dict[int, str | None]:
    """Match new clusters to previous ones by centroid similarity.

    Returns mapping: new_cluster_id -> old_label or None (needs new label).
    """
    if not previous_clusters:
        return {cid: None for cid in new_centroids}

    old_clusters = previous_clusters.get("clusters", [])
    old_centroids = {}
    old_labels = {}
    for c in old_clusters:
        cid = c["id"]
        if c.get("centroid"):
            old_centroids[cid] = np.array(c["centroid"])
            old_labels[cid] = c.get("label", "")

    result: dict[int, str | None] = {}
    for new_id, new_centroid in new_centroids.items():
        best_sim = -1.0
        best_label = None
        for old_id, old_centroid in old_centroids.items():
            norm = np.linalg.norm(new_centroid) * np.linalg.norm(old_centroid)
            if norm == 0:
                continue
            sim = float(np.dot(new_centroid, old_centroid) / norm)
            if sim > best_sim:
                best_sim = sim
                best_label = old_labels.get(old_id)

        result[new_id] = best_label if best_sim >= threshold else None

    return result


def build_cluster_map(
    articles: list[dict],
    labels: np.ndarray,
    embeddings: np.ndarray,
    cluster_names: dict[int, str],
) -> dict:
    """Build JSON-serializable cluster map."""
    centroids = compute_centroids(embeddings, labels)

    clusters = []
    for cid in sorted(set(labels)):
        if cid == -1:
            continue
        mask = labels == cid
        cluster_articles = [articles[i] for i in np.where(mask)[0]]
        clusters.append({
            "id": int(cid),
            "label": cluster_names.get(int(cid), f"Cluster {cid}"),
            "article_count": int(np.sum(mask)),
            "centroid": centroids[int(cid)].tolist() if int(cid) in centroids else [],
            "articles": cluster_articles,
        })

    clusters.sort(key=lambda x: x["article_count"], reverse=True)

    noise_count = int(np.sum(labels == -1))
    return {
        "generated_at": datetime.utcnow().isoformat(),
        "algorithm": "hdbscan",
        "n_clusters": len(clusters),
        "n_articles": len(articles),
        "n_noise": noise_count,
        "clusters": clusters,
    }


def save_cluster_map(cluster_map: dict, output_path: str) -> None:
    Path(output_path).parent.mkdir(parents=True, exist_ok=True)
    with open(output_path, "w") as f:
        json.dump(cluster_map, f, indent=2, default=str)


def load_previous_clusters(path: str) -> dict | None:
    try:
        with open(path) as f:
            return json.load(f)
    except (FileNotFoundError, json.JSONDecodeError):
        return None


def split_large_clusters(
    embeddings: np.ndarray,
    labels: np.ndarray,
    max_cluster_size: int = 60,
    target_subcluster_size: int = 40,
) -> np.ndarray:
    """Split clusters with >max_cluster_size articles using KMeans.

    Large clusters are broken into sub-clusters of ~target_subcluster_size articles
    for more specific topic labeling.
    """
    from sklearn.cluster import KMeans

    new_labels = labels.copy()
    next_label = int(labels.max()) + 1 if len(labels) > 0 else 0

    for cid in sorted(set(labels)):
        if cid == -1:
            continue
        mask = labels == cid
        count = int(np.sum(mask))
        if count <= max_cluster_size:
            continue

        indices = np.where(mask)[0]
        cluster_embs = embeddings[indices]
        k = max(2, count // target_subcluster_size)

        log.info("Splitting cluster %d (%d articles) into %d sub-clusters", cid, count, k)
        kmeans = KMeans(n_clusters=k, random_state=42, n_init=10)
        sub_labels = kmeans.fit_predict(cluster_embs)

        # Reassign: first sub-cluster keeps original ID, rest get new IDs
        for sub_id in range(k):
            sub_mask = sub_labels == sub_id
            sub_indices = indices[sub_mask]
            if sub_id == 0:
                new_labels[sub_indices] = cid
            else:
                new_labels[sub_indices] = next_label
                next_label += 1

    n_before = len(set(labels) - {-1})
    n_after = len(set(new_labels) - {-1})
    if n_after > n_before:
        log.info("Split large clusters: %d -> %d clusters", n_before, n_after)

    return new_labels


def merge_similar_labels(
    cluster_names: dict[int, str],
    llm_provider=None,
    model: str = "",
    ollama_url: str = "",
) -> dict[int, str]:
    """Ask LLM to merge duplicate/overlapping topic labels.

    Returns updated cluster_names with duplicates mapped to canonical labels.
    """
    if len(cluster_names) < 2:
        return cluster_names

    labels_list = "\n".join(
        f"- {name}" for name in sorted(set(cluster_names.values()))
    )

    prompt = f"""Here are topic labels from a news clustering system.
Identify any that clearly refer to the same topic and should be merged.
Return ONLY a JSON object mapping duplicate labels to the canonical (best) label.
Choose the most descriptive label as canonical.
If no merges needed, return {{}}

Labels:
{labels_list}

JSON:"""

    try:
        msgs = [{"role": "user", "content": prompt}]
        content = _call_llm_for_label(msgs, llm_provider, model, ollama_url)

        # Extract JSON from response
        content = re.sub(r"<think>.*?</think>", "", content, flags=re.DOTALL).strip()
        # Find JSON object in response
        match = re.search(r"\{[^{}]*\}", content)
        if not match:
            log.info("No merges suggested by LLM")
            return cluster_names

        merge_map = json.loads(match.group())
        if not merge_map:
            return cluster_names

        log.info("LLM suggested %d label merges: %s", len(merge_map), merge_map)

        # Apply merges
        merged = {}
        for cid, name in cluster_names.items():
            merged[cid] = merge_map.get(name, name)
        return merged

    except Exception as e:
        log.warning("Label merge failed: %s", e)
        return cluster_names


def run_clustering(config, llm_provider=None, embed_provider=None, on_progress=None) -> dict | None:
    """Full clustering pipeline. Returns cluster map or None."""
    from .store import init_store, update_article_tags
    from .providers import OllamaEmbed

    def progress(msg):
        if on_progress:
            on_progress(msg)

    ep = embed_provider or OllamaEmbed(config.embed_model, config.ollama_url)
    collection, _ = init_store(config.chroma_path, ep)

    progress("Extracting article embeddings...")
    articles, embeddings = get_article_embeddings(collection)

    if len(articles) < config.min_cluster_size:
        log.info("Not enough articles for clustering (%d)", len(articles))
        progress(f"Not enough articles ({len(articles)}, need {config.min_cluster_size})")
        return None

    log.info("Clustering %d articles...", len(articles))
    progress(f"Clustering {len(articles)} articles...")
    labels = cluster_hdbscan(embeddings, min_cluster_size=config.min_cluster_size)

    centroids = compute_centroids(embeddings, labels)
    if not centroids:
        log.warning("No clusters found")
        return None

    # Reassign noise
    labels = assign_noise_to_nearest(embeddings, labels, centroids)

    # Split oversized clusters into sub-topics
    progress("Splitting large clusters...")
    labels = split_large_clusters(embeddings, labels, max_cluster_size=60, target_subcluster_size=40)
    centroids = compute_centroids(embeddings, labels)

    n_clusters = len(set(labels) - {-1})
    progress(f"Found {n_clusters} clusters, selecting representatives...")
    representatives = get_representative_articles(embeddings, labels, articles, n=10)

    # Label stability
    output_path = str(Path(config.chroma_path).parent / "clusters.json")
    previous = load_previous_clusters(output_path)
    stable = stabilize_labels(centroids, previous)

    # LLM-label new/changed clusters AND clusters with generic "Cluster N" labels
    to_label = {
        cid: representatives[cid]
        for cid, old in stable.items()
        if (old is None or re.match(r"^Cluster \d+$", old or "")) and cid in representatives
    }
    n_clusters = len(set(labels) - {-1})
    new_labels = {}
    if to_label:
        progress(f"Found {n_clusters} clusters, labeling {len(to_label)}...")
        new_labels = label_clusters(to_label, llm_provider=llm_provider, model=config.label_model, ollama_url=config.ollama_url, on_progress=progress)

    # Merge labels — treat generic "Cluster N" as unlabeled
    cluster_names: dict[int, str] = {}
    for cid in set(labels):
        if cid == -1:
            continue
        cid_int = int(cid)
        old_label = stable.get(cid_int)
        is_generic = old_label and re.match(r"^Cluster \d+$", old_label)
        if old_label and not is_generic:
            cluster_names[cid_int] = old_label
        elif cid_int in new_labels:
            cluster_names[cid_int] = new_labels[cid_int]
        else:
            if cid_int in representatives:
                log.info("Force-labeling cluster %d (was generic)", cid_int)
                forced = label_clusters({cid_int: representatives[cid_int]}, llm_provider=llm_provider, model=config.label_model, ollama_url=config.ollama_url)
                cluster_names[cid_int] = forced.get(cid_int, f"Cluster {cid_int}")
            else:
                cluster_names[cid_int] = f"Cluster {cid_int}"

    # Merge duplicate labels via LLM
    progress("Merging similar topic labels...")
    import time
    time.sleep(8)  # Rate limit buffer before merge call
    cluster_names = merge_similar_labels(cluster_names, llm_provider=llm_provider, model=config.label_model, ollama_url=config.ollama_url)

    cluster_map = build_cluster_map(articles, labels, embeddings, cluster_names)
    save_cluster_map(cluster_map, output_path)

    # Update ChromaDB tags
    progress(f"Updating tags for {len(articles)} articles...")
    for cid_int, name in cluster_names.items():
        mask = labels == cid_int
        for i in np.where(mask)[0]:
            url = articles[i]["url"]
            update_article_tags(collection, url, name)

    log.info("Clustering complete: %d clusters", len(cluster_names))
    progress(f"Complete: {len(cluster_names)} clusters found")
    return cluster_map
