"""FastAPI routes for the AI service."""

import json
import logging
import re
from datetime import datetime, timedelta
from fastapi import APIRouter, BackgroundTasks, Request
from fastapi.responses import JSONResponse
from sse_starlette.sse import EventSourceResponse

from .briefing import generate_briefing
from .cluster import run_clustering
from .dedup import find_dedup_groups
from .indexer import index_items, index_and_assign_items, reindex_all
from .search import build_bm25_index

log = logging.getLogger(__name__)
router = APIRouter()


def parse_since(since: str | None) -> float | None:
    """Parse a since parameter into a Unix timestamp.

    Supports: Nh (hours), Nd (days), Nw (weeks), Nm (months), ISO date.
    """
    if not since:
        return None
    since = since.strip()
    m = re.match(r"^(\d+)([hdwm])$", since)
    if m:
        n = int(m.group(1))
        unit = m.group(2)
        delta = {
            "h": timedelta(hours=n),
            "d": timedelta(days=n),
            "w": timedelta(weeks=n),
            "m": timedelta(days=n * 30),
        }[unit]
        return (datetime.utcnow() - delta).timestamp()
    # Try ISO date
    try:
        return datetime.fromisoformat(since).timestamp()
    except ValueError:
        return None


def _get_state(request: Request):
    """Get shared state from app."""
    return (
        request.app.state.config,
        request.app.state.collection,
        request.app.state.chat_engine,
    )


@router.post("/index")
async def webhook_index(request: Request, background_tasks: BackgroundTasks):
    """Webhook: index new articles by ID."""
    config = request.app.state.config
    collection = request.app.state.collection
    body = await request.json()
    item_ids = body.get("item_ids", [])
    if not item_ids:
        return JSONResponse({"status": "no items"}, status_code=200)

    def do_index():
        count, _new_urls = index_and_assign_items(config, collection, item_ids)
        if count > 0:
            bm25, docs = build_bm25_index(collection)
            engine = request.app.state.chat_engine
            if engine:
                engine.rebuild_index(bm25, docs)

    background_tasks.add_task(do_index)
    return JSONResponse({"status": "accepted", "items": len(item_ids)}, status_code=202)


@router.post("/index-feeds")
async def webhook_index_feeds(request: Request, background_tasks: BackgroundTasks):
    """Webhook from yarr: index recent articles from specific feeds."""
    config = request.app.state.config
    collection = request.app.state.collection
    body = await request.json()
    feed_ids = body.get("feed_ids", [])
    if not feed_ids:
        return JSONResponse({"status": "no feeds"}, status_code=200)

    def do_index():
        from .yarr_db import open_db, list_items as db_list_items

        conn = open_db(config.yarr_db)
        try:
            all_item_ids = []
            for fid in feed_ids:
                items = db_list_items(conn, feed_id=fid, limit=50)
                all_item_ids.extend(item["id"] for item in items)
        finally:
            conn.close()

        if all_item_ids:
            count, _new_urls = index_and_assign_items(config, collection, all_item_ids)
            if count > 0:
                bm25, docs = build_bm25_index(collection)
                engine = request.app.state.chat_engine
                if engine:
                    engine.rebuild_index(bm25, docs)
                log.info("Indexed %d new articles from %d feeds", count, len(feed_ids))

    background_tasks.add_task(do_index)
    return JSONResponse({"status": "accepted", "feeds": len(feed_ids)}, status_code=202)


@router.get("/task-status")
async def task_status(request: Request):
    """Return current AI background task status."""
    task = getattr(request.app.state, "ai_task", None)
    if task is None:
        return JSONResponse({"type": None, "started_at": None, "detail": ""})
    return JSONResponse(task)


@router.post("/reindex")
async def full_reindex(request: Request, background_tasks: BackgroundTasks):
    """Trigger full reindex of all yarr articles."""
    config = request.app.state.config
    collection = request.app.state.collection
    task_state = request.app.state.ai_task

    # Don't start if another task is running
    if task_state.get("type"):
        return JSONResponse({"status": "busy", "detail": task_state["detail"]}, status_code=409)

    def progress(msg):
        task_state["detail"] = msg

    embed_prov = request.app.state.embed_provider

    def do_reindex():
        task_state.update({"type": "reindex", "started_at": datetime.utcnow().isoformat(), "detail": "Starting..."})
        try:
            count = reindex_all(config, collection, embed_provider=embed_prov, on_progress=progress)
            if count > 0:
                progress("Rebuilding search index...")
                bm25, docs = build_bm25_index(collection)
                engine = request.app.state.chat_engine
                if engine:
                    engine.rebuild_index(bm25, docs)
        finally:
            task_state.update({"type": None, "started_at": None, "detail": ""})

    background_tasks.add_task(do_reindex)
    return JSONResponse({"status": "accepted"}, status_code=202)


@router.post("/recluster")
async def recluster(request: Request, background_tasks: BackgroundTasks):
    """Trigger re-clustering."""
    config = request.app.state.config
    task_state = request.app.state.ai_task

    # Don't start if another task is running
    if task_state.get("type"):
        return JSONResponse({"status": "busy", "detail": task_state["detail"]}, status_code=409)

    llm_prov = request.app.state.label_provider
    embed_prov = request.app.state.embed_provider

    def progress(msg):
        task_state["detail"] = msg

    def do_recluster():
        task_state.update({"type": "cluster", "started_at": datetime.utcnow().isoformat(), "detail": "Starting..."})
        try:
            run_clustering(config, llm_provider=llm_prov, embed_provider=embed_prov, on_progress=progress)
        finally:
            task_state.update({"type": None, "started_at": None, "detail": ""})

    background_tasks.add_task(do_recluster)
    return JSONResponse({"status": "accepted"}, status_code=202)


@router.post("/chat")
async def chat(request: Request):
    """SSE streaming chat response."""
    config, collection, engine = _get_state(request)
    if engine is None:
        return JSONResponse({"error": "AI engine not initialized"}, status_code=503)

    body = await request.json()
    query = body.get("query", "")
    history = body.get("history", [])
    topic = body.get("topic")
    tag = body.get("tag")
    since = parse_since(body.get("since"))

    if not query:
        return JSONResponse({"error": "query required"}, status_code=400)

    # Search
    results = engine.search(query, topic=topic, tag=tag, since_ts=since)

    async def event_stream():
        try:
            async for token in engine.generate_response(query, results, history):
                yield {"data": token}
            # Send sources metadata
            sources = [
                {
                    "title": r.get("title", ""),
                    "url": r.get("url", ""),
                    "published": r.get("published", ""),
                    "folder": r.get("folder", ""),
                    "feed_name": r.get("feed_name", ""),
                }
                for r in results
            ]
            yield {"data": json.dumps({"sources": sources}), "event": "sources"}
            yield {"data": "[DONE]"}
        except Exception as e:
            log.error("Chat stream error: %s", e)
            yield {"data": json.dumps({"error": str(e)}), "event": "error"}

    return EventSourceResponse(event_stream())


@router.get("/briefing")
async def briefing(request: Request):
    """SSE streaming briefing digest."""
    config = request.app.state.config
    collection = request.app.state.collection

    topic = request.query_params.get("topic")
    tag = request.query_params.get("tag")
    since = parse_since(request.query_params.get("since", "24h"))

    llm_prov = request.app.state.chat_provider

    async def event_stream():
        try:
            async for token in generate_briefing(
                llm_prov, collection,
                context_window=config.context_window,
                temperature=config.temperature,
                topic=topic, tag=tag, since_ts=since
            ):
                yield {"data": token}
            yield {"data": "[DONE]"}
        except Exception as e:
            log.error("Briefing stream error: %s", e)
            yield {"data": json.dumps({"error": str(e)}), "event": "error"}

    return EventSourceResponse(event_stream())


@router.post("/search")
async def search(request: Request):
    """Hybrid search returning ranked results."""
    config, collection, engine = _get_state(request)
    if engine is None:
        return JSONResponse({"error": "AI engine not initialized"}, status_code=503)

    body = await request.json()
    query = body.get("query", "")
    n_results = body.get("n_results", config.n_results)
    topic = body.get("topic")
    tag = body.get("tag")
    since = parse_since(body.get("since"))

    if not query:
        return JSONResponse({"error": "query required"}, status_code=400)

    results = engine.search(query, topic=topic, tag=tag, since_ts=since)
    return JSONResponse(results[:n_results])


@router.get("/clusters")
async def clusters(request: Request):
    """List topic clusters — now served directly by Go via /api/ai/clusters.
    This endpoint is kept as a fallback but should not be reached in normal operation.
    """
    return JSONResponse({"clusters": [], "message": "Clusters are served by the Go server at /api/ai/clusters."})


@router.get("/tags")
async def tags(request: Request):
    """List all AI-generated tags with article counts."""
    from .store import list_tags
    collection = request.app.state.collection
    if not collection:
        return JSONResponse([])
    return JSONResponse(list_tags(collection))


@router.get("/articles")
async def articles(request: Request):
    """Article listing is now served directly by Go at /api/ai/articles.
    This stub is kept so the Python router doesn't 404 if hit directly.
    """
    return JSONResponse([])


@router.get("/dedup-groups")
async def dedup_groups(request: Request):
    """List dedup groups for recent articles."""
    config = request.app.state.config
    collection = request.app.state.collection
    since = parse_since(request.query_params.get("since", "48h"))

    groups = find_dedup_groups(collection, threshold=config.dedup_threshold, since_ts=since)
    return JSONResponse(groups)


@router.get("/health")
async def health(request: Request):
    """Health check with subsystem status."""
    config = request.app.state.config
    collection = request.app.state.collection

    # Check Ollama
    import httpx
    ollama_ok = False
    try:
        resp = httpx.get(f"{config.ollama_url}/api/tags", timeout=5.0)
        ollama_ok = resp.status_code == 200
    except Exception:
        pass

    chroma_docs = 0
    try:
        chroma_docs = collection.count() if collection else 0
    except Exception:
        pass

    bm25_docs = 0
    engine = request.app.state.chat_engine
    if engine and engine.bm25_docs:
        bm25_docs = len(engine.bm25_docs)

    # Centroid count — indicates whether incremental topic assignment is active
    centroid_count = 0
    if config.yarr_api_url:
        from .cluster import fetch_previous_centroids
        try:
            prev = fetch_previous_centroids(config.yarr_api_url)
            centroid_count = len(prev) if prev else 0
        except Exception:
            pass

    # n_clusters is now served by Go directly — not tracked here
    status = "ok" if ollama_ok else "degraded"
    if not collection:
        status = "error"

    return JSONResponse({
        "status": status,
        "ollama": ollama_ok,
        "chroma_docs": chroma_docs,
        "bm25_docs": bm25_docs,
        "centroid_count": centroid_count,
    })


@router.get("/settings")
async def get_settings(request: Request):
    """Return current AI config with active providers."""
    config = request.app.state.config
    chat_prov = request.app.state.chat_provider
    label_prov = request.app.state.label_provider
    embed_prov = request.app.state.embed_provider
    return JSONResponse({
        "embed_provider": type(embed_prov).__name__,
        "chat_provider": chat_prov.model_name() if hasattr(chat_prov, 'model_name') else "unknown",
        "label_provider": label_prov.model_name() if hasattr(label_prov, 'model_name') else "unknown",
        "ollama_url": config.ollama_url,
        "gemini_configured": bool(config.gemini_api_key),
        "grok_configured": bool(config.grok_api_key),
        "chunk_size": config.chunk_size,
        "n_results": config.n_results,
        "distance_threshold": config.distance_threshold,
        "context_window": config.context_window,
        "temperature": config.temperature,
        "dedup_threshold": config.dedup_threshold,
        "min_cluster_size": config.min_cluster_size,
    })
