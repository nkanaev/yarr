# Stage 1: Build Go binary
FROM golang:1.23-bookworm AS go-builder

RUN apt-get update && apt-get install -y gcc libc6-dev

WORKDIR /build
COPY go.mod go.sum ./
COPY vendor/ vendor/
COPY cmd/ cmd/
COPY src/ src/

RUN CGO_ENABLED=1 go build \
    -tags "sqlite_foreign_keys sqlite_json" \
    -ldflags="-s -w" \
    -o /yarr ./cmd/yarr

# Stage 2: Python dependencies
FROM python:3.11-slim-bookworm AS py-builder

WORKDIR /build
COPY ai/pyproject.toml ./
RUN pip install --no-cache-dir --target=/pylibs \
    fastapi 'uvicorn[standard]' sse-starlette httpx \
    chromadb rank-bm25 hdbscan umap-learn scikit-learn numpy \
    trafilatura beautifulsoup4

# Stage 3: Runtime
FROM python:3.11-slim-bookworm

RUN apt-get update && apt-get install -y --no-install-recommends \
    s6 ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy Go binary
COPY --from=go-builder /yarr /usr/local/bin/yarr

# Copy Python libs and app
COPY --from=py-builder /pylibs /usr/local/lib/python3.11/site-packages
COPY ai/ /app/ai/

# Create s6 service directories
RUN mkdir -p /etc/s6/yarr /etc/s6/yarr-ai /data

# s6 service: yarr (Go)
RUN printf '#!/bin/sh\nexec /usr/local/bin/yarr \\\n  -addr 0.0.0.0:7070 \\\n  -db /data/yarr.db \\\n  -ai-url http://127.0.0.1:8484\n' > /etc/s6/yarr/run && \
    chmod +x /etc/s6/yarr/run

# s6 service: yarr-ai (Python)
RUN printf '#!/bin/sh\nexport YARR_DB=/data/yarr.db\nexport CHROMA_PATH=/data/chroma\nexport OLLAMA_URL="${OLLAMA_URL:-http://host.docker.internal:11434}"\nexport EMBED_MODEL="${EMBED_MODEL:-nomic-embed-text}"\nexport CHAT_MODEL="${CHAT_MODEL:-deepseek-r1:7b}"\ncd /app\nexec python -m ai.main\n' > /etc/s6/yarr-ai/run && \
    chmod +x /etc/s6/yarr-ai/run

EXPOSE 7070
VOLUME /data

# s6 manages both processes
CMD ["s6-svscan", "/etc/s6"]
