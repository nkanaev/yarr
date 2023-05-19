package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nkanaev/yarr/src/platform"
	"github.com/nkanaev/yarr/src/server"
	"github.com/nkanaev/yarr/src/storage"
)

var Version string = "0.0"
var GitHash string = "unknown"

var OptList = make([]string, 0)

func opt(envVar, defaultValue string) string {
	OptList = append(OptList, envVar)
	value := os.Getenv(envVar)
	if value != "" {
		return value
	}
	return defaultValue
}

func parseAuthfile(authfile io.Reader) (username, password string, err error) {
	scanner := bufio.NewScanner(authfile)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return "", "", fmt.Errorf("wrong syntax (expected `username:password`)")
		}
		username = parts[0]
		password = parts[1]
		break
	}
	return username, password, nil
}

func main() {
	platform.FixConsoleIfNeeded()

	var addr, db, authfile, auth, certfile, keyfile, basepath, logfile, title string
	var ver, open bool

	flag.CommandLine.SetOutput(os.Stdout)

	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(out, "\nThe environmental variables, if present, will be used to provide\nthe default values for the params above:")
		fmt.Fprintln(out, " ", strings.Join(OptList, ", "))
	}

	flag.StringVar(&addr, "addr", opt("YARR_ADDR", "127.0.0.1:7070"), "address to run server on")
	flag.StringVar(&basepath, "base", opt("YARR_BASE", ""), "base path of the service url")
	flag.StringVar(&authfile, "auth-file", opt("YARR_AUTHFILE", ""), "`path` to a file containing username:password. Takes precedence over --auth (or YARR_AUTH)")
	flag.StringVar(&auth, "auth", opt("YARR_AUTH", ""), "string with username and password in the format `username:password`")
	flag.StringVar(&certfile, "cert-file", opt("YARR_CERTFILE", ""), "`path` to cert file for https")
	flag.StringVar(&keyfile, "key-file", opt("YARR_KEYFILE", ""), "`path` to key file for https")
	flag.StringVar(&db, "db", opt("YARR_DB", ""), "storage file `path`")
	flag.StringVar(&logfile, "log-file", opt("YARR_LOGFILE", ""), "`path` to log file to use instead of stdout")
	flag.StringVar(&title, "title", opt("YARR_TITLE", "Yarr!"), "title of the served page")
	flag.BoolVar(&ver, "version", false, "print application version")
	flag.BoolVar(&open, "open", false, "open the server in browser")
	flag.Parse()

	// Sanitize title as its used in the template.
	title = template.HTMLEscapeString(title)

	if ver {
		fmt.Printf("v%s (%s)\n", Version, GitHash)
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

	configPath, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("Failed to get config dir: ", err)
	}

	if db == "" {
		storagePath := filepath.Join(configPath, "yarr")
		if err := os.MkdirAll(storagePath, 0755); err != nil {
			log.Fatal("Failed to create app config dir: ", err)
		}
		db = filepath.Join(storagePath, "storage.db")
	}

	log.Printf("using db file %s", db)

	var username, password string
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

	srv := server.NewServer(store, addr, title)

	if basepath != "" {
		srv.BasePath = "/" + strings.Trim(basepath, "/")
	}

	if certfile != "" && keyfile != "" {
		srv.CertFile = certfile
		srv.KeyFile = keyfile
	}

	if username != "" && password != "" {
		srv.Username = username
		srv.Password = password
	}

	log.Printf("starting server at %s", srv.GetAddr())
	if open {
		platform.Open(srv.GetAddr())
	}
	platform.Start(srv)
}
