package server

import (
	"html/template"
	"io/fs"
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/nkanaev/yarr/src/storage"
	"github.com/nkanaev/yarr/src/storage/model"
)

type FeedScheduler interface {
	FeedsPending() int32
	RefreshFeeds()
	SetRefreshRate(minutes int64)
}

type FeedLink struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	// TODO: bad idea, remove
	TitleOverride string `json:"title_override,omitempty"`
}

type AddFeedResult struct {
	Feed    *model.Feed
	Choices []FeedLink
}

type FeedIngestor interface {
	AddFeed(ctx context.Context, store storage.Storage, feed FeedLink) (AddFeedResult, error)
	AddOPML(ctx context.Context, store storage.Storage, file io.Reader) error
}

type Server struct {
	Addr     string
	BasePath string

	Storage   StorageProvider
	Scheduler FeedScheduler
	Auth      AuthProvider
	Ingestor  FeedIngestor

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
