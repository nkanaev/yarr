"""Configuration from environment variables."""

import os
from dataclasses import dataclass


@dataclass
class Config:
    yarr_db: str = ""
    chroma_path: str = "./data/chroma"

    # Ollama (local, default)
    ollama_url: str = "http://localhost:11434"
    embed_model: str = "nomic-embed-text"
    chat_model: str = "deepseek-r1:7b"
    label_model: str = "llama3.2:latest"

    # Provider selection: "auto", "ollama", "gemini", "grok"
    embed_provider: str = "auto"
    llm_provider: str = "auto"

    # Gemini (Google)
    gemini_api_key: str = ""
    gemini_embed_model: str = "gemini-embedding-001"
    gemini_chat_model: str = "gemini-2.5-flash"

    # Grok (xAI)
    grok_api_key: str = ""
    grok_chat_model: str = "grok-3-mini"

    # AI parameters
    chunk_size: int = 500
    chunk_overlap: int = 50
    n_results: int = 10
    distance_threshold: float = 0.7
    context_window: int = 4096
    temperature: float = 0.3
    max_history: int = 3
    dedup_threshold: float = 0.92
    min_cluster_size: int = 10
    host: str = "0.0.0.0"
    port: int = 8484

    @classmethod
    def from_env(cls) -> "Config":
        ollama_url = os.environ.get("OLLAMA_URL", "http://localhost:11434")
        if ollama_url and not ollama_url.startswith(("http://", "https://")):
            ollama_url = "http://" + ollama_url

        return cls(
            yarr_db=os.environ.get("YARR_DB", ""),
            chroma_path=os.environ.get("CHROMA_PATH", "./data/chroma"),
            ollama_url=ollama_url,
            embed_model=os.environ.get("EMBED_MODEL", "nomic-embed-text"),
            chat_model=os.environ.get("CHAT_MODEL", "deepseek-r1:7b"),
            label_model=os.environ.get("LABEL_MODEL", "llama3.2:latest"),
            embed_provider=os.environ.get("EMBED_PROVIDER", "auto"),
            llm_provider=os.environ.get("LLM_PROVIDER", "auto"),
            gemini_api_key=os.environ.get("GEMINI_API_KEY", ""),
            gemini_embed_model=os.environ.get("GEMINI_EMBED_MODEL", "gemini-embedding-001"),
            gemini_chat_model=os.environ.get("GEMINI_CHAT_MODEL", "gemini-2.5-flash"),
            grok_api_key=os.environ.get("GROK_API_KEY", ""),
            grok_chat_model=os.environ.get("GROK_CHAT_MODEL", "grok-3-mini"),
            chunk_size=int(os.environ.get("CHUNK_SIZE", "500")),
            chunk_overlap=int(os.environ.get("CHUNK_OVERLAP", "50")),
            n_results=int(os.environ.get("N_RESULTS", "10")),
            distance_threshold=float(os.environ.get("DISTANCE_THRESHOLD", "0.7")),
            context_window=int(os.environ.get("CONTEXT_WINDOW", "4096")),
            temperature=float(os.environ.get("TEMPERATURE", "0.3")),
            max_history=int(os.environ.get("MAX_HISTORY", "3")),
            dedup_threshold=float(os.environ.get("DEDUP_THRESHOLD", "0.92")),
            min_cluster_size=int(os.environ.get("MIN_CLUSTER_SIZE", "10")),
            host=os.environ.get("AI_HOST", "0.0.0.0"),
            port=int(os.environ.get("AI_PORT", "8484")),
        )
