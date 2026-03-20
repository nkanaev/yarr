package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecret(t *testing.T) {
	// secret is deterministic: same inputs produce same output
	s1 := secret("alice", "pass123", "")
	s2 := secret("alice", "pass123", "")
	if s1 != s2 {
		t.Fatal("secret is not deterministic")
	}

	// different inputs produce different outputs
	s3 := secret("bob", "pass123", "")
	if s1 == s3 {
		t.Fatal("different usernames should produce different secrets")
	}

	s4 := secret("alice", "otherpass", "")
	if s1 == s4 {
		t.Fatal("different passwords should produce different secrets")
	}

	// output is hex-encoded (64 chars for SHA-256)
	if len(s1) != 64 {
		t.Fatalf("expected 64 hex chars, got %d", len(s1))
	}
}

func TestStringsEqual(t *testing.T) {
	if !StringsEqual("hello", "hello") {
		t.Fatal("equal strings should match")
	}
	if StringsEqual("hello", "world") {
		t.Fatal("different strings should not match")
	}
	if StringsEqual("hello", "Hello") {
		t.Fatal("comparison should be case-sensitive")
	}
	if StringsEqual("", "notempty") {
		t.Fatal("empty vs non-empty should not match")
	}
	if !StringsEqual("", "") {
		t.Fatal("two empty strings should match")
	}
}

func TestAuthenticateAndIsAuthenticated(t *testing.T) {
	username := "admin"
	password := "secret"
	basepath := ""

	// Authenticate sets a cookie
	recorder := httptest.NewRecorder()
	Authenticate(recorder, username, password, basepath, "", true)

	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != "auth" {
		t.Fatalf("expected cookie name 'auth', got %q", cookie.Name)
	}
	if cookie.MaxAge != 604800 {
		t.Fatalf("expected MaxAge 604800, got %d", cookie.MaxAge)
	}
	if !cookie.Secure {
		t.Fatal("cookie should be secure")
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatal("cookie should have SameSite=Lax")
	}

	// Request with the auth cookie should be authenticated
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(cookie)
	if !IsAuthenticated(req, username, password, "") {
		t.Fatal("should be authenticated with valid cookie")
	}
}

func TestIsAuthenticated_NoCookie(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	if IsAuthenticated(req, "admin", "secret", "") {
		t.Fatal("should not be authenticated without cookie")
	}
}

func TestIsAuthenticated_TamperedHMAC(t *testing.T) {
	recorder := httptest.NewRecorder()
	Authenticate(recorder, "admin", "secret", "", "", true)

	cookie := recorder.Result().Cookies()[0]
	cookie.Value = "admin:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(cookie)
	if IsAuthenticated(req, "admin", "secret", "") {
		t.Fatal("tampered HMAC should not authenticate")
	}
}

func TestIsAuthenticated_WrongUsername(t *testing.T) {
	recorder := httptest.NewRecorder()
	Authenticate(recorder, "admin", "secret", "", "", true)

	cookie := recorder.Result().Cookies()[0]

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(cookie)
	if IsAuthenticated(req, "other", "secret", "") {
		t.Fatal("wrong username should not authenticate")
	}
}

func TestIsAuthenticated_WrongPassword(t *testing.T) {
	recorder := httptest.NewRecorder()
	Authenticate(recorder, "admin", "secret", "", "", true)

	cookie := recorder.Result().Cookies()[0]

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(cookie)
	if IsAuthenticated(req, "admin", "wrongpass", "") {
		t.Fatal("wrong password should not authenticate")
	}
}

func TestIsAuthenticated_MalformedCookie(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"no colon", "adminsecret"},
		{"empty value", ""},
		{"too many colons", "a:b:c"},
		{"only colon", ":"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(&http.Cookie{Name: "auth", Value: tt.value})
			if IsAuthenticated(req, "admin", "secret", "") {
				t.Fatal("malformed cookie should not authenticate")
			}
		})
	}
}

func TestLogout(t *testing.T) {
	recorder := httptest.NewRecorder()
	Logout(recorder, "")

	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	if cookies[0].MaxAge != -1 {
		t.Fatalf("expected MaxAge -1, got %d", cookies[0].MaxAge)
	}
}

func TestAuthenticate_WithBasePath(t *testing.T) {
	recorder := httptest.NewRecorder()
	Authenticate(recorder, "admin", "secret", "/app", "", true)

	cookie := recorder.Result().Cookies()[0]
	if cookie.Path != "/app" {
		t.Fatalf("expected cookie path '/app', got %q", cookie.Path)
	}
}

func TestUnsafeMethod(t *testing.T) {
	unsafe := []string{"POST", "PUT", "DELETE"}
	safe := []string{"GET", "HEAD", "OPTIONS", "PATCH"}

	for _, m := range unsafe {
		if !unsafeMethod(m) {
			t.Errorf("%s should be unsafe", m)
		}
	}
	for _, m := range safe {
		if unsafeMethod(m) {
			t.Errorf("%s should not be unsafe", m)
		}
	}
}

func TestSecret_WithBaseKey(t *testing.T) {
	withoutBase := secret("alice", "pass123", "")
	withBase := secret("alice", "pass123", "my-secret-key")
	if withoutBase == withBase {
		t.Fatal("baseKey should change the secret output")
	}
	withBase2 := secret("alice", "pass123", "my-secret-key")
	if withBase != withBase2 {
		t.Fatal("secret with baseKey should be deterministic")
	}
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

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(cookie)
	if !IsAuthenticated(req, username, password, baseKey) {
		t.Fatal("should authenticate with matching baseKey")
	}

	req2 := httptest.NewRequest("GET", "/", nil)
	req2.AddCookie(cookie)
	if IsAuthenticated(req2, username, password, "wrong-key") {
		t.Fatal("should not authenticate with different baseKey")
	}

	req3 := httptest.NewRequest("GET", "/", nil)
	req3.AddCookie(cookie)
	if IsAuthenticated(req3, username, password, "") {
		t.Fatal("should not authenticate when baseKey was used but not provided")
	}
}

func TestAuthenticate_SecureCookieFlag(t *testing.T) {
	rec1 := httptest.NewRecorder()
	Authenticate(rec1, "admin", "pass", "", "", true)
	if !rec1.Result().Cookies()[0].Secure {
		t.Fatal("cookie should be secure when secureCookie=true")
	}

	rec2 := httptest.NewRecorder()
	Authenticate(rec2, "admin", "pass", "", "", false)
	if rec2.Result().Cookies()[0].Secure {
		t.Fatal("cookie should not be secure when secureCookie=false")
	}
}
