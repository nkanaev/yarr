package server

import (
	"log"
	"net/http"

	"github.com/nkanaev/yarr/src/storage"
	"github.com/nkanaev/yarr/src/worker"
)

var BasePath string = ""

type Server struct {
	Addr        string
	db          *storage.Storage
	worker      *worker.Worker
	// auth
	Username string
	Password string
	// https
	CertFile string
	KeyFile  string
}

func NewServer(db *storage.Storage, addr string) *Server {
	return &Server{
		db:     db,
		Addr:   addr,
		worker: worker.NewWorker(db),
	}
}

func (h *Server) GetAddr() string {
	proto := "http"
	if h.CertFile != "" && h.KeyFile != "" {
		proto = "https"
	}
	return proto + "://" + h.Addr + BasePath
}

func (s *Server) Start() {
	s.worker.Start()

	httpserver := &http.Server{Addr: s.Addr, Handler: s.handler()}

	var err error
	if s.CertFile != "" && s.KeyFile != "" {
		err = httpserver.ListenAndServeTLS(s.CertFile, s.KeyFile)
	} else {
		err = httpserver.ListenAndServe()
	}
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}


/*
func (h Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	reqPath := req.URL.Path
	if BasePath != "" {
		if !strings.HasPrefix(reqPath, BasePath) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		reqPath = strings.TrimPrefix(req.URL.Path, BasePath)
		if reqPath == "" {
			http.Redirect(rw, req, BasePath+"/", http.StatusFound)
			return
		}
	}
	route, vars := getRoute(reqPath)
	if route == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	if h.requiresAuth() && !route.manualAuth {
		if unsafeMethod(req.Method) && req.Header.Get("X-Requested-By") != "yarr" {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !userIsAuthenticated(req, h.Username, h.Password) {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	ctx := context.WithValue(req.Context(), ctxHandler, &h)
	ctx = context.WithValue(ctx, ctxVars, vars)
	route.handler(rw, req.WithContext(ctx))
}
*/

func (h Server) requiresAuth() bool {
	return h.Username != "" && h.Password != ""
}
