package router

import (
	"net/http"
	"regexp"
	"strings"
)

type Handler func(*Context)

type Router struct {
	middle []Handler
	routes []Route
	base   string
}

type Route struct {
	regex *regexp.Regexp
	chain []Handler
}

func NewRouter(base string) *Router {
	router := &Router{}
	router.middle = make([]Handler, 0)
	router.routes = make([]Route, 0)
	router.base = base
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

func (r *Router) resolve(path string) *Route {
	for _, route := range r.routes {
		if route.regex.MatchString(path) {
			return &route
		}
	}
	return nil
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// autoclose open base url
	if r.base != "" {
		if r.base == req.URL.Path {
			http.Redirect(rw, req, r.base+"/", http.StatusFound)
			return
		}
		if !strings.HasPrefix(req.URL.Path, r.base) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
	}

	path := strings.TrimPrefix(req.URL.Path, r.base)

	route := r.resolve(path)
	if route == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	context := &Context{}
	context.Req = req
	context.Out = rw
	context.Vars = regexGroups(path, route.regex)
	context.index = -1
	context.chain = route.chain
	context.Next()
}
