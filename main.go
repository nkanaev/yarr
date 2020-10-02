package main

import (
	"flag"
	"fmt"
	"github.com/nkanaev/yarr/server"
	"github.com/nkanaev/yarr/storage"
	"github.com/nkanaev/yarr/platform"
	"log"
	"os"
	"path/filepath"
)

var Version string = "v0.0"
var GitHash string = "unknown"

func main() {
	var addr, storageFile string
	var ver bool
	flag.StringVar(&addr, "addr", "127.0.0.1:7070", "address to run server on")
	flag.StringVar(&storageFile, "db", "", "storage file path")
	flag.BoolVar(&ver, "version", false, "print application version")
	flag.Parse()

	if ver {
		fmt.Printf("%s (%s)\n", Version, GitHash)
		return
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	configPath, err := os.UserConfigDir()
	if err != nil {
		logger.Fatal("Failed to get config dir: ", err)
	}

	if storageFile == "" {
		storagePath := filepath.Join(configPath, "yarr")
		if err := os.MkdirAll(storagePath, 0755); err != nil {
			logger.Fatal("Failed to create app config dir: ", err)
		}
		storageFile = filepath.Join(storagePath, "storage.db")
	}

	db, err := storage.New(storageFile, logger)
	if err != nil {
		logger.Fatal("Failed to initialise database: ", err)
	}

	srv := server.New(db, logger, addr)
	logger.Printf("starting server at http://%s", addr)
	platform.Start(srv)
}
