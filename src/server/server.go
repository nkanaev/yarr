package server

import (
	"html/template"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/nkanaev/yarr/src/storage"
)

type Server struct {
	Addr     string
	BasePath string

	Storage   StorageProvider
	Scheduler FeedScheduler
	Auth      AuthProvider

	StaticFS fs.FS
	Template *template.Template

	// https
	CertFile string
	KeyFile  string
}

func NewServer(addr string) *Server {
	return &Server{
		Addr: addr,
	}
}

func (h *Server) GetAddr() string {
	proto := "http"
	if h.CertFile != "" && h.KeyFile != "" {
		proto = "https"
	}
	return proto + "://" + h.Addr + h.BasePath
}

func (s *Server) db(r *http.Request) storage.Storage {
	return s.Storage.GetStorage(r)
}

func (s *Server) Start() {
	var ln net.Listener
	var err error

	if path, isUnix := strings.CutPrefix(s.Addr, "unix:"); isUnix {
		err = os.Remove(path)
		if err != nil {
			log.Print(err)
		}
		ln, err = net.Listen("unix", path)
	} else {
		ln, err = net.Listen("tcp", s.Addr)
	}

	if err != nil {
		log.Fatal(err)
	}

	httpserver := &http.Server{Handler: s.Handler()}
	if s.CertFile != "" && s.KeyFile != "" {
		err = httpserver.ServeTLS(ln, s.CertFile, s.KeyFile)
		ln.Close()
	} else {
		err = httpserver.Serve(ln)
	}

	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
