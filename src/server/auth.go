package server

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/nkanaev/yarr/src/server/middleware"
)

type Auth interface {
	Middleware(next http.Handler) http.HandlerFunc
	IsAuthenticated(request *http.Request) bool
	Authenticate(rw http.ResponseWriter, username, password string) bool
	Logout(rw http.ResponseWriter)
	FeverAPIKey() string
}

type LocalAuth struct {
	Username string
	Password string
	BasePath string
}

func (a LocalAuth) Middleware(next http.Handler) http.HandlerFunc {
	return middleware.LocalAuth(next, a.Username, a.Password)
}

func (a LocalAuth) IsAuthenticated(r *http.Request) bool {
	return middleware.IsAuthenticated(r, a.Username, a.Password)
}

func (a LocalAuth) Authenticate(rw http.ResponseWriter, username, password string) bool {
	if !middleware.StringsEqual(username, a.Username) || !middleware.StringsEqual(password, a.Password) {
		return false
	}
	middleware.Authenticate(rw, a.Username, a.Password, a.BasePath)
	return true
}

func (a LocalAuth) Logout(rw http.ResponseWriter) {
	middleware.Logout(rw, a.BasePath)
}

func (a LocalAuth) FeverAPIKey() string {
	md5HashValue := md5.Sum(fmt.Appendf(nil, "%s:%s", a.Username, a.Password))
	return fmt.Sprintf("%x", md5HashValue[:])
}
