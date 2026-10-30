package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func authCookieValue(username, password string, issued time.Time) string {
	ts := strconv.FormatInt(issued.Unix(), 10)
	return username + ":" + ts + ":" + secret(username+":"+ts, password)
}

func TestStringsEqual(t *testing.T) {
	tests := []struct {
		p1, p2 string
		equal  bool
	}{
		{"secret", "secret", true},
		{"secret", "Secret", false},
		{"secret", "", false},
	}
	for _, test := range tests {
		if equal := StringsEqual(test.p1, test.p2); equal != test.equal {
			t.Errorf("StringsEqual(%q, %q) = %v; want %v", test.p1, test.p2, equal, test.equal)
		}
	}
}

func TestNewAuthCookie(t *testing.T) {
	username, password, basepath := "user", "pass", "/base"
	issued := time.Unix(time.Now().Unix(), 0)
	cookie := newAuthCookie(username, password, basepath, issued)
	if cookie.Name != authCookieName {
		t.Errorf("name = %q; want %q", cookie.Name, authCookieName)
	}
	if cookie.Path != basepath {
		t.Errorf("path = %q; want %q", cookie.Path, basepath)
	}
	if cookie.MaxAge != int(authMaxAge/time.Second) {
		t.Errorf("max age = %d; want %d", cookie.MaxAge, int(authMaxAge/time.Second))
	}
	if !cookie.HttpOnly || !cookie.Secure {
		t.Errorf("HttpOnly = %v, Secure = %v; want true, true", cookie.HttpOnly, cookie.Secure)
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(cookie)
	got, ok := parseAuthCookie(req, username, password)
	if !ok {
		t.Fatal("parseAuthCookie() = false; want true")
	}
	if got.Unix() != issued.Unix() {
		t.Errorf("issued = %v; want %v", got, issued)
	}
}

func TestIsAuthenticated(t *testing.T) {
	username, password := "user", "pass"
	now := time.Now()
	cookie := func(value string) *http.Request {
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: authCookieName, Value: value})
		return req
	}

	tests := []struct {
		name string
		req  *http.Request
		want bool
	}{
		{"no cookie", httptest.NewRequest("GET", "/", nil), false},
		{"malformed cookie", cookie("invalid"), false},
		{"wrong username", cookie(authCookieValue("other", password, now)), false},
		{"wrong secret", cookie(username + ":" + strconv.FormatInt(now.Unix(), 10) + ":invalid"), false},
		{"invalid timestamp", cookie(username + ":invalid:" + secret(username+":invalid", password)), false},
		{"expired cookie", cookie(authCookieValue(username, password, now.Add(-authMaxAge-time.Minute))), false},
		{"valid cookie", cookie(authCookieValue(username, password, now)), true},
	}
	for _, test := range tests {
		if got := IsAuthenticated(test.req, username, password); got != test.want {
			t.Errorf("%s: IsAuthenticated = %v; want %v", test.name, got, test.want)
		}
	}
}

func TestLocalAuth(t *testing.T) {
	username, password := "user", "pass"
	handler := LocalAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), username, password, "")

	tests := []struct {
		name   string
		cookie string
		status int
	}{
		{"no cookie", "", http.StatusUnauthorized},
		{"wrong secret", username + ":" + strconv.FormatInt(time.Now().Unix(), 10) + ":invalid", http.StatusUnauthorized},
		{"expired cookie", authCookieValue(username, password, time.Now().Add(-authMaxAge-time.Minute)), http.StatusUnauthorized},
		{"valid cookie", authCookieValue(username, password, time.Now()), http.StatusOK},
	}
	for _, test := range tests {
		req := httptest.NewRequest("GET", "/", nil)
		if test.cookie != "" {
			req.AddCookie(&http.Cookie{Name: authCookieName, Value: test.cookie})
		}
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != test.status {
			t.Errorf("%s: status = %d; want %d", test.name, rec.Code, test.status)
		}
	}
}

func TestLocalAuthRefreshesCookie(t *testing.T) {
	username, password, basepath := "user", "pass", "/base"
	handler := LocalAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), username, password, basepath)

	tests := []struct {
		name        string
		issued      time.Time
		wantRefresh bool
	}{
		{"fresh cookie", time.Now(), false},
		{"stale cookie", time.Now().Add(-authMaxAge/2 - time.Minute), true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(&http.Cookie{Name: authCookieName, Value: authCookieValue(username, password, test.issued)})
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			cookies := rec.Result().Cookies()
			if got := len(cookies) > 0; got != test.wantRefresh {
				t.Fatalf("refreshed = %v; want %v", got, test.wantRefresh)
			}
			if !test.wantRefresh {
				return
			}
			if cookies[0].Path != basepath {
				t.Errorf("path = %q; want %q", cookies[0].Path, basepath)
			}
			req = httptest.NewRequest("GET", "/", nil)
			req.AddCookie(cookies[0])
			if _, ok := parseAuthCookie(req, username, password); !ok {
				t.Error("refreshed cookie is not valid")
			}
		})
	}
}
