package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nkanaev/yarr/src/parser"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: <script> [url|filepath]")
		return
	}
	url := os.Args[1]
	var r io.Reader

	if strings.HasPrefix(url, "http") {
		res, err := http.Get(url)
		if err != nil {
			log.Fatalf("failed to get url %s: %s", url, err)
		}
		r = res.Body
	} else {
		var err error
		r, err = os.Open(url)
		if err != nil {
			log.Fatalf("failed to open file: %s", err)
		}
	}
	feed, err := parser.Parse(r)
	if err != nil {
		log.Fatalf("failed to parse feed: %s", err)
	}
	body, err := json.MarshalIndent(feed, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshall feed: %s", err)
	}
	fmt.Println(string(body))
}
