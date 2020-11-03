package server

import (
	"net/http"
	"regexp"
)

type Route struct {
	url      string
	urlRegex *regexp.Regexp
	handler  func(http.ResponseWriter, *http.Request)
	skipAuth bool
}

func (r Route) SkipAuth() Route {
	r.skipAuth = true
	return r
}

func p(path string, handler func(http.ResponseWriter, *http.Request)) Route {
	var urlRegexp string
	urlRegexp = regexp.MustCompile(`[\*\:]\w+`).ReplaceAllStringFunc(path, func(m string) string {
		if m[0:1] == `*` {
			return "(?P<" + m[1:] + ">.+)"
		}
		return "(?P<" + m[1:] + ">[^/]+)"
	})
	urlRegexp = "^" + urlRegexp + "$"
	return Route{
		url:      path,
		urlRegex: regexp.MustCompile(urlRegexp),
		handler:  handler,
	}
}

func getRoute(req *http.Request) (*Route, map[string]string) {
	vars := make(map[string]string)
	for _, route := range routes {
		if route.urlRegex.MatchString(req.URL.Path) {
			matches := route.urlRegex.FindStringSubmatchIndex(req.URL.Path)
			for i, key := range route.urlRegex.SubexpNames()[1:] {
				vars[key] = req.URL.Path[matches[i*2+2]:matches[i*2+3]]
			}
			return &route, vars
		}
	}
	return nil, nil
}
