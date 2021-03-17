package server

import (
	"net/http"
	"strings"

	"github.com/nkanaev/yarr/src/assets"
	"github.com/nkanaev/yarr/src/auth"
	"github.com/nkanaev/yarr/src/router"
)

type authMiddleware struct {
	username string
	password string
	basepath string
	public   string
}

func unsafeMethod(method string) bool {
	return method == "POST" || method == "PUT" || method == "DELETE"
}

func (m *authMiddleware) handler(c *router.Context) {
	if strings.HasPrefix(c.Req.URL.Path, m.basepath + m.public) {
		c.Next()
		return
	}
	if auth.IsAuthenticated(c.Req, m.username, m.password) {
		c.Next()
		return
	}

	rootUrl := m.basepath + "/"

	if c.Req.URL.Path != rootUrl {
		c.Out.WriteHeader(http.StatusUnauthorized)
		return
	}

	if c.Req.Method == "POST" {
		username := c.Req.FormValue("username")
		password := c.Req.FormValue("password")
		if auth.StringsEqual(username, m.username) && auth.StringsEqual(password, m.password) {
			auth.Authenticate(c.Out, m.username, m.password, m.basepath)
			c.Redirect(rootUrl)
			return
		} else {
			c.HTML(http.StatusOK, assets.Template("login.html"), map[string]string{
				"username": username,
				"error": "Invalid username/password",
			})
			return
		}
	}
	c.HTML(http.StatusOK, assets.Template("login.html"), nil)
}
