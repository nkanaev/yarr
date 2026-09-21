package server

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/nkanaev/yarr/src/server/middleware"
)

type AuthProvider interface {
	Middleware(next http.Handler) http.HandlerFunc
	IsAuthenticated(request *http.Request) bool
	Authenticate(rw http.ResponseWriter, username, password string) bool
	Logout(rw http.ResponseWriter)
	FeverAPIKey(request *http.Request) string
}

type localAuth struct {
	Username string
	Password string
	BasePath string
}

func NewLocalAuthProvider(username, password, basepath string) AuthProvider {
	return &localAuth{Username: username, Password: password, BasePath: basepath}
}

func (a *localAuth) Middleware(next http.Handler) http.HandlerFunc {
	return middleware.LocalAuth(next, a.Username, a.Password)
}

func (a *localAuth) IsAuthenticated(r *http.Request) bool {
	return middleware.IsAuthenticated(r, a.Username, a.Password)
}

func (a *localAuth) Authenticate(rw http.ResponseWriter, username, password string) bool {
	if !middleware.StringsEqual(username, a.Username) || !middleware.StringsEqual(password, a.Password) {
		return false
	}
	middleware.Authenticate(rw, a.Username, a.Password, a.BasePath)
	return true
}

func (a *localAuth) Logout(rw http.ResponseWriter) {
	middleware.Logout(rw, a.BasePath)
}

func (a *localAuth) FeverAPIKey(r *http.Request) string {
	md5HashValue := md5.Sum(fmt.Appendf(nil, "%s:%s", a.Username, a.Password))
	return fmt.Sprintf("%x", md5HashValue[:])
}
