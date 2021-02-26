package server

import (
	"net/http"
	"regexp"
)

var BasePath string = ""

type Route struct {
	url        string
	urlRegex   *regexp.Regexp
	handler    func(http.ResponseWriter, *http.Request)
	manualAuth bool
}

func (r Route) ManualAuth() Route {
	r.manualAuth = true
	return r
}

func p(path string, handler http.HandlerFunc) Route {
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

func getRoute(reqPath string) (*Route, map[string]string) {
	vars := make(map[string]string)
	for _, route := range routes {
		if route.urlRegex.MatchString(reqPath) {
			matches := route.urlRegex.FindStringSubmatchIndex(reqPath)
			for i, key := range route.urlRegex.SubexpNames()[1:] {
				vars[key] = reqPath[matches[i*2+2]:matches[i*2+3]]
			}
			return &route, vars
		}
	}
	return nil, nil
}
