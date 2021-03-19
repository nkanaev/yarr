package server

import (
	"log"
	"net/http"

	"github.com/nkanaev/yarr/src/storage"
	"github.com/nkanaev/yarr/src/worker"
)

var BasePath string = ""

type Server struct {
	Addr   string
	db     *storage.Storage
	worker *worker.Worker
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
