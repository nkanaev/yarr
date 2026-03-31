# Save to Instapaper — Design Spec

## Summary

Add a "Save to Instapaper" button to yarr's item detail toolbar. When clicked, the article is saved to the user's Instapaper account via the Full API, the item is marked as saved and read in yarr, and the button reflects the saved state persistently across sessions.

## Decisions

- **Authentication:** Full API with stored username/password, exchanged for OAuth xAuth tokens
- **UX:** Single toolbar button, no folder picker or context menu
- **Credential setup:** Added to existing settings panel (not a separate integrations section)
- **Post-save behavior:** Persistent saved state + auto-mark as read
- **Tracking:** Boolean column on items table (not a separate table)

## Data Layer

### Migration

New migration adds a boolean column to the items table:

```sql
ALTER TABLE items ADD COLUMN instapaper_saved BOOLEAN NOT NULL DEFAULT 0;
```

### Settings

Add to `settingsDefaults` in `storage/settings.go`:

- `instapaper_username` — string, default `""`
- `instapaper_password` — string, default `""`
- `instapaper_oauth_token` — string, default `""` (cached after xAuth exchange)
- `instapaper_oauth_secret` — string, default `""` (cached after xAuth exchange)

### Storage Methods

- `SetItemInstapaperSaved(id int64, saved bool)` — updates the `instapaper_saved` column
- Existing `GetItem`/`GetItems` queries updated to include `instapaper_saved` in SELECT

## Backend API

### Endpoint

`POST /api/items/:id/instapaper`

### Handler Logic

1. Look up item by ID to get its `link` and `title`
2. Read Instapaper credentials from settings
3. If credentials empty, return `400` with error message
4. If no cached OAuth token, exchange username/password for token via Instapaper xAuth (`POST https://www.instapaper.com/api/1/oauth/access_token`)
5. Cache the OAuth token/secret in settings
6. Call `POST https://www.instapaper.com/api/1/bookmarks/add` with item URL and title
7. On success: set `instapaper_saved = true` and `status = READ`
8. Return updated item as JSON

### OAuth 1.0 Implementation

Implemented using Go standard library (`crypto/hmac`, `crypto/sha1`, `net/http`). No external dependencies.

Instapaper's Full API requires a consumer key/secret (registered app) plus user access tokens obtained via xAuth. Consumer key/secret are provided via environment variables `INSTAPAPER_CLIENT_KEY` and `INSTAPAPER_CLIENT_SECRET`, following the existing pattern of `YARR_AUTH` and other env-based config in `cmd/yarr/main.go`.

### Error Responses

- `400` — credentials not configured
- `401` — Instapaper rejected credentials (clear cached token, prompt re-auth)
- `502` — Instapaper API unreachable or returned unexpected error

## Frontend

### Settings Panel

Add to existing settings area in `index.html`:

- "Instapaper Username" text input
- "Instapaper Password" password input

Wired to `api.settings.update()` on change, same as theme/font settings.

### Item Detail Toolbar

New button after existing star/unread buttons:

- **Default state:** Instapaper/save icon
- **Saved state:** Checkmark or filled icon, visually disabled
- **Already-saved items:** Button renders in saved state on load (driven by `item.instapaper_saved`)

### API Client (`api.js`)

```javascript
items: {
  saveToInstapaper: function(id) {
    return request('/api/items/' + id + '/instapaper', {method: 'POST'})
  }
}
```

### Vue State Updates

On successful save:
1. `item.instapaper_saved = true`
2. `item.status = READ`
3. Button transitions to saved state

### Error Handling

- Credentials not configured: alert/message pointing to settings
- API failure: toast/alert with failure reason

## Files Modified

- `src/storage/migration.go` — new migration
- `src/storage/item.go` — add `InstapaperSaved` field to Item struct, update queries
- `src/storage/settings.go` — add Instapaper settings defaults
- `src/server/routes.go` — new endpoint handler, register route
- `src/server/forms.go` — response struct updates if needed
- `src/assets/index.html` — settings inputs, toolbar button
- `src/assets/javascripts/app.js` — Vue method for save action
- `src/assets/javascripts/api.js` — API client method

## New Files

- `src/server/instapaper.go` — OAuth 1.0 client and Instapaper API wrapper (keeps HTTP handler clean)
