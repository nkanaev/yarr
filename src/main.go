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

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate|log.Ltime|log.Lshortfile)

	var addr, db, authfile, certfile, keyfile string
	var ver, open bool
	flag.StringVar(&addr, "addr", "127.0.0.1:7070", "address to run server on")
	flag.StringVar(&authfile, "auth-file", "", "path to a file containing username:password")
	flag.StringVar(&server.BasePath, "base", "", "base path of the service url")
	flag.StringVar(&certfile, "cert-file", "", "path to cert file for https")
	flag.StringVar(&keyfile, "key-file", "", "path to key file for https")
	flag.StringVar(&db, "db", "", "storage file path")
	flag.BoolVar(&ver, "version", false, "print application version")
	flag.BoolVar(&open, "open", false, "open the server in browser")
	flag.Parse()

	if ver {
		fmt.Printf("v%s (%s)\n", Version, GitHash)
		return
	}

	if server.BasePath != "" && !strings.HasPrefix(server.BasePath, "/") {
		server.BasePath = "/" + server.BasePath
	}

	if server.BasePath != "" && strings.HasSuffix(server.BasePath, "/") {
		server.BasePath = strings.TrimSuffix(server.BasePath, "/")
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

	srv := server.New(store, addr)

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
