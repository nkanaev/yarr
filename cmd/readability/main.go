package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nkanaev/yarr/src/content/readability"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: <script> [url]")
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

	content, err := readability.ExtractContent(r)
	if err != nil {
		log.Fatalf("failed to extract content: %s", err)
	}
	fmt.Println(content)
}
