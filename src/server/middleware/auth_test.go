package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func TestIsAuthenticated(t *testing.T) {
	username, password := "user", "pass"
	cookie := func(value string) *http.Request {
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: "auth", Value: value})
		return req
	}

	tests := []struct {
		name string
		req  *http.Request
		want bool
	}{
		{"no cookie", httptest.NewRequest("GET", "/", nil), false},
		{"malformed cookie", cookie("invalid"), false},
		{"wrong username", cookie("other:" + secret(username, password)), false},
		{"wrong secret", cookie(username + ":invalid"), false},
		{"valid cookie", cookie(username + ":" + secret(username, password)), true},
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
	}), username, password)

	tests := []struct {
		name   string
		cookie string
		status int
	}{
		{"no cookie", "", http.StatusUnauthorized},
		{"wrong secret", username + ":invalid", http.StatusUnauthorized},
		{"valid cookie", username + ":" + secret(username, password), http.StatusOK},
	}
	for _, test := range tests {
		req := httptest.NewRequest("GET", "/", nil)
		if test.cookie != "" {
			req.AddCookie(&http.Cookie{Name: "auth", Value: test.cookie})
		}
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != test.status {
			t.Errorf("%s: status = %d; want %d", test.name, rec.Code, test.status)
		}
	}
}
