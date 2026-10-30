package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	authCookieName = "auth"
	authMaxAge     = 7 * 24 * time.Hour
)

func parseAuthCookie(req *http.Request, username, password string) (time.Time, bool) {
	cookie, _ := req.Cookie(authCookieName)
	if cookie == nil {
		return time.Time{}, false
	}
	parts := strings.Split(cookie.Value, ":")
	if len(parts) != 3 || !StringsEqual(parts[0], username) {
		return time.Time{}, false
	}
	if !StringsEqual(parts[2], secret(parts[0]+":"+parts[1], password)) {
		return time.Time{}, false
	}
	issuedUnix, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return time.Time{}, false
	}
	issued := time.Unix(issuedUnix, 0)
	if time.Since(issued) > authMaxAge {
		return time.Time{}, false
	}
	return issued, true
}

func newAuthCookie(username, password, basepath string, issued time.Time) *http.Cookie {
	ts := strconv.FormatInt(issued.Unix(), 10)
	return &http.Cookie{
		Name:     authCookieName,
		Value:    username + ":" + ts + ":" + secret(username+":"+ts, password),
		MaxAge:   int(authMaxAge / time.Second),
		Path:     basepath,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}
}

func IsAuthenticated(req *http.Request, username, password string) bool {
	_, ok := parseAuthCookie(req, username, password)
	return ok
}

func Authenticate(rw http.ResponseWriter, username, password, basepath string) {
	http.SetCookie(rw, newAuthCookie(username, password, basepath, time.Now()))
}

func Logout(rw http.ResponseWriter, basepath string) {
	http.SetCookie(rw, &http.Cookie{
		Name:     authCookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     basepath,
		HttpOnly: true,
	})
}

func StringsEqual(p1, p2 string) bool {
	return subtle.ConstantTimeCompare([]byte(p1), []byte(p2)) == 1
}

func secret(msg, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(msg))
	src := mac.Sum(nil)
	return hex.EncodeToString(src)
}

func LocalAuth(next http.Handler, username, password, basepath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		issued, ok := parseAuthCookie(r, username, password)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if time.Since(issued) > authMaxAge/2 {
			http.SetCookie(w, newAuthCookie(username, password, basepath, time.Now()))
		}
		next.ServeHTTP(w, r)
	}
}
