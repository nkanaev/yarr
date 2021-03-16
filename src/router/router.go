package router

import (
	"regexp"
)

type Handler func(*Context)

type Router struct {
	middle []Handler
	routes []Route
}

type Route struct {
	regex *regexp.Regexp
	chain []Handler
}

func NewRouter() *Router {
	router := &Router{}
	router.middle = make([]Handler, 0)
	router.routes = make([]Route, 0)
	return router
}

func (r *Router) Use(h Handler) {
	r.middle = append(r.middle, h)
}

func (r *Router) For(path string, handler Handler) {
	x := Route{}
	x.regex = routeRegexp(path)
	x.chain = append(r.middle, handler)

	r.routes = append(r.routes, x)
}

func (r *Router) resolve(path string) *route {
	for _, r := range r.routes {
		if r.regex.MatchString(path) {
			return &r
		}
	}
	return nil
}

func regexGroups(input string, regex *regexp.Regexp) map[string]string {
	groups := make(map[string]string)
	matches := regex.FindStringSubmatchIndex(input)
	for i, key := range regex.SubexpNames()[1:] {
		groups[key] = input[matches[i*2+2]:matches[i*2+3]]
	}
	return groups
}

func routeRegexp(route string) *regexp.Regexp {
	chunks := regexp.MustCompile(`[\*\:]\w+`)
	output := chunks.ReplaceAllStringFunc(route, func(m string) string {
		if m[0:1] == `*` {
			return "(?P<" + m[1:] + ">.+)"
		}
		return "(?P<" + m[1:] + ">[^/]+)"
	})
	output = "^" + output + "$"
	return regexp.MustCompile(output)
}
