# Once Compatibility Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make yarr fully compatible with the Basecamp Once self-hosting platform (all three tiers).

**Architecture:** Shared Go code changes (health endpoint, graceful shutdown, `SECRET_KEY_BASE`, `DISABLE_SSL`) benefit all deployment modes. A separate `etc/dockerfile.once` handles Once-specific conventions (port 80, `/storage`, hooks, OCI labels). Existing Dockerfile and behavior are untouched.

**Tech Stack:** Go 1.23, SQLite, Docker, GitHub Actions

**Spec:** `docs/superpowers/specs/2026-03-20-once-compatibility-design.md`

---

## File Map

| File | Action | Responsibility |
|---|---|---|
| `src/storage/storage.go` | Modify | Add `Ping()` method for health checks |
| `src/server/server.go` | Modify | Add `SecretKeyBase`/`SecureCookie` fields, graceful shutdown |
| `src/server/routes.go` | Modify | Add `/up` health endpoint, wire new auth fields |
| `src/server/auth/auth.go` | Modify | Add `baseKey` param to `secret`, `secureCookie` param to `Authenticate` |
| `src/server/auth/middleware.go` | Modify | Add `SecretKeyBase`/`SecureCookie` fields, update call sites |
| `cmd/yarr/main.go` | Modify | Read `SECRET_KEY_BASE` and `DISABLE_SSL` env vars |
| `src/server/auth/auth_test.go` | Modify | Add tests for `SECRET_KEY_BASE` and `DISABLE_SSL` |
| `src/server/routes_test.go` | Modify | Add health endpoint tests |
| `etc/hooks/pre-backup` | Create | SQLite backup hook script |
| `etc/hooks/post-restore` | Create | Post-restore cleanup script |
| `etc/dockerfile.once` | Create | Once-compatible Dockerfile |
| `.github/workflows/build-docker-once.yml` | Create | CI workflow for Once Docker image |

---

### Task 1: Add `Ping()` to storage

**Files:**
- Modify: `src/storage/storage.go:11-13`

- [ ] **Step 1: Add `Ping` method**

In `src/storage/storage.go`, add after line 13 (closing brace of `Storage` struct):

```go
func (s *Storage) Ping() error {
	return s.db.Ping()
}
```

- [ ] **Step 2: Run existing tests to verify no regression**

Run: `go test -tags "sqlite_foreign_keys sqlite_json" ./src/storage/...`
Expected: All tests PASS

- [ ] **Step 3: Commit**

```bash
git add src/storage/storage.go
git commit -m "add Ping method to storage for health checks"
```

---

### Task 2: Update auth to support `SECRET_KEY_BASE` and `DISABLE_SSL`

**Files:**
- Modify: `src/server/auth/auth.go:12-53`
- Modify: `src/server/auth/middleware.go:12-18,31,47`
- Modify: `src/server/auth/auth_test.go`

- [ ] **Step 1: Write failing tests for `SECRET_KEY_BASE`**

Append to `src/server/auth/auth_test.go`:

```go
func TestSecret_WithBaseKey(t *testing.T) {
	// With baseKey, secret should differ from without
	withoutBase := secret("alice", "pass123", "")
	withBase := secret("alice", "pass123", "my-secret-key")
	if withoutBase == withBase {
		t.Fatal("baseKey should change the secret output")
	}

	// Same baseKey is deterministic
	withBase2 := secret("alice", "pass123", "my-secret-key")
	if withBase != withBase2 {
		t.Fatal("secret with baseKey should be deterministic")
	}

	// Different baseKey produces different output
	withBase3 := secret("alice", "pass123", "other-key")
	if withBase == withBase3 {
		t.Fatal("different baseKeys should produce different secrets")
	}
}

func TestAuthenticateAndIsAuthenticated_WithBaseKey(t *testing.T) {
	username := "admin"
	password := "secret"
	baseKey := "test-secret-key-base"

	recorder := httptest.NewRecorder()
	Authenticate(recorder, username, password, "", baseKey, true)
	cookie := recorder.Result().Cookies()[0]

	// Should authenticate with same baseKey
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(cookie)
	if !IsAuthenticated(req, username, password, baseKey) {
		t.Fatal("should authenticate with matching baseKey")
	}

	// Should NOT authenticate with different baseKey
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.AddCookie(cookie)
	if IsAuthenticated(req2, username, password, "wrong-key") {
		t.Fatal("should not authenticate with different baseKey")
	}

	// Should NOT authenticate with empty baseKey (backward compat path)
	req3 := httptest.NewRequest("GET", "/", nil)
	req3.AddCookie(cookie)
	if IsAuthenticated(req3, username, password, "") {
		t.Fatal("should not authenticate when baseKey was used but not provided")
	}
}

func TestAuthenticate_SecureCookieFlag(t *testing.T) {
	// secureCookie=true -> Secure: true
	rec1 := httptest.NewRecorder()
	Authenticate(rec1, "admin", "pass", "", "", true)
	if !rec1.Result().Cookies()[0].Secure {
		t.Fatal("cookie should be secure when secureCookie=true")
	}

	// secureCookie=false -> Secure: false
	rec2 := httptest.NewRecorder()
	Authenticate(rec2, "admin", "pass", "", "", false)
	if rec2.Result().Cookies()[0].Secure {
		t.Fatal("cookie should not be secure when secureCookie=false")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test -tags "sqlite_foreign_keys sqlite_json" ./src/server/auth/...`
Expected: FAIL — `secret` has wrong arity, `Authenticate`/`IsAuthenticated` have wrong arity

- [ ] **Step 3: Update `secret` function in `auth.go`**

Replace the `secret` function (lines 48-53) in `src/server/auth/auth.go`:

```go
func secret(msg, key, baseKey string) string {
	hmacKey := []byte(key)
	if baseKey != "" {
		// Derive key using HMAC-SHA256(baseKey, key) to avoid
		// cryptographic weakness of simple concatenation
		derivation := hmac.New(sha256.New, []byte(baseKey))
		derivation.Write([]byte(key))
		hmacKey = derivation.Sum(nil)
	}
	mac := hmac.New(sha256.New, hmacKey)
	mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}
```

- [ ] **Step 4: Update `IsAuthenticated` signature**

Replace `IsAuthenticated` (lines 12-22) in `src/server/auth/auth.go`:

```go
func IsAuthenticated(req *http.Request, username, password, baseKey string) bool {
	cookie, _ := req.Cookie("auth")
	if cookie == nil {
		return false
	}
	parts := strings.Split(cookie.Value, ":")
	if len(parts) != 2 || !StringsEqual(parts[0], username) {
		return false
	}
	return StringsEqual(parts[1], secret(username, password, baseKey))
}
```

- [ ] **Step 5: Update `Authenticate` signature**

Replace `Authenticate` (lines 24-33) in `src/server/auth/auth.go`:

```go
func Authenticate(rw http.ResponseWriter, username, password, basepath, baseKey string, secureCookie bool) {
	http.SetCookie(rw, &http.Cookie{
		Name:     "auth",
		Value:    username + ":" + secret(username, password, baseKey),
		MaxAge:   604800, // 1 week
		Path:     basepath,
		Secure:   secureCookie,
		SameSite: http.SameSiteLaxMode,
	})
}
```

- [ ] **Step 6: Update `Middleware` struct and call sites in `middleware.go`**

Replace the `Middleware` struct (lines 12-18) in `src/server/auth/middleware.go`:

```go
type Middleware struct {
	Username      string
	Password      string
	BasePath      string
	Public        []string
	DB            *storage.Storage
	SecretKeyBase string
	SecureCookie  bool
}
```

Update `IsAuthenticated` call (line 31):

```go
	if IsAuthenticated(c.Req, m.Username, m.Password, m.SecretKeyBase) {
```

Update `Authenticate` call (line 47):

```go
			Authenticate(c.Out, m.Username, m.Password, m.BasePath, m.SecretKeyBase, m.SecureCookie)
```

- [ ] **Step 7: Fix existing tests that call old signatures**

Use `replace_all` for each pattern in `src/server/auth/auth_test.go`:

**`secret` calls (3-arg → 4-arg):**
- `secret("alice", "pass123")` → `secret("alice", "pass123", "")` (lines 11, 12)
- `secret("bob", "pass123")` → `secret("bob", "pass123", "")` (line 19)
- `secret("alice", "otherpass")` → `secret("alice", "otherpass", "")` (line 23)

**`Authenticate` calls (4-arg → 6-arg):**
- Line 59: `Authenticate(recorder, username, password, basepath)` → `Authenticate(recorder, username, password, basepath, "", true)`
- Line 97: `Authenticate(recorder, "admin", "secret", "")` → `Authenticate(recorder, "admin", "secret", "", "", true)`
- Line 111: `Authenticate(recorder, "admin", "secret", "")` → `Authenticate(recorder, "admin", "secret", "", "", true)`
- Line 124: `Authenticate(recorder, "admin", "secret", "")` → `Authenticate(recorder, "admin", "secret", "", "", true)`
- Line 172: `Authenticate(recorder, "admin", "secret", "/app")` → `Authenticate(recorder, "admin", "secret", "/app", "", true)`

**`IsAuthenticated` calls (3-arg → 4-arg):**
- Line 83: `IsAuthenticated(req, username, password)` → `IsAuthenticated(req, username, password, "")`
- Line 90: `IsAuthenticated(req, "admin", "secret")` → `IsAuthenticated(req, "admin", "secret", "")`
- Line 104: `IsAuthenticated(req, "admin", "secret")` → `IsAuthenticated(req, "admin", "secret", "")`
- Line 117: `IsAuthenticated(req, "other", "secret")` → `IsAuthenticated(req, "other", "secret", "")`
- Line 130: `IsAuthenticated(req, "admin", "wrongpass")` → `IsAuthenticated(req, "admin", "wrongpass", "")`
- Line 150: `IsAuthenticated(req, "admin", "secret")` → `IsAuthenticated(req, "admin", "secret", "")`

- [ ] **Step 8: Run tests to verify they pass**

Run: `go test -tags "sqlite_foreign_keys sqlite_json" ./src/server/auth/...`
Expected: All tests PASS

- [ ] **Step 9: Commit**

```bash
git add src/server/auth/auth.go src/server/auth/middleware.go src/server/auth/auth_test.go
git commit -m "add SECRET_KEY_BASE and DISABLE_SSL support to auth"
```

---

### Task 3: Add health check endpoint

**Files:**
- Modify: `src/server/routes.go:37,42-62`
- Modify: `src/server/routes_test.go`

- [ ] **Step 1: Write failing test for health endpoint**

Append to `src/server/routes_test.go`:

```go
func TestHealthEndpoint(t *testing.T) {
	log.SetOutput(io.Discard)
	db, _ := storage.New(":memory:")
	log.SetOutput(os.Stderr)

	handler := NewServer(db, "127.0.0.1:8000").handler()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/up", nil)
	handler.ServeHTTP(recorder, request)

	if recorder.Result().StatusCode != 200 {
		t.Fatalf("expected 200, got %d", recorder.Result().StatusCode)
	}
	body, _ := io.ReadAll(recorder.Result().Body)
	if string(body) != "OK" {
		t.Fatalf("expected body 'OK', got %q", string(body))
	}
}

func TestHealthEndpoint_NotAuthGated(t *testing.T) {
	log.SetOutput(io.Discard)
	db, _ := storage.New(":memory:")
	log.SetOutput(os.Stderr)

	srv := NewServer(db, "127.0.0.1:8000")
	srv.Username = "admin"
	srv.Password = "secret"
	srv.SecureCookie = true
	handler := srv.handler()

	// Request without auth cookie should still get 200 on /up
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/up", nil)
	handler.ServeHTTP(recorder, request)

	if recorder.Result().StatusCode != 200 {
		t.Fatalf("health endpoint should not require auth, got %d", recorder.Result().StatusCode)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test -tags "sqlite_foreign_keys sqlite_json" ./src/server/...`
Expected: FAIL — route `/up` not registered, `SecureCookie` field unknown

- [ ] **Step 3: Add `SecretKeyBase` and `SecureCookie` fields to `Server` struct**

In `src/server/server.go`, add after line 29 (`KeyFile string`):

```go
	// once
	SecretKeyBase string
	SecureCookie  bool
```

- [ ] **Step 4: Add health handler and wire auth fields in `routes.go`**

In `src/server/routes.go`, add `/up` to the `Public` slice (line 37):

```go
		Public:   []string{"/static", "/fever", "/manifest.json", "/up"},
```

Pass new fields to `auth.Middleware` (add after line 38, the `DB` field):

```go
			SecretKeyBase: s.SecretKeyBase,
			SecureCookie:  s.SecureCookie,
```

Add the `/up` route after the auth middleware `if` block (line 41-42 closes with `}`), before `r.For("/", ...)` on line 43:

```go
	r.For("/up", s.handleHealth)
```

Add the handler method at the end of `routes.go`:

```go
func (s *Server) handleHealth(c *router.Context) {
	if err := s.db.Ping(); err != nil {
		c.Out.WriteHeader(http.StatusServiceUnavailable)
		c.Out.Write([]byte("ERROR"))
		return
	}
	c.Out.WriteHeader(http.StatusOK)
	c.Out.Write([]byte("OK"))
}
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test -tags "sqlite_foreign_keys sqlite_json" ./src/server/...`
Expected: All tests PASS

- [ ] **Step 6: Commit**

```bash
git add src/server/server.go src/server/routes.go src/server/routes_test.go
git commit -m "add /up health check endpoint"
```

---

### Task 4: Add graceful shutdown

**Files:**
- Modify: `src/server/server.go:50-87`

- [ ] **Step 1: Add signal handling to `Start()`**

Replace the `Start()` method (lines 50-87) in `src/server/server.go` with:

```go
func (s *Server) Start() {
	refreshRate := s.db.GetSettingsValueInt64("refresh_rate")
	s.worker.FindFavicons()
	s.worker.StartFeedCleaner()
	s.worker.SetRefreshRate(refreshRate)
	if refreshRate > 0 {
		s.worker.RefreshFeeds()
	}

	var ln net.Listener
	var err error

	if path, isUnix := strings.CutPrefix(s.Addr, "unix:"); isUnix {
		err = os.Remove(path)
		if err != nil {
			log.Print(err)
		}
		ln, err = net.Listen("unix", path)
	} else {
		ln, err = net.Listen("tcp", s.Addr)
	}

	if err != nil {
		log.Fatal(err)
	}

	httpserver := &http.Server{Handler: s.handler()}

	// Graceful shutdown: listen for SIGTERM/SIGINT
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Print("shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := httpserver.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown error: %v", err)
		}
	}()

	if s.CertFile != "" && s.KeyFile != "" {
		err = httpserver.ServeTLS(ln, s.CertFile, s.KeyFile)
	} else {
		err = httpserver.Serve(ln)
	}

	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
```

**Note:** This replacement removes the `ln.Close()` call that was on line 79 after `ServeTLS`. `Shutdown()` already closes the listener, so the explicit close is no longer needed. The existing `if err != http.ErrServerClosed` guard (line 84-86) is retained — `Shutdown()` causes `Serve`/`ServeTLS` to return `http.ErrServerClosed`, which is treated as a clean exit.

- [ ] **Step 2: Update imports in `server.go`**

Replace the import block (lines 3-13) with:

```go
import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/nkanaev/yarr/src/storage"
	"github.com/nkanaev/yarr/src/worker"
)
```

- [ ] **Step 3: Run full test suite**

Run: `go test -tags "sqlite_foreign_keys sqlite_json" ./...`
Expected: All tests PASS

- [ ] **Step 4: Commit**

```bash
git add src/server/server.go
git commit -m "add graceful shutdown on SIGTERM/SIGINT"
```

---

### Task 5: Wire `SECRET_KEY_BASE` and `DISABLE_SSL` in `main.go`

**Files:**
- Modify: `cmd/yarr/main.go:1-162`

- [ ] **Step 1: Add env var reading after flag parsing**

In `cmd/yarr/main.go`, add after line 132-133 (the `log.Fatalf("Both cert & key files are required")` block), before the `store, err := storage.New(db)` line (135):

```go
	secretKeyBase := os.Getenv("SECRET_KEY_BASE")
	secureCookie := true
	if disableSSL := os.Getenv("DISABLE_SSL"); disableSSL != "" {
		if parsed, err := strconv.ParseBool(disableSSL); err != nil {
			log.Printf("invalid DISABLE_SSL value %q, defaulting to false", disableSSL)
		} else if parsed {
			secureCookie = false
		}
	}
```

- [ ] **Step 2: Add `strconv` to imports**

The `strconv` import is not yet in `main.go`. Add `"strconv"` to the import block.

- [ ] **Step 3: Pass values to Server**

After the existing password/username assignment block (lines 152-155), add:

```go
	srv.SecretKeyBase = secretKeyBase
	srv.SecureCookie = secureCookie
```

- [ ] **Step 4: Run full test suite**

Run: `go test -tags "sqlite_foreign_keys sqlite_json" ./...`
Expected: All tests PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/yarr/main.go
git commit -m "read SECRET_KEY_BASE and DISABLE_SSL env vars"
```

---

### Task 6: Create hook scripts

**Files:**
- Create: `etc/hooks/pre-backup`
- Create: `etc/hooks/post-restore`

- [ ] **Step 1: Create pre-backup hook**

Create `etc/hooks/pre-backup`:

```sh
#!/bin/sh
set -e
sqlite3 /storage/db/yarr.db ".backup /storage/db/yarr.db.bak"
```

- [ ] **Step 2: Create post-restore hook**

Create `etc/hooks/post-restore`:

```sh
#!/bin/sh
set -e
rm -f /storage/db/yarr.db.bak
rm -f /storage/db/yarr.db-wal /storage/db/yarr.db-shm
```

- [ ] **Step 3: Make scripts executable**

```bash
chmod +x etc/hooks/pre-backup etc/hooks/post-restore
```

- [ ] **Step 4: Commit**

```bash
git add etc/hooks/pre-backup etc/hooks/post-restore
git commit -m "add Once backup/restore hook scripts"
```

---

### Task 7: Create Once Dockerfile

**Files:**
- Create: `etc/dockerfile.once`

- [ ] **Step 1: Create the Once Dockerfile**

Create `etc/dockerfile.once`:

```dockerfile
FROM golang:alpine3.21 AS builder
RUN apk add --no-cache build-base git
WORKDIR /src
COPY . .
RUN go build -tags "sqlite_foreign_keys sqlite_json" \
    -ldflags="-s -w -X 'main.Version=$(git describe --tags --abbrev=0 2>/dev/null || echo dev)' -X 'main.GitHash=$(git rev-parse --short=8 HEAD)'" \
    -o /usr/local/bin/yarr ./cmd/yarr

FROM alpine:latest
RUN apk add --no-cache ca-certificates sqlite

COPY --from=builder /usr/local/bin/yarr /usr/local/bin/yarr

RUN mkdir -p /storage/db /hooks

COPY etc/hooks/pre-backup /hooks/pre-backup
COPY etc/hooks/post-restore /hooks/post-restore
RUN chmod +x /hooks/pre-backup /hooks/post-restore

ENV DISABLE_SSL=true

LABEL org.opencontainers.image.title="yarr"
LABEL org.opencontainers.image.description="Yet Another RSS Reader - self-hosted feed aggregator"
LABEL org.opencontainers.image.version="2.6"
LABEL org.opencontainers.image.source="https://github.com/sroberts/yarr"
LABEL org.opencontainers.image.licenses="MIT"

EXPOSE 80

ENTRYPOINT ["/usr/local/bin/yarr"]
CMD ["-addr", "0.0.0.0:80", "-db", "/storage/db/yarr.db"]
```

- [ ] **Step 2: Verify Dockerfile builds locally**

```bash
docker build -f etc/dockerfile.once -t yarr:once-test .
```

Expected: Successful build

- [ ] **Step 3: Verify container runs and responds**

```bash
mkdir -p /tmp/once-test-storage/db
docker run --rm -d --name yarr-once-test \
  -p 8080:80 \
  -v /tmp/once-test-storage:/storage \
  -e SECRET_KEY_BASE="test-key" \
  -e DISABLE_SSL="true" \
  yarr:once-test

# Wait for startup
sleep 2

# Test health endpoint
curl -f http://localhost:8080/up

# Cleanup
docker stop yarr-once-test
rm -rf /tmp/once-test-storage
```

Expected: `curl` returns `OK` with HTTP 200

- [ ] **Step 4: Commit**

```bash
git add etc/dockerfile.once
git commit -m "add Once-compatible Dockerfile"
```

---

### Task 8: Create GitHub Actions workflow for Once Docker image

**Files:**
- Create: `.github/workflows/build-docker-once.yml`

- [ ] **Step 1: Create the workflow file**

Create `.github/workflows/build-docker-once.yml`:

```yaml
name: Publish Once Docker Image
on:
  push:
    tags:
      - v*
  workflow_dispatch:

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: sroberts/yarr

jobs:
  build-and-push-image:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Set up QEMU
        uses: docker/setup-qemu-action@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to the Container registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract version from tag
        id: version
        run: echo "version=${GITHUB_REF_NAME#v}" >> "$GITHUB_OUTPUT"

      - name: Build and push Docker image
        uses: docker/build-push-action@v6
        with:
          context: .
          file: ./etc/dockerfile.once
          push: true
          tags: |
            ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:once-latest
            ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:once-${{ steps.version.outputs.version }}
          platforms: linux/amd64,linux/arm64
```

- [ ] **Step 2: Commit**

```bash
git add .github/workflows/build-docker-once.yml
git commit -m "add GitHub Actions workflow for Once Docker image"
```

---

### Task 9: Final verification

- [ ] **Step 1: Run full test suite with race detection**

```bash
go test -tags "sqlite_foreign_keys sqlite_json" -race -count=1 ./...
```

Expected: All tests PASS, no race conditions

- [ ] **Step 2: Verify build**

```bash
make host
```

Expected: Binary builds successfully

- [ ] **Step 3: Review Once compatibility checklist**

Verify all items against the spec checklist:

**Tier 1:** Port 80 (dockerfile.once CMD), `/storage` (dockerfile.once CMD), Docker image (dockerfile.once + workflow)

**Tier 2:** `/hooks/pre-backup` (etc/hooks/pre-backup), `/hooks/post-restore` (etc/hooks/post-restore), `SECRET_KEY_BASE` (main.go → server → auth), `DISABLE_SSL` (main.go → server → auth cookie), graceful degradation (all env vars optional)

**Tier 3:** `GET /up` (routes.go), `/storage/db/` layout (dockerfile.once), `SIGTERM` handling (server.go), stdout logging (unchanged), OCI labels (dockerfile.once)
