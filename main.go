package main

import (
	"github.com/nkanaev/yarr/server"
	"github.com/nkanaev/yarr/storage"
	"github.com/shibukawa/configdir"
	"log"
	"os"
	"path/filepath"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	configDirs := configdir.New("", "yarr")
	storageDir := configDirs.QueryFolders(configdir.Global)[0].Path
	storageFile := filepath.Join(storageDir, "storage.db")

	if err := os.MkdirAll(storageDir, 0755); err != nil {
		logger.Fatal(err)
	}

	db, err := storage.New(storageFile, logger)
	if err != nil {
		logger.Fatal(err)
	}

	srv := server.New(db, logger)
	srv.Start("127.0.0.1:8000")
}
