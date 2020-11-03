package server

import (
	"net/http"
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

}
