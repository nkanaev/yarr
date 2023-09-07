package htmlutil

import (
	"net/url"
	"strings"
)

func Any(els []string, el string, match func(string, string) bool) bool {
	for _, x := range els {
		if match(x, el) {
			return true
		}
	}
	return false
}

func AbsoluteUrl(href, base string) string {
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

func URLDomain(val string) string {
	if u, err := url.Parse(val); err == nil {
		return u.Host
	}
	return val
}

func IsAPossibleLink(val string) bool {
	return strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://")
}
