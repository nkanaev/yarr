package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
)

func userIsAuthenticated(req *http.Request, username, password string) bool {
	cookie, _ := req.Cookie("auth")
	if cookie == nil {
		return false
	}
	parts := strings.Split(cookie.Value, ":")
	if len(parts) != 2 || !stringsEqual(parts[0], username) {
		return false
	}
	return stringsEqual(parts[1], secret(username, password))
}

func userAuthenticate(rw http.ResponseWriter, username, password string) {
	expires := time.Now().Add(time.Hour * 24 * 7) // 1 week
	var cookiePath string
	if BasePath != "" {
		cookiePath = BasePath
	} else {
		cookiePath = "/"
	}
	cookie := http.Cookie{
		Name:    "auth",
		Value:   username + ":" + secret(username, password),
		Expires: expires,
		Path:    cookiePath,
	}
	http.SetCookie(rw, &cookie)
}

func stringsEqual(p1, p2 string) bool {
	return subtle.ConstantTimeCompare([]byte(p1), []byte(p2)) == 1
}

func secret(msg, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(msg))
	src := mac.Sum(nil)
	return hex.EncodeToString(src)
}
