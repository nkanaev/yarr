package server

import (
	"net/http"
	"crypto/subtle"
	"time"
)


func userIsAuthenticated(req *http.Request, username, password string) bool {
	cookie, _ :=  req.Cookie("auth")
	if cookie == nil {
		return false
	}
	// TODO: change to something sane
	if cookie.Value != username {
		return false
	}
	return true
}

func userAuthenticate(rw http.ResponseWriter, username, password string) {
	expires := time.Now().Add(time.Hour * 24 * 7)  // 1 week
	cookie := http.Cookie{Name: "auth", Value: username, Expires: expires}
	http.SetCookie(rw, &cookie)
}

func safeCompare(p1, p2 string) bool {
	return subtle.ConstantTimeCompare([]byte(p1), []byte(p2)) == 1
}
