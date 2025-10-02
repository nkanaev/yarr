package server

import (
	"net"
	"net/url"
	"strings"
)

func isInternalFromURL(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	host := parsedURL.Host

	// Handle "host:port" format
	if strings.Contains(host, ":") {
		host, _, err = net.SplitHostPort(host)
		if err != nil {
			return false
		}
	}

	if host == "localhost" {
		return true
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}

	return ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast()
}
