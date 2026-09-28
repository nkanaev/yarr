package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nkanaev/yarr/src/assets"
	"github.com/nkanaev/yarr/src/platform"
	"github.com/nkanaev/yarr/src/server"
	"github.com/nkanaev/yarr/src/storage"
	"github.com/nkanaev/yarr/src/worker"
)

var Version string = "0.0"
var GitHash string = "unknown"

func opt(envVar, defaultValue string) string {
	value := os.Getenv(envVar)
	if value != "" {
		return value
	}
	return defaultValue
}

func parseAuthfile(authfile io.Reader) (username, password string, err error) {
	if scanner := bufio.NewScanner(authfile); scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return "", "", fmt.Errorf("wrong syntax (expected `username:password`)")
		}
		username = strings.TrimSpace(parts[0])
		password = strings.TrimSpace(parts[1])

		if username == "" || password == "" {
			return "", "", fmt.Errorf("missing username and/or password")
		}
	} else {
		return "", "", fmt.Errorf("empty auth line")
	}
	return username, password, nil
}

func main() {
	platform.FixConsoleIfNeeded()

	var addr, db, authfile, auth, certfile, keyfile, basepath, logfile string
	var ver, open bool

	flag.CommandLine.SetOutput(os.Stdout)
	flag.StringVar(
		&addr,
		"addr",
		opt("YARR_ADDR", "127.0.0.1:7070"),
		"listen `address` (host:port or unix:/path/socket)\n(env: YARR_ADDR)",
	)
	flag.StringVar(
		&basepath,
		"base",
		opt("YARR_BASE", ""),
		"`path` prefix for the service url, e.g. /news\n(env: YARR_BASE)",
	)
	flag.StringVar(
		&authfile,
		"auth-file",
		opt("YARR_AUTHFILE", ""),
		"`path` to a file with username:password,\ntakes precedence over -auth\n(env: YARR_AUTHFILE)",
	)
	flag.StringVar(
		&auth,
		"auth",
		opt("YARR_AUTH", ""),
		"credentials as `username:password`\n(env: YARR_AUTH)",
	)
	flag.StringVar(
		&certfile,
		"cert-file",
		opt("YARR_CERTFILE", ""),
		"`path` to the TLS certificate file, requires -key-file\n(env: YARR_CERTFILE)",
	)
	flag.StringVar(
		&keyfile,
		"key-file",
		opt("YARR_KEYFILE", ""),
		"`path` to the TLS private key file, requires -cert-file\n(env: YARR_KEYFILE)",
	)
	flag.StringVar(&db, "db", opt("YARR_DB", ""), "storage database `path`\n(env: YARR_DB)")
	flag.StringVar(
		&logfile,
		"log-file",
		opt("YARR_LOGFILE", ""),
		"`path` to the log file, defaults to stdout\n(env: YARR_LOGFILE)",
	)
	flag.BoolVar(&ver, "version", false, "print version and exit")
	flag.BoolVar(&open, "open", false, "open the server url in the default browser")
	flag.Parse()

	if ver {
		fmt.Printf("%s (%s)\n", Version, GitHash)
		return
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	if logfile != "" {
		file, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal("Failed to setup log file: ", err)
		}
		defer file.Close()
		log.SetOutput(file)
	} else {
		log.SetOutput(os.Stdout)
	}

	if open && strings.HasPrefix(addr, "unix:") {
		log.Fatal("Cannot open ", addr, " in browser")
	}

	if db == "" {
		configPath, err := os.UserConfigDir()
		if err != nil {
			log.Fatal("Failed to get config dir: ", err)
		}

		storagePath := filepath.Join(configPath, "yarr")
		if err := os.MkdirAll(storagePath, 0755); err != nil {
			log.Fatal("Failed to create app config dir: ", err)
		}
		db = filepath.Join(storagePath, "storage.db")
	}

	log.Printf("using db file %s", db)

	var username, password string
	var err error
	if authfile != "" {
		f, err := os.Open(authfile)
		if err != nil {
			log.Fatal("Failed to open auth file: ", err)
		}
		defer f.Close()
		username, password, err = parseAuthfile(f)
		if err != nil {
			log.Fatal("Failed to parse auth file: ", err)
		}
	} else if auth != "" {
		username, password, err = parseAuthfile(strings.NewReader(auth))
		if err != nil {
			log.Fatal("Failed to parse auth literal: ", err)
		}
	}

	if (certfile != "" || keyfile != "") && (certfile == "" || keyfile == "") {
		log.Fatalf("Both cert & key files are required")
	}

	store, err := storage.New(db)
	if err != nil {
		log.Fatal("Failed to initialise database: ", err)
	}

	worker.SetVersion(Version)
	srv := server.NewServer(addr)
	srv.StaticFS = assets.StaticFS()
	srv.Template = assets.Templates()

	if basepath != "" {
		srv.BasePath = "/" + strings.Trim(basepath, "/")
	}

	if certfile != "" && keyfile != "" {
		srv.CertFile = certfile
		srv.KeyFile = keyfile
	}

	wrk := worker.NewWorker(store)

	if username != "" && password != "" {
		srv.Auth = server.NewLocalAuthProvider(username, password, basepath)
	}
	srv.Storage = server.NewLocalStorage(store)
	srv.Scheduler = wrk

	log.Printf("starting server at %s", srv.GetAddr())
	if open {
		platform.Open(srv.GetAddr())
	}
	wrk.StartFeedCleaner()
	wrk.SetRefreshRate(store.GetSettings().RefreshRate)
	platform.Start(srv)
}
