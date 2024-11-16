package nostr

import (
	"net/url"
	"strings"
)

// NormalizeURL normalizes the url and replaces http://, https:// schemes with ws://, wss://
// and normalizes the path.
func NormalizeURL(u string) string {
	if u == "" {
		return ""
	}

	u = strings.TrimSpace(u)
	u = strings.ToLower(u)

	if fqn := strings.Split(u, ":")[0]; fqn == "localhost" || fqn == "127.0.0.1" {
		u = "ws://" + u
	} else if !strings.HasPrefix(u, "http") && !strings.HasPrefix(u, "ws") {
		u = "wss://" + u
	}

	p, err := url.Parse(u)
	if err != nil {
		return ""
	}

	if p.Scheme == "http" {
		p.Scheme = "ws"
	} else if p.Scheme == "https" {
		p.Scheme = "wss"
	}

	p.Path = strings.TrimRight(p.Path, "/")

	return p.String()
}

// NormalizeOKMessage takes a string message that is to be sent in an `OK` or `CLOSED` command
// and prefixes it with "<prefix>: " if it doesn't already have an acceptable prefix.
func NormalizeOKMessage(reason string, prefix string) string {
	if idx := strings.Index(reason, ": "); idx == -1 || strings.IndexByte(reason[0:idx], ' ') != -1 {
		return prefix + ": " + reason
	}
	return reason
}
