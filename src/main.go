package main

import (
	"bufio"
	"flag"
	"fmt"
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

func opt(envVar, defaultValue string) string {
	value := os.Getenv(envVar)
	if value != "" {
		return value
	}
	return defaultValue
}

func main() {
	var addr, db, authfile, certfile, keyfile, basepath, logfile string
	var ver, open bool
	flag.StringVar(&addr, "addr", opt("YARR_ADDR", "127.0.0.1:7070"), "address to run server on")
	flag.StringVar(&authfile, "auth-file", opt("YARR_AUTHFILE", ""), "path to a file containing username:password")
	flag.StringVar(&basepath, "base", opt("YARR_BASE", ""), "base path of the service url")
	flag.StringVar(&certfile, "cert-file", opt("YARR_CERTFILE", ""), "path to cert file for https")
	flag.StringVar(&keyfile, "key-file", opt("YARR_KEYFILE", ""), "path to key file for https")
	flag.StringVar(&db, "db", opt("YARR_DB", ""), "storage file path")
	flag.StringVar(&logfile, "log-file", opt("YARR_LOGFILE", ""), "path to log file to use instead of stdout")
	flag.BoolVar(&ver, "version", false, "print application version")
	flag.BoolVar(&open, "open", false, "open the server in browser")
	flag.Parse()

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
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				log.Fatalf("Invalid auth: %v (expected `username:password`)", line)
			}
			username = parts[0]
			password = parts[1]
			break
		}
	}

	if (certfile != "" || keyfile != "") && (certfile == "" || keyfile == "") {
		log.Fatalf("Both cert & key files are required")
	}

	store, err := storage.New(db)
	if err != nil {
		log.Fatal("Failed to initialise database: ", err)
	}

	srv := server.NewServer(store, addr)

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
