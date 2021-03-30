package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nkanaev/yarr/src/reader"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: <script> [url]")
		return
	}
	url := os.Args[1]
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("failed to get url %s: %s", url, err)
	}
	defer res.Body.Close()

	content, err := reader.ExtractContent(res.Body)
	if err != nil {
		log.Fatalf("failed to extract content: %s", err)
	}
	fmt.Println(content)
}
