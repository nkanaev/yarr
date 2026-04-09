# AGENTS.md

Architecture and contribution guide for AI coding agents working on the yarr codebase.

## Build Commands

```bash
# Development
make serve          # Run dev server with local.db (no AI features)
make test           # Run all Go tests
make host           # Build for current OS → out/yarr

# Cross-compile (requires zig >= 0.14.0)
make linux_amd64
make linux_arm64
make darwin_arm64   # Native build only

# Docker (full stack: Go + Python AI)
docker compose up -d --build
```

### Build Requirements

- Go >= 1.23, C compiler (GCC/Clang) — CGO_ENABLED=1 required for SQLite
- Python >= 3.10 (for AI service in `ai/`)
- Zig >= 0.14.0 (for cross-compilation only)

### Build Tags

Always required: `sqlite_foreign_keys sqlite_json`. GUI builds add: `gui`.

### Running a Single Test

```bash
go test -tags "sqlite_foreign_keys sqlite_json" ./src/storage/ -run TestName
go test -tags "sqlite_foreign_keys sqlite_json" ./src/storage/ -v -run TestListItems
go test -tags "sqlite_foreign_keys sqlite_json" ./src/server/ -run TestFever
```

### Full Build Check

```bash
go build -tags "sqlite_foreign_keys sqlite_json" ./...
go test  -tags "sqlite_foreign_keys sqlite_json" ./...
```

## Project Structure

```
cmd/yarr/              # Entry point (CLI flags, auth, server startup)
src/
  server/              # HTTP layer
    routes.go          #   Route registration + all JSON API handlers
    ai_proxy.go        #   Reverse proxy to Python AI service (SSE passthrough)
    ai_clusters.go     #   Cluster/tag/article endpoints served from SQLite
  storage/             # SQLite data layer
    storage.go         #   DB connection + migrations runner
    feed.go            #   Feed/folder CRUD
    item.go            #   Item CRUD, search, status updates
    ai_cluster.go      #   AI cluster/tag tables + queries
    migration.go       #   Migration registry (append-only)
  worker/              # Background jobs (feed refresh, cleanup, AI webhook)
  parser/              # Feed parsing (RSS, Atom, JSON Feed, RDF)
  content/             # Content processing (sanitizer, readability, scraper)
  assets/              # Embedded web UI (Vue 2 SPA)
    index.html         #   Main SPA template (Go template delimiters {% %})
    javascripts/       #   vue.min.js, api.js, app.js, key.js
    stylesheets/       #   bootstrap.min.css, app.css
    graphicarts/       #   SVG icons (inline'd via {% inline "name.svg" %})
  platform/            # OS-specific code (gui/guiless builds)
ai/                    # Python AI service (FastAPI)
  providers.py         #   Embed/LLM provider abstraction (Ollama, Gemini, Grok)
  cluster.py           #   HDBSCAN clustering + LLM labeling
  chat.py              #   RAG chat engine with SSE streaming
  indexer.py           #   Article embedding pipeline
  routes.py            #   FastAPI endpoints + background task management
  store.py             #   ChromaDB + BM25 hybrid retrieval
```

## Go Code Style

- **Error handling:** Log via `log.Print()`, return nil/false. Never `panic()` in handlers.
- **Database:** Parameterized queries only (`?` placeholders). Transactions for bulk ops. Always `defer tx.Rollback()` before work, call `tx.Commit()` on success.
- **Route handlers:** Extract params → switch on `c.Req.Method` → fetch/mutate → respond. Use `c.VarInt64()`, `c.JSON()`, `c.HTML()`.
- **Naming:** `PascalCase` exported, `camelCase` unexported. Structs: `Feed`, `Item`, `Folder`, `ArticleResult`.
- **Imports:** Standard lib → third-party → internal, separated by blank lines.
- **JSON tags:** Always add `json:"snake_case"` to struct fields returned by API.
- **Migrations:** NEVER modify existing migration functions. ALWAYS append new functions to the `migrations` array in `migration.go`.

## Vue 2 Frontend

- **Architecture:** Single-page app. Vue 2.6 instance in `app.js`, mounted on `#app`. No build step — all JS served directly.
- **Template syntax:** Go templates with `{% %}` delimiters render initial state; Vue handles all client-side reactivity with `{{ }}` mustaches.
- **Components:** `drag`, `dropdown`, `modal`, `relative-time` registered globally via `Vue.component()`.
- **State:** All app state lives in `vm.data`. Settings sync to server via `api.settings.update()`.
- **API client:** `api.js` defines `window.api` with namespaced methods (`api.feeds.*`, `api.items.*`, `api.ai.*`). Uses `fetch()` with `x-requested-by: yarr` header on mutations.
- **Keyboard shortcuts:** `key.js` maps keys to functions on the global `vm` instance.
- **CSS theming:** Class-based themes on `<body>`: `.theme-light`, `.theme-sepia`, `.theme-night`. No CSS custom properties — use theme-scoped selectors (e.g., `.theme-night #chat-panel`).
- **AI features:** Conditionally rendered with `v-if="aiEnabled"`. AI panels (chat, briefing) live outside `#app` div and are toggled via DOM manipulation + CSS classes.
- **Icons:** SVG files in `graphicarts/`, inlined at build time via `{% inline "icon.svg" %}`.

## Python AI Service Style

- **Type hints:** Modern syntax (`list[str]`, `dict[int, str]`, `str | None`).
- **Error handling:** Log via `logging`, raise exceptions for API failures. Fallback to Ollama when external providers error.
- **Providers:** Use `providers.py` factory functions. Never call Ollama/Gemini/Grok HTTP directly outside `providers.py`.
- **Naming:** `snake_case` for functions/variables, `PascalCase` for classes.
- **Async:** `async/await` for streaming chat/briefing endpoints. `httpx` for synchronous embed calls.
- **Cluster pipeline:** Results are POSTed to Go endpoints (`/api/ai/clusters/centroids`, `/api/ai/articles`) for SQLite persistence.

## AI Service Configuration

Env vars: `EMBED_PROVIDER` (ollama|gemini|auto), `LLM_PROVIDER` (ollama|gemini|grok|auto), `GEMINI_API_KEY`, `GROK_API_KEY`, `OLLAMA_URL`, `YARR_DB`, `CHROMA_PATH`, `EMBED_MODEL`, `CHAT_MODEL`.

The Go server receives `-ai-url` flag (or `YARR_AI_URL` env) pointing to the Python service (default `http://127.0.0.1:8484`). When set, `aiEnabled` is true and AI UI buttons render.

## Never Do

- **Never** modify existing migration functions — append new ones
- **Never** change `parser.Feed` or `parser.Item` structs without updating `worker.ConvertItems()`
- **Never** use `panic()` in request handlers
- **Never** commit secrets or credentials (`.env` is gitignored)
- **Never** bypass HTML sanitization for user content
- **Never** call Ollama/Gemini HTTP directly outside `ai/providers.py`
- **Never** add a JS build step or framework — the frontend is vanilla JS + Vue 2 CDN

## Quick Reference

```go
// Status values: storage.UNREAD=0, READ=1, STARRED=2
// Route params: c.VarInt64("id"), c.QueryInt64("id"), c.Vars["param"]
// HTML response: c.HTML(http.StatusOK, assets.Template("index.html"), data)
// JSON response: c.JSON(http.StatusOK, item)
```

```javascript
// Vue instance: vm (global)
// API calls:    api.feeds.list(), api.items.get(id), api.ai.chat(query, history)
// Toast:        vm.toast("message") or vm.toast("error msg", "error")
// Settings:     api.settings.update({key: value})
```

## Deployment

- **Docker Compose:** `docker compose up -d --build` — runs Go + Python via s6-overlay
- **Kubernetes:** Image at `ghcr.io/cargaona/yarr/yarr:<branch>`. CI builds on push to `main`, `master`, `feat/*`.
- **Manifests:** External repo, deployment references the GHCR image with `imagePullPolicy: Always`.

## License

See `license` file in repository root.
