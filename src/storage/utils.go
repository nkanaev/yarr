package storage

import (
	"strings"
	"golang.org/x/net/html"
)

func HTMLText(s string) string {
	tokenizer := html.NewTokenizer(strings.NewReader(s))
	contents := make([]string, 0)
	for {
		token := tokenizer.Next()
		if token == html.ErrorToken {
			break
		}
		if token == html.TextToken {
			content := strings.TrimSpace(html.UnescapeString(string(tokenizer.Text())))
			if len(content) > 0 {
				contents = append(contents, content)
			}
		}
	}
	return strings.Join(contents, " ")
}
