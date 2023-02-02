package silo

import (
	"net/url"
	"strings"
)

func RedirectURL(link string) string {
	if strings.HasPrefix(link, "https://www.google.com/url?") {
		if u, err := url.Parse(link); err == nil {
			if u2 := u.Query().Get("url"); u2 != "" {
				return u2
			}
		}
	}
	return link
}
