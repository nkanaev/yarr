package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"strings"
)

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

func Logout(rw http.ResponseWriter, basepath string) {
	http.SetCookie(rw, &http.Cookie{
		Name:   "auth",
		Value:  "",
		MaxAge: -1,
		Path:   basepath,
	})
}

func StringsEqual(p1, p2 string) bool {
	return subtle.ConstantTimeCompare([]byte(p1), []byte(p2)) == 1
}

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
