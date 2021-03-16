package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
)

func IsAuthenticated(req *http.Request, username, password string) bool {
	cookie, _ := req.Cookie("auth")
	if cookie == nil {
		return false
	}
	parts := strings.Split(cookie.Value, ":")
	if len(parts) != 2 || !StringsEqual(parts[0], username) {
		return false
	}
	return StringsEqual(parts[1], secret(username, password))
}

func Authenticate(rw http.ResponseWriter, username, password, basepath string) {
	expires := time.Now().Add(time.Hour * 24 * 7) // 1 week

	var cookiePath string
	if basepath != "" {
		cookiePath = basepath
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

func Logout(rw http.ResponseWriter) {
	cookie := http.Cookie{
		Name:    "auth",
		Value:   "",
		MaxAge:  -1,
	}
	http.SetCookie(rw, &cookie)
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
