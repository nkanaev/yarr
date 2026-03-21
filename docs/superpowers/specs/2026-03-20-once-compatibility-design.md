# Once Compatibility Design

**Date:** 2026-03-20
**Branch:** `feature/once-compatibility`
**Approach:** Separate Once Dockerfile + shared Go code changes

---

## Summary

Make yarr compatible with all three tiers of the Basecamp Once spec. Go code changes (health endpoint, graceful shutdown, `SECRET_KEY_BASE`, `DISABLE_SSL`) are universal improvements. Docker-specific conventions (port 80, `/storage`, hooks, OCI labels) live in a dedicated `etc/dockerfile.once` to avoid breaking existing users.

---

## Go Code Changes

### Health check endpoint

- **Route:** `GET /up` in `src/server/routes.go`
- **Behavior:** Runs `SELECT 1` against the SQLite database. Returns HTTP 200 with body `OK` on success, HTTP 503 on failure.
- **Auth:** Not gated by auth middleware. Added to the `Public` slice in the `auth.Middleware` construction in `routes.go` (alongside existing entries like `/static`, `/fever`). The route itself is registered via `r.For("/up", ...)` — the router handles `BasePath` stripping internally, consistent with all other routes.

### Graceful shutdown

- **File:** `src/server/server.go`
- **Current state:** `server.go` already uses `&http.Server{Handler: s.handler()}` with `httpserver.Serve(ln)` / `httpserver.ServeTLS(...)`. No signal handling exists.
- **Change:** Add signal handling goroutine. Before calling `Serve`/`ServeTLS`, create a signal context via `signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)`. Launch a goroutine that waits on the signal context, then calls `httpserver.Shutdown(shutdownCtx)` with a 30-second deadline. The main goroutine calls `Serve`/`ServeTLS` which blocks until shutdown completes, returning `http.ErrServerClosed` (not an error in this case).
- **`ln.Close()` removal:** The existing `ln.Close()` call after `ServeTLS` (line 79) must be removed — `Shutdown()` already closes the listener. Double-closing would be harmless but incorrect.
- **Error handling:** The existing `if err != http.ErrServerClosed { log.Fatal(err) }` guard (line 84-86) is retained as-is — `Shutdown()` causes `Serve`/`ServeTLS` to return `http.ErrServerClosed`, which this guard correctly treats as a clean exit.
- **Scope:** HTTP-only drain — feed refresh goroutines are not coordinated (they are idempotent and SQLite WAL protects against corruption).
- **Affects:** Both GUI and non-GUI paths via `platform.Start` → `server.Start`.

### `SECRET_KEY_BASE` integration

- **Read:** `os.Getenv("SECRET_KEY_BASE")` in `cmd/yarr/main.go`. No CLI flag — this is Once-injected only.
- **Data flow:** Add `SecretKeyBase string` and `SecureCookie bool` fields to `server.Server` struct. Set them in `main.go`. In `server.handler()` (where `auth.Middleware` is constructed in `routes.go`), pass these values to `auth.Middleware{SecretKeyBase: s.SecretKeyBase, SecureCookie: s.SecureCookie, ...}`.
- **HMAC change:** In `src/server/auth/auth.go`, `secret(msg, key)` becomes `secret(msg, key, baseKey)`. When `baseKey` is non-empty, derive the HMAC key using `HMAC-SHA256(baseKey, key)` — i.e., use `baseKey` as the HMAC key to hash `key`, producing a derived key. This avoids the cryptographic weakness of simple string concatenation. When `baseKey` is empty, behavior is unchanged (backward compatible).
- **Call site updates in `middleware.go`:** `IsAuthenticated(c.Req, m.Username, m.Password)` (line 31) becomes `IsAuthenticated(c.Req, m.Username, m.Password, m.SecretKeyBase)`. `Authenticate(c.Out, m.Username, m.Password, m.BasePath)` (line 47) becomes `Authenticate(c.Out, m.Username, m.Password, m.BasePath, m.SecretKeyBase, m.SecureCookie)`.
- **Session impact:** Existing sessions are invalidated when `SECRET_KEY_BASE` is set, since the HMAC key changes. This is acceptable — Once is a fresh install context.

### `DISABLE_SSL` support

- **Read:** `os.Getenv("DISABLE_SSL")` in `cmd/yarr/main.go`. Parsed with `strconv.ParseBool` (accepts `"true"`, `"1"`, `"TRUE"`, etc. — permissive, since Once's injection format may vary). On parse error (e.g., `DISABLE_SSL=maybe`), defaults to `false` (SSL enabled) and logs a warning.
- **Pass through:** `SecureCookie bool` flows through `server.Server` → `auth.Middleware` (see data flow above).
- **Cookie change:** In `src/server/auth/auth.go`, `Authenticate()` sets `Secure: m.SecureCookie` instead of hardcoded `true`.
- **Server behavior:** When `DISABLE_SSL=true` and no TLS cert/key provided, yarr already serves plain HTTP. No server-side change needed — only the cookie flag matters.
- **Auth disabled edge case:** When no username/password is configured, `auth.Middleware` is never constructed (`routes.go` line 32), so `DISABLE_SSL` has no effect on cookies. This is intentional — without auth, there is no session cookie to secure.

---

## Dockerfile and Hook Scripts

### `etc/dockerfile.once`

Same build stage as `etc/dockerfile`. Differences in runtime stage:

| Aspect | Existing Dockerfile | Once Dockerfile |
|---|---|---|
| Port | 7070 | 80 |
| DB path | `/data/yarr.db` | `/storage/db/yarr.db` |
| Storage dirs | none | `/storage/db` created |
| Hook scripts | none | `/hooks/pre-backup`, `/hooks/post-restore` |
| OCI labels | none | title, description, version, source, licenses |
| Default env | none | `DISABLE_SSL=true` |
| Extra packages | none | `sqlite3` (for hook scripts) |

### Hook scripts

**`etc/hooks/pre-backup`**
```sh
#!/bin/sh
set -e
sqlite3 /storage/db/yarr.db ".backup /storage/db/yarr.db.bak"
```

**`etc/hooks/post-restore`**
```sh
#!/bin/sh
set -e
rm -f /storage/db/yarr.db.bak
rm -f /storage/db/yarr.db-wal /storage/db/yarr.db-shm
```

Both `chmod +x` in the Dockerfile. Both scripts require `sqlite3` CLI, which is added via `apk add sqlite` in the Once Dockerfile's runtime stage.

**Note:** During backup, both `yarr.db` and `yarr.db.bak` exist in `/storage/db/`. Once will back up both files. The `.bak` is cleaned up by `post-restore` after a restore operation.

### GitHub Actions workflow

New file: `.github/workflows/build-docker-once.yml`
- Trigger: tags (`v*`) and manual dispatch (same as existing `build-docker.yml`)
- Builds from `etc/dockerfile.once`
- Pushes to `ghcr.io/sroberts/yarr` with tags: `once-latest`, `once-<version>`
- Multi-platform: `linux/amd64`, `linux/arm64`

---

## What Doesn't Change

- **Existing Dockerfile** (`etc/dockerfile`): port 7070, `/data/yarr.db`, unchanged
- **Existing env vars** (`YARR_*`): all still work
- **SMTP/VAPID/NUM_CPUS**: yarr ignores unknown env vars, so these are harmless. No code needed.
- **`/rails/storage`**: non-Rails app, skip per spec
- **Existing tests**: not modified
- **Worker coordination**: feed refresh goroutines are not drained on shutdown (idempotent, WAL-protected)

---

## New Tests

- **Health endpoint:** `GET /up` returns 200 with a valid DB, verify it's not auth-gated
- **Graceful shutdown:** Verify `SIGTERM` triggers clean exit (may be integration-level)
- **Auth with `SECRET_KEY_BASE`:** Round-trip test with base key set, verify old cookies don't authenticate
- **Auth with `DISABLE_SSL`:** Verify cookie `Secure` flag respects the setting
- **Hook scripts:** Validate scripts are syntactically correct (shellcheck or execution test)

---

## Compatibility Checklist (from Once spec)

### Tier 1 — Required
- [ ] Docker image published to a public container registry
- [ ] Application binds to `0.0.0.0:80`
- [ ] All persistent data written to `/storage`
- [ ] Application boots with no volumes pre-populated

### Tier 2 — Recommended
- [ ] `/hooks/pre-backup` creates safe SQLite snapshot
- [ ] `/hooks/post-restore` cleans up backup artifacts
- [ ] `SECRET_KEY_BASE` used for HMAC signing
- [ ] `DISABLE_SSL` respected for cookie security flag
- [ ] Application does not crash if Once env vars are absent

### Tier 3 — Full Parity
- [ ] `GET /up` returns 200 when DB is healthy
- [ ] `/storage/db/` separates database from other paths
- [ ] `SIGTERM` handled with 30-second graceful drain
- [ ] All logs to stdout/stderr
- [ ] Docker image labeled with OCI metadata
