package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/nkanaev/yarr/platform"
	"github.com/nkanaev/yarr/server"
	"github.com/nkanaev/yarr/storage"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var Version string = "0.0"
var GitHash string = "unknown"

func main() {
	var addr, db, authfile, certfile, keyfile string
	var ver bool
	flag.StringVar(&addr, "addr", "127.0.0.1:7070", "address to run server on")
	flag.StringVar(&authfile, "auth-file", "", "path to a file containing username:password")
	flag.StringVar(&certfile, "cert-file", "", "path to cert file for https")
	flag.StringVar(&keyfile, "key-file", "", "path to key file for https")
	flag.StringVar(&db, "db", "", "storage file path")
	flag.BoolVar(&ver, "version", false, "print application version")
	flag.Parse()

	if ver {
		fmt.Printf("v%s (%s)\n", Version, GitHash)
		return
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	configPath, err := os.UserConfigDir()
	if err != nil {
		logger.Fatal("Failed to get config dir: ", err)
	}

	if db == "" {
		storagePath := filepath.Join(configPath, "yarr")
		if err := os.MkdirAll(storagePath, 0755); err != nil {
			logger.Fatal("Failed to create app config dir: ", err)
		}
		db = filepath.Join(storagePath, "storage.db")
	}

	logger.Printf("using db file %s", db)

	var username, password string
	if authfile != "" {
		f, err := os.Open(authfile)
		if err != nil {
			logger.Fatal("Failed to open auth file: ", err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				logger.Fatalf("Invalid auth: %v (expected `username:password`)", line)
			}
			username = parts[0]
			password = parts[1]
			break
		}
	}

	if (certfile != "" || keyfile != "") && (certfile == "" || keyfile == "") {
		logger.Fatalf("Both cert & key files are required")
	}

	store, err := storage.New(db, logger)
	if err != nil {
		logger.Fatal("Failed to initialise database: ", err)
	}

	srv := server.New(store, logger, addr)

	proto := "http"
	if certfile != "" && keyfile != "" {
		srv.CertFile = certfile
		srv.KeyFile = keyfile
		proto = "https"
	}

	if username != "" && password != "" {
		srv.Username = username
		srv.Password = password
	}

	logger.Printf("starting server at %s://%s", proto, addr)
	platform.Start(srv)
}
