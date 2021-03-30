package scraper

import (
	"net/url"
)

func any(els []string, el string, match func(string, string) bool) bool {
	for _, x := range els {
		if match(x, el) {
			return true
		}
	}
	return false
}

func absoluteUrl(href, base string) string {
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	hrefUrl, err := url.Parse(href)
	if err != nil {
		return ""
	}
	return baseUrl.ResolveReference(hrefUrl).String()
}

func urlDomain(val string) string {
	if u, err := url.Parse(val); err == nil {
		return u.Host
	}
	return val
}
