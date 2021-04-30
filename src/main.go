package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nkanaev/yarr/src/env"
	"github.com/nkanaev/yarr/src/platform"
	"github.com/nkanaev/yarr/src/server"
	"github.com/nkanaev/yarr/src/storage"
)

const APP = "yarr"

var Version string = "0.0"
var GitHash string = "unknown"

type Config struct {
	Address     string `json:"address"`
	Database    string `json:"database"`
	AuthFile    string `json:"auth-file"`
	CertFile    string `json:"cert-file"`
	KeyFile     string `json:"key-file"`
	BasePath    string `json:"base-path"`
	LogPath     string `json:"log-path"`
	OpenBrowser bool   `json:"open-browser"`
}

func initConfig(appPath string) *Config {
	config := &Config{
		Address:  "127.0.0.1:7070",
		Database: filepath.Join(appPath, "storage.db"),
	}
	appConfigPath := filepath.Join(appPath, "yarr.json")
	if _, err := os.Stat(appConfigPath); err == nil {
		// log.Printf("config path: %s", appConfigPath)
		f, err := os.Open(appConfigPath)
		if err != nil {
			log.Fatal("Failed to open config file: ", err)
		}
		defer f.Close()
		body, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal("Failed to read config file: ", err)
		}

		json.Unmarshal(body, config)
	}
	return config
}

func main() {
	configPath, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("Failed to get config dir: ", err)
	}
	appPath := filepath.Join(configPath, APP)
	if err := os.MkdirAll(appPath, 0755); err != nil {
		log.Fatal("Failed to create app config dir: ", err)
	}
	config := initConfig(appPath)
	if err := env.Fill(APP, config); err != nil {
		log.Fatal("Failed to fill env: ", err)
	}

	var addr, db, authfile, certfile, keyfile, basepath, logpath string
	var ver, open bool
	flag.StringVar(&addr, "addr", config.Address, "address to run server on")
	flag.StringVar(&authfile, "auth-file", config.AuthFile, "path to a file containing username:password")
	flag.StringVar(&basepath, "base", config.BasePath, "base path of the service url")
	flag.StringVar(&certfile, "cert-file", config.CertFile, "path to cert file for https")
	flag.StringVar(&keyfile, "key-file", config.KeyFile, "path to key file for https")
	flag.StringVar(&logpath, "log-path", config.LogPath, "server log path")
	flag.StringVar(&db, "db", config.Database, "storage file path")
	flag.BoolVar(&open, "open", config.OpenBrowser, "open the server in browser")
	flag.BoolVar(&ver, "version", false, "print application version")
	flag.Parse()

	if ver {
		fmt.Printf("v%s (%s)\n", Version, GitHash)
		return
	}

	if logpath != "" {
		f, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(f)
	} else {
		log.SetOutput(os.Stdout)
	}

	log.Printf("config %+v", config)

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
