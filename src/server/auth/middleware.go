package auth

import (
	"net/http"
	"strings"

	"github.com/nkanaev/yarr/src/server/router"
	"github.com/nkanaev/yarr/src/storage"
)

type Middleware struct {
	Username string
	Password string
	BasePath string
	Public   []string
	DB       storage.Storage
}

func (m *Middleware) Handler(c *router.Context) {
	for _, path := range m.Public {
		if strings.HasPrefix(c.Req.URL.Path, m.BasePath+path) {
			c.Next()
			return
		}
	}
	if IsAuthenticated(c.Req, m.Username, m.Password) {
		c.Next()
	} else {
		c.Out.WriteHeader(http.StatusUnauthorized)
	}
}
