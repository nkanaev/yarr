# Save to Instapaper Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a "Save to Instapaper" button to yarr's item detail toolbar that saves articles via Instapaper's Full API, tracks saved state persistently, and auto-marks items as read.

**Architecture:** New migration adds `instapaper_saved` boolean column to items table. A new `instapaper.go` file in `src/server/` handles OAuth 1.0 xAuth token exchange and bookmark creation. Instapaper credentials stored in existing settings key-value store. Frontend gets a toolbar button and settings fields. Consumer key/secret come from environment variables.

**Tech Stack:** Go standard library (`crypto/hmac`, `crypto/sha1`, `net/http`, `net/url`), SQLite, Vue 2, existing yarr patterns.

---

## File Structure

| Action | File | Responsibility |
|--------|------|---------------|
| Modify | `src/storage/migration.go` | New migration m11: add `instapaper_saved` column |
| Modify | `src/storage/item.go` | Add `InstapaperSaved` field to `Item` struct, update queries |
| Modify | `src/storage/settings.go` | Add Instapaper settings defaults |
| Create | `src/server/instapaper.go` | OAuth 1.0 xAuth client + Instapaper API wrapper |
| Create | `src/server/instapaper_test.go` | Tests for OAuth signing and API integration |
| Modify | `src/server/server.go` | Add Instapaper consumer key/secret fields to `Server` struct |
| Modify | `src/server/routes.go` | New endpoint handler + route registration |
| Modify | `src/assets/javascripts/api.js` | Add `saveToInstapaper` API method |
| Modify | `src/assets/index.html` | Toolbar button + settings fields |
| Modify | `src/assets/javascripts/app.js` | Vue method + loading state |
| Modify | `src/assets/javascripts/key.js` | Keyboard shortcut for save |
| Modify | `cmd/yarr/main.go` | Read `INSTAPAPER_CLIENT_KEY` / `INSTAPAPER_CLIENT_SECRET` env vars |
| Create | `src/assets/graphicarts/inbox.svg` | Instapaper save icon (Feather "inbox" icon) |

---

### Task 1: Database Migration — Add `instapaper_saved` Column

**Files:**
- Modify: `src/storage/migration.go:10-21` (migrations slice and function)

- [ ] **Step 1: Add migration function to migrations slice**

In `src/storage/migration.go`, add `m11_add_instapaper_saved` to the migrations slice:

```go
var migrations = []func(*sql.Tx) error{
	m01_initial,
	m02_feed_states_and_errors,
	m03_on_delete_actions,
	m04_item_podcasturl,
	m05_move_description_to_content,
	m06_fill_missing_dates,
	m07_add_feed_size,
	m08_normalize_datetime,
	m09_change_item_index,
	m10_add_item_medialinks,
	m11_add_instapaper_saved,
}
```

- [ ] **Step 2: Add migration function**

Append to the bottom of `src/storage/migration.go`:

```go
func m11_add_instapaper_saved(tx *sql.Tx) error {
	_, err := tx.Exec(`alter table items add column instapaper_saved boolean not null default 0`)
	return err
}
```

- [ ] **Step 3: Verify the migration compiles**

Run: `cd /Users/sroberts/Developer/yarr && go build ./...`
Expected: No errors.

- [ ] **Step 4: Commit**

```bash
git add src/storage/migration.go
git commit -m "feat: add migration for instapaper_saved column"
```

---

### Task 2: Update Item Struct and Queries

**Files:**
- Modify: `src/storage/item.go:71-81` (Item struct)
- Modify: `src/storage/item.go:261` (ListItems SELECT)
- Modify: `src/storage/item.go:279-285` (ListItems Scan)
- Modify: `src/storage/item.go:297-306` (GetItem SELECT and Scan)

- [ ] **Step 1: Add `InstapaperSaved` field to Item struct**

In `src/storage/item.go`, change the `Item` struct (line 71-81):

```go
type Item struct {
	Id              int64      `json:"id"`
	GUID            string     `json:"guid"`
	FeedId          int64      `json:"feed_id"`
	Title           string     `json:"title"`
	Link            string     `json:"link"`
	Content         string     `json:"content,omitempty"`
	Date            time.Time  `json:"date"`
	Status          ItemStatus `json:"status"`
	MediaLinks      MediaLinks `json:"media_links"`
	InstapaperSaved bool       `json:"instapaper_saved"`
}
```

- [ ] **Step 2: Update ListItems query SELECT columns**

In `src/storage/item.go`, update the `selectCols` variable in `ListItems` (around line 261):

```go
selectCols := "i.id, i.guid, i.feed_id, i.title, i.link, i.date, i.status, i.media_links, i.instapaper_saved"
if withContent {
	selectCols += ", i.content"
} else {
	selectCols += ", '' as content"
}
```

- [ ] **Step 3: Update ListItems Scan**

In `src/storage/item.go`, update the `rows.Scan` call in `ListItems` (around line 281):

```go
err = rows.Scan(
	&x.Id, &x.GUID, &x.FeedId,
	&x.Title, &x.Link, &x.Date,
	&x.Status, &x.MediaLinks, &x.InstapaperSaved, &x.Content,
)
```

- [ ] **Step 4: Update GetItem query and Scan**

In `src/storage/item.go`, update `GetItem` (around line 296-306):

```go
func (s *Storage) GetItem(id int64) *Item {
	i := &Item{}
	err := s.db.QueryRow(`
		select
			i.id, i.guid, i.feed_id, i.title, i.link, i.content,
			i.date, i.status, i.media_links, i.instapaper_saved
		from items i
		where i.id = ?
	`, id).Scan(
		&i.Id, &i.GUID, &i.FeedId, &i.Title, &i.Link, &i.Content,
		&i.Date, &i.Status, &i.MediaLinks, &i.InstapaperSaved,
	)
	if err != nil {
		log.Print(err)
		return nil
	}
	return i
}
```

- [ ] **Step 5: Add SetItemInstapaperSaved method**

Append after the `UpdateItemStatus` method (around line 317):

```go
func (s *Storage) SetItemInstapaperSaved(id int64, saved bool) bool {
	_, err := s.db.Exec(`update items set instapaper_saved = ? where id = ?`, saved, id)
	return err == nil
}
```

- [ ] **Step 6: Verify compilation**

Run: `cd /Users/sroberts/Developer/yarr && go build ./...`
Expected: No errors.

- [ ] **Step 7: Run existing tests**

Run: `cd /Users/sroberts/Developer/yarr && go test ./...`
Expected: All tests pass.

- [ ] **Step 8: Commit**

```bash
git add src/storage/item.go
git commit -m "feat: add InstapaperSaved field to Item and update queries"
```

---

### Task 3: Update Settings Defaults

**Files:**
- Modify: `src/storage/settings.go:8-20` (settingsDefaults)

- [ ] **Step 1: Add Instapaper settings to defaults**

In `src/storage/settings.go`, update `settingsDefaults()`:

```go
func settingsDefaults() map[string]interface{} {
	return map[string]interface{}{
		"filter":                  "",
		"feed":                    "",
		"feed_list_width":         300,
		"item_list_width":         300,
		"sort_newest_first":       true,
		"theme_name":              "light",
		"theme_font":              "",
		"theme_size":              1,
		"refresh_rate":            0,
		"instapaper_username":     "",
		"instapaper_password":     "",
		"instapaper_oauth_token":  "",
		"instapaper_oauth_secret": "",
	}
}
```

- [ ] **Step 2: Verify compilation**

Run: `cd /Users/sroberts/Developer/yarr && go build ./...`
Expected: No errors.

- [ ] **Step 3: Commit**

```bash
git add src/storage/settings.go
git commit -m "feat: add Instapaper credential settings defaults"
```

---

### Task 4: OAuth 1.0 xAuth Client and Instapaper API

**Files:**
- Create: `src/server/instapaper.go`
- Create: `src/server/instapaper_test.go`

- [ ] **Step 1: Write test for OAuth signature generation**

Create `src/server/instapaper_test.go`:

```go
package server

import (
	"testing"
)

func TestOAuthSignatureBaseString(t *testing.T) {
	params := map[string]string{
		"oauth_consumer_key":     "testkey",
		"oauth_nonce":            "testnonce",
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        "1234567890",
		"oauth_version":          "1.0",
		"x_auth_mode":            "client_auth",
		"x_auth_username":        "user@example.com",
		"x_auth_password":        "password123",
	}

	base := signatureBaseString("POST", "https://www.instapaper.com/api/1/oauth/access_token", params)

	// Should be: METHOD&encoded_url&encoded_params
	if base == "" {
		t.Fatal("signature base string should not be empty")
	}
	if base[:4] != "POST" {
		t.Errorf("base string should start with POST, got %s", base[:4])
	}
}

func TestOAuthHMACSHA1(t *testing.T) {
	sig := hmacSHA1Sign("consumerSecret&", "base")
	if sig == "" {
		t.Fatal("signature should not be empty")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/sroberts/Developer/yarr && go test ./src/server/ -run TestOAuth -v`
Expected: FAIL — functions not defined.

- [ ] **Step 3: Write the OAuth 1.0 and Instapaper client**

Create `src/server/instapaper.go`:

```go
package server

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	instapaperAccessTokenURL = "https://www.instapaper.com/api/1/oauth/access_token"
	instapaperAddBookmarkURL = "https://www.instapaper.com/api/1/bookmarks/add"
)

type InstapaperClient struct {
	ConsumerKey    string
	ConsumerSecret string
}

func percentEncode(s string) string {
	return url.QueryEscape(s)
}

func signatureBaseString(method, baseURL string, params map[string]string) string {
	// Sort parameter keys
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build parameter string
	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, percentEncode(k)+"="+percentEncode(params[k]))
	}
	paramStr := strings.Join(pairs, "&")

	return method + "&" + percentEncode(baseURL) + "&" + percentEncode(paramStr)
}

func hmacSHA1Sign(key, data string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func generateNonce() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func (c *InstapaperClient) signedRequest(method, endpoint string, extraParams map[string]string, oauthToken, oauthTokenSecret string) (*http.Response, error) {
	params := map[string]string{
		"oauth_consumer_key":     c.ConsumerKey,
		"oauth_nonce":            generateNonce(),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        fmt.Sprintf("%d", time.Now().Unix()),
		"oauth_version":          "1.0",
	}
	if oauthToken != "" {
		params["oauth_token"] = oauthToken
	}
	for k, v := range extraParams {
		params[k] = v
	}

	baseStr := signatureBaseString(method, endpoint, params)
	signingKey := percentEncode(c.ConsumerSecret) + "&" + percentEncode(oauthTokenSecret)
	params["oauth_signature"] = hmacSHA1Sign(signingKey, baseStr)

	// Build POST body (all params)
	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	req, err := http.NewRequest(method, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return http.DefaultClient.Do(req)
}

// GetAccessToken exchanges username/password for OAuth access token via xAuth.
func (c *InstapaperClient) GetAccessToken(username, password string) (token, secret string, err error) {
	extra := map[string]string{
		"x_auth_mode":     "client_auth",
		"x_auth_username": username,
		"x_auth_password": password,
	}

	resp, err := c.signedRequest("POST", instapaperAccessTokenURL, extra, "", "")
	if err != nil {
		return "", "", fmt.Errorf("instapaper auth request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("instapaper auth read failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("instapaper auth rejected (status %d): %s", resp.StatusCode, string(body))
	}

	vals, err := url.ParseQuery(string(body))
	if err != nil {
		return "", "", fmt.Errorf("instapaper auth parse failed: %w", err)
	}

	return vals.Get("oauth_token"), vals.Get("oauth_token_secret"), nil
}

// AddBookmark saves a URL to Instapaper.
func (c *InstapaperClient) AddBookmark(oauthToken, oauthTokenSecret, articleURL, title string) error {
	extra := map[string]string{
		"url": articleURL,
	}
	if title != "" {
		extra["title"] = title
	}

	resp, err := c.signedRequest("POST", instapaperAddBookmarkURL, extra, oauthToken, oauthTokenSecret)
	if err != nil {
		return fmt.Errorf("instapaper add bookmark failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("instapaper add bookmark rejected (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd /Users/sroberts/Developer/yarr && go test ./src/server/ -run TestOAuth -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add src/server/instapaper.go src/server/instapaper_test.go
git commit -m "feat: add OAuth 1.0 xAuth client for Instapaper API"
```

---

### Task 5: Server Struct and Environment Variables

**Files:**
- Modify: `src/server/server.go:19-38` (Server struct)
- Modify: `cmd/yarr/main.go:49-176` (main function)

- [ ] **Step 1: Add Instapaper fields to Server struct**

In `src/server/server.go`, add fields to the `Server` struct after the `SecureCookie` field (around line 37):

```go
type Server struct {
	Addr        string
	db          *storage.Storage
	worker      *worker.Worker
	cache       map[string]interface{}
	cache_mutex *sync.Mutex

	BasePath string

	// auth
	Username string
	Password string
	// https
	CertFile string
	KeyFile  string

	// once
	SecretKeyBase string
	SecureCookie  bool

	// instapaper
	InstapaperClientKey    string
	InstapaperClientSecret string
}
```

- [ ] **Step 2: Read environment variables in main.go**

In `cmd/yarr/main.go`, after the `srv.SecureCookie = secureCookie` line (around line 169), add:

```go
srv.InstapaperClientKey = os.Getenv("INSTAPAPER_CLIENT_KEY")
srv.InstapaperClientSecret = os.Getenv("INSTAPAPER_CLIENT_SECRET")
```

- [ ] **Step 3: Verify compilation**

Run: `cd /Users/sroberts/Developer/yarr && go build ./...`
Expected: No errors.

- [ ] **Step 4: Commit**

```bash
git add src/server/server.go cmd/yarr/main.go
git commit -m "feat: add Instapaper client key/secret to server config"
```

---

### Task 6: API Endpoint Handler

**Files:**
- Modify: `src/server/routes.go:28-69` (handler route registration)
- Modify: `src/server/routes.go` (add new handler function)

- [ ] **Step 1: Register the new route**

In `src/server/routes.go`, inside the `handler()` method, add the new route after the `/api/items/:id` line (around line 60):

```go
r.For("/api/items/:id/instapaper", s.handleItemInstapaper)
```

- [ ] **Step 2: Add the handler function**

Append the handler function to `src/server/routes.go`, before `handleSettings`:

```go
func (s *Server) handleItemInstapaper(c *router.Context) {
	if c.Req.Method != "POST" {
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if s.InstapaperClientKey == "" || s.InstapaperClientSecret == "" {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Instapaper client key/secret not configured. Set INSTAPAPER_CLIENT_KEY and INSTAPAPER_CLIENT_SECRET environment variables.",
		})
		return
	}

	id, err := c.VarInt64("id")
	if err != nil {
		c.Out.WriteHeader(http.StatusBadRequest)
		return
	}

	item := s.db.GetItem(id)
	if item == nil {
		c.Out.WriteHeader(http.StatusNotFound)
		return
	}

	// Get Instapaper credentials from settings
	username, _ := s.db.GetSettingsValue("instapaper_username").(string)
	password, _ := s.db.GetSettingsValue("instapaper_password").(string)
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Instapaper credentials not configured. Add your username and password in Settings.",
		})
		return
	}

	client := &InstapaperClient{
		ConsumerKey:    s.InstapaperClientKey,
		ConsumerSecret: s.InstapaperClientSecret,
	}

	// Get or create OAuth token
	oauthToken, _ := s.db.GetSettingsValue("instapaper_oauth_token").(string)
	oauthSecret, _ := s.db.GetSettingsValue("instapaper_oauth_secret").(string)
	if oauthToken == "" {
		oauthToken, oauthSecret, err = client.GetAccessToken(username, password)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Instapaper authentication failed. Check your username and password.",
			})
			return
		}
		s.db.UpdateSettings(map[string]interface{}{
			"instapaper_oauth_token":  oauthToken,
			"instapaper_oauth_secret": oauthSecret,
		})
	}

	// Save bookmark
	err = client.AddBookmark(oauthToken, oauthSecret, item.Link, item.Title)
	if err != nil {
		// Token may be stale — clear and retry once
		log.Printf("instapaper save failed, retrying with fresh token: %v", err)
		s.db.UpdateSettings(map[string]interface{}{
			"instapaper_oauth_token":  "",
			"instapaper_oauth_secret": "",
		})
		oauthToken, oauthSecret, err = client.GetAccessToken(username, password)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Instapaper authentication failed. Check your username and password.",
			})
			return
		}
		s.db.UpdateSettings(map[string]interface{}{
			"instapaper_oauth_token":  oauthToken,
			"instapaper_oauth_secret": oauthSecret,
		})
		err = client.AddBookmark(oauthToken, oauthSecret, item.Link, item.Title)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusBadGateway, map[string]string{
				"error": "Failed to save to Instapaper: " + err.Error(),
			})
			return
		}
	}

	// Update item state
	s.db.SetItemInstapaperSaved(id, true)
	s.db.UpdateItemStatus(id, storage.READ)

	// Return updated item
	updatedItem := s.db.GetItem(id)
	c.JSON(http.StatusOK, updatedItem)
}
```

- [ ] **Step 3: Verify compilation**

Run: `cd /Users/sroberts/Developer/yarr && go build ./...`
Expected: No errors.

- [ ] **Step 4: Run existing tests**

Run: `cd /Users/sroberts/Developer/yarr && go test ./...`
Expected: All tests pass.

- [ ] **Step 5: Commit**

```bash
git add src/server/routes.go
git commit -m "feat: add POST /api/items/:id/instapaper endpoint"
```

---

### Task 7: Frontend — Icon, API Method, and Toolbar Button

**Files:**
- Create: `src/assets/graphicarts/inbox.svg`
- Modify: `src/assets/javascripts/api.js:73-86`
- Modify: `src/assets/index.html:310-345`
- Modify: `src/assets/javascripts/app.js:205-272` (data)
- Modify: `src/assets/javascripts/app.js:425` (methods)

- [ ] **Step 1: Add the inbox icon**

Create `src/assets/graphicarts/inbox.svg` (Feather icon):

```svg
<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-inbox"><polyline points="22 12 16 12 14 15 10 15 8 12 2 12"></polyline><path d="M5.45 5.11L2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11z"></path></svg>
```

- [ ] **Step 2: Add API method**

In `src/assets/javascripts/api.js`, add `saveToInstapaper` to the `items` object after the `mark_read` method (around line 84):

```javascript
saveToInstapaper: function(id) {
  return api('post', './api/items/' + id + '/instapaper')
},
```

- [ ] **Step 3: Add loading state to Vue data**

In `src/assets/javascripts/app.js`, add `instapaper` to the `loading` object (around line 240):

```javascript
'loading': {
  'feeds': 0,
  'newfeed': false,
  'items': false,
  'readability': false,
  'instapaper': false,
},
```

- [ ] **Step 4: Add `saveToInstapaper` Vue method**

In `src/assets/javascripts/app.js`, add the method inside `methods:` (after `toggleReadability`, around line 702):

```javascript
saveToInstapaper: function(item) {
  if (!item || !item.link || item.instapaper_saved) return
  this.loading.instapaper = true
  api.items.saveToInstapaper(item.id).then(function(resp) {
    vm.loading.instapaper = false
    if (!resp.ok) {
      return resp.json().then(function(data) {
        alert(data.error || 'Failed to save to Instapaper')
      })
    }
    return resp.json().then(function(data) {
      vm.itemSelectedDetails.instapaper_saved = true
      vm.itemSelectedDetails.status = 'read'
      var itemInList = vm.items.find(function(i) { return i.id == item.id })
      if (itemInList) {
        itemInList.status = 'read'
        itemInList.instapaper_saved = true
      }
      if (vm.feedStats[item.feed_id]) {
        // Decrement unread count if it was unread before
        var stat = vm.feedStats[item.feed_id]
        if (item.status == 'unread' && stat.unread > 0) {
          stat.unread -= 1
        }
      }
    })
  }.bind(this))
},
```

- [ ] **Step 5: Add toolbar button to index.html**

In `src/assets/index.html`, add the Instapaper button after the "Read Here" button and before the external link (between line 342 and 343):

```html
<button class="toolbar-item"
        :class="{active: itemSelectedDetails.instapaper_saved}"
        :disabled="itemSelectedDetails.instapaper_saved || loading.instapaper"
        @click="saveToInstapaper(itemSelectedDetails)"
        title="Save to Instapaper">
    <span class="icon" :class="{'icon-loading': loading.instapaper}">
        <template v-if="itemSelectedDetails.instapaper_saved">{% inline "check.svg" %}</template>
        <template v-else>{% inline "inbox.svg" %}</template>
    </span>
</button>
```

- [ ] **Step 6: Verify the Go build succeeds (templates are embedded)**

Run: `cd /Users/sroberts/Developer/yarr && go build ./...`
Expected: No errors.

- [ ] **Step 7: Commit**

```bash
git add src/assets/graphicarts/inbox.svg src/assets/javascripts/api.js src/assets/index.html src/assets/javascripts/app.js
git commit -m "feat: add Instapaper save button to item toolbar"
```

---

### Task 8: Frontend — Settings Fields

**Files:**
- Modify: `src/assets/index.html:102-110` (settings dropdown, after "Show first" section)

- [ ] **Step 1: Add Instapaper settings section to the dropdown**

In `src/assets/index.html`, after the "Show first" section (around line 110, after the `</div>` that closes the New/Old buttons), add:

```html
<div class="dropdown-divider"></div>
<header class="dropdown-header" role="heading" aria-level="2">Instapaper</header>
<div class="px-3 py-1">
    <input type="text"
           class="form-control form-control-sm mb-1"
           placeholder="Username or email"
           :value="instapaperUsername"
           @change="updateInstapaperCredentials('instapaper_username', $event.target.value)">
    <input type="password"
           class="form-control form-control-sm"
           placeholder="Password"
           :value="instapaperPassword"
           @change="updateInstapaperCredentials('instapaper_password', $event.target.value)">
</div>
```

- [ ] **Step 2: Add Vue data properties for Instapaper credentials**

In `src/assets/javascripts/app.js`, add to the `data` function return object (after `refreshRateOptions`, around line 271):

```javascript
'instapaperUsername': s.instapaper_username || '',
'instapaperPassword': s.instapaper_password || '',
```

- [ ] **Step 3: Add `updateInstapaperCredentials` Vue method**

In `src/assets/javascripts/app.js`, add the method inside `methods:` (after `saveToInstapaper`):

```javascript
updateInstapaperCredentials: function(key, value) {
  if (key === 'instapaper_username') this.instapaperUsername = value
  if (key === 'instapaper_password') this.instapaperPassword = value
  var update = {}
  update[key] = value
  // Clear cached OAuth tokens when credentials change
  update['instapaper_oauth_token'] = ''
  update['instapaper_oauth_secret'] = ''
  api.settings.update(update)
},
```

- [ ] **Step 4: Verify the Go build succeeds**

Run: `cd /Users/sroberts/Developer/yarr && go build ./...`
Expected: No errors.

- [ ] **Step 5: Commit**

```bash
git add src/assets/index.html src/assets/javascripts/app.js
git commit -m "feat: add Instapaper credential fields to settings dropdown"
```

---

### Task 9: Keyboard Shortcut

**Files:**
- Modify: `src/assets/javascripts/key.js:17-75` (shortcutFunctions)
- Modify: `src/assets/javascripts/key.js:78-114` (keybindings/codebindings)
- Modify: `src/assets/index.html:420-442` (shortcuts help modal)

- [ ] **Step 1: Add shortcut function**

In `src/assets/javascripts/key.js`, add to `shortcutFunctions` (after `toggleItemStarred`, around line 41):

```javascript
saveToInstapaper: function() {
  if (vm.itemSelected != null) {
    vm.saveToInstapaper(vm.itemSelectedDetails)
  }
},
```

- [ ] **Step 2: Add keybinding**

In `src/assets/javascripts/key.js`, add to `keybindings` (around line 85):

```javascript
"I": shortcutFunctions.saveToInstapaper,
```

And add to `codebindings` (note: Shift+I is already different from `i` since keybindings uses `event.key`):

No codebinding needed — `"I"` (uppercase) in `keybindings` is distinct from `"i"` (lowercase) because `event.key` is case-sensitive. The `codebindings` use `event.code` which doesn't distinguish case. Since `KeyI` is already mapped to `toggleReadability`, we only add to `keybindings`.

- [ ] **Step 3: Update shortcuts help modal**

In `src/assets/index.html`, add a row to the shortcuts table (after the `<kbd>i</kbd>` / "read here" row, around line 438):

```html
<tr><td><kbd>I</kbd></td>               <td>save to Instapaper</td></tr>
```

- [ ] **Step 4: Verify the Go build succeeds**

Run: `cd /Users/sroberts/Developer/yarr && go build ./...`
Expected: No errors.

- [ ] **Step 5: Commit**

```bash
git add src/assets/javascripts/key.js src/assets/index.html
git commit -m "feat: add Shift+I keyboard shortcut for save to Instapaper"
```

---

### Task 10: Final Integration Test

**Files:**
- No new files

- [ ] **Step 1: Run full test suite**

Run: `cd /Users/sroberts/Developer/yarr && go test ./... -v`
Expected: All tests pass.

- [ ] **Step 2: Build the binary**

Run: `cd /Users/sroberts/Developer/yarr && go build -o yarr ./cmd/yarr/`
Expected: Binary builds successfully.

- [ ] **Step 3: Verify the binary starts**

Run: `cd /Users/sroberts/Developer/yarr && ./yarr -version`
Expected: Prints version string.

- [ ] **Step 4: Clean up binary**

Run: `rm /Users/sroberts/Developer/yarr/yarr`

- [ ] **Step 5: Commit (if any fixes were needed)**

Only commit if fixes were applied in earlier steps.
