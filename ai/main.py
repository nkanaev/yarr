"""yarr-ai: AI service for yarr RSS reader."""

import logging
from contextlib import asynccontextmanager

import httpx
import uvicorn
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from .config import Config
from .chat import ChatEngine
from .providers import create_embed_provider, create_llm_provider, OllamaEmbed, OllamaLLM
from .routes import router
from .search import build_bm25_index
from .store import init_store

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(name)s %(levelname)s: %(message)s",
)
log = logging.getLogger(__name__)

# Suppress noisy trafilatura warnings for empty/invalid HTML
logging.getLogger("trafilatura").setLevel(logging.CRITICAL)


async def ensure_models(config: Config):
    """Check Ollama for required models and pull any that are missing."""
    models = set([config.embed_model, config.chat_model, config.label_model])
    url = config.ollama_url.rstrip("/")

    try:
        resp = httpx.get(f"{url}/api/tags", timeout=10.0)
        resp.raise_for_status()
        available = {m["name"] for m in resp.json().get("models", [])}
        available_base = set()
        for name in available:
            available_base.add(name)
            if ":" in name:
                available_base.add(name.split(":")[0])
    except Exception as e:
        log.warning("Could not check Ollama models (is Ollama running?): %s", e)
        return

    for model in models:
        model_base = model.split(":")[0] if ":" in model else model
        if model in available_base or model_base in available_base:
            log.info("Model %s: available", model)
        else:
            log.info("Model %s: not found, pulling...", model)
            try:
                with httpx.stream(
                    "POST",
                    f"{url}/api/pull",
                    json={"name": model},
                    timeout=600.0,
                ) as pull_resp:
                    pull_resp.raise_for_status()
                    last_status = ""
                    for line in pull_resp.iter_lines():
                        if not line:
                            continue
                        import json
                        try:
                            data = json.loads(line)
                            status = data.get("status", "")
                            if status != last_status:
                                log.info("Model %s: %s", model, status)
                                last_status = status
                        except json.JSONDecodeError:
                            pass
                log.info("Model %s: pulled successfully", model)
            except Exception as e:
                log.error("Failed to pull model %s: %s", model, e)


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Initialize AI subsystems on startup."""
    config = Config.from_env()
    app.state.config = config

    log.info("yarr-ai starting...")
    log.info("ChromaDB: %s", config.chroma_path)
    log.info("yarr DB: %s", config.yarr_db or "(not configured)")

    # Create providers
    embed_provider = create_embed_provider(config)
    chat_provider = create_llm_provider(config, purpose="chat")
    label_provider = create_llm_provider(config, purpose="label")
    app.state.embed_provider = embed_provider
    app.state.chat_provider = chat_provider
    app.state.label_provider = label_provider

    # Ensure Ollama models are available (only for Ollama providers)
    needs_ollama = (
        isinstance(embed_provider, OllamaEmbed)
        or isinstance(chat_provider, OllamaLLM)
        or isinstance(label_provider, OllamaLLM)
    )
    # Check inside FallbackEmbed/FallbackLLM too
    from .providers import FallbackEmbed, FallbackLLM
    if isinstance(embed_provider, FallbackEmbed):
        needs_ollama = True  # fallback is always Ollama
    if isinstance(chat_provider, FallbackLLM) or isinstance(label_provider, FallbackLLM):
        needs_ollama = True

    if needs_ollama:
        await ensure_models(config)

    # Init ChromaDB
    app.state.collection = None
    try:
        collection, _ = init_store(config.chroma_path, embed_provider)
        app.state.collection = collection
        doc_count = collection.count()
        log.info("ChromaDB initialized: %d documents", doc_count)
    except Exception as e:
        log.error("Failed to initialize ChromaDB: %s", e)

    # Build BM25 index
    bm25_index = None
    bm25_docs = []
    if app.state.collection and app.state.collection.count() > 0:
        try:
            bm25_index, bm25_docs = build_bm25_index(app.state.collection)
            log.info("BM25 index built: %d documents", len(bm25_docs))
        except Exception as e:
            log.warning("BM25 index build failed: %s", e)

    # Init chat engine
    try:
        engine = ChatEngine(
            config, app.state.collection, bm25_index, bm25_docs,
            embed_provider=embed_provider,
            llm_provider=chat_provider,
        )
        app.state.chat_engine = engine
        log.info("Chat engine initialized")
    except Exception as e:
        log.warning("Chat engine init failed: %s", e)
        app.state.chat_engine = None

    # Task status tracking (survives page reloads)
    app.state.ai_task = {"type": None, "started_at": None, "detail": ""}

    log.info("yarr-ai ready on %s:%d", config.host, config.port)
    yield
    log.info("yarr-ai shutting down")


def create_app() -> FastAPI:
    app = FastAPI(title="yarr-ai", version="0.1.0", lifespan=lifespan)
    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_methods=["*"],
        allow_headers=["*"],
    )
    app.include_router(router)
    return app


app = create_app()


def main():
    config = Config.from_env()
    uvicorn.run(
        "ai.main:app",
        host=config.host,
        port=config.port,
        log_level="info",
    )


if __name__ == "__main__":
    main()
