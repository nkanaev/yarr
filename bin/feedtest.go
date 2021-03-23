package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nkanaev/yarr/src/parser"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: <script> url")
		return
	}
	url := os.Args[1]
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("failed to get url %s: %s", url, err)
	}
	feed, err := parser.Parse(res.Body)
	if err != nil {
		log.Fatalf("failed to parse feed: %s", err)
	}
	body, err := json.MarshalIndent(feed, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshall feed: %s", err)
	}
	fmt.Println(string(body))
}
