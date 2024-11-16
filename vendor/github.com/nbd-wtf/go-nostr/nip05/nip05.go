package nip05

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/nbd-wtf/go-nostr"
)

var NIP05_REGEX = regexp.MustCompile(`^(?:([\w.+-]+)@)?([\w_-]+(\.[\w_-]+)+)$`)

type WellKnownResponse struct {
	Names  map[string]string   `json:"names"`
	Relays map[string][]string `json:"relays,omitempty"`
	NIP46  map[string][]string `json:"nip46,omitempty"`
}

func IsValidIdentifier(input string) bool {
	return NIP05_REGEX.MatchString(input)
}

func ParseIdentifier(fullname string) (name string, domain string, err error) {
	res := NIP05_REGEX.FindStringSubmatch(fullname)
	if len(res) == 0 {
		return "", "", fmt.Errorf("invalid identifier")
	}
	if res[1] == "" {
		res[1] = "_"
	}
	return res[1], res[2], nil
}

func QueryIdentifier(ctx context.Context, fullname string) (*nostr.ProfilePointer, error) {
	result, name, err := Fetch(ctx, fullname)
	if err != nil {
		return nil, err
	}

	pubkey, ok := result.Names[name]
	if !ok {
		return nil, fmt.Errorf("no entry for name '%s'", name)
	}

	if !nostr.IsValidPublicKey(pubkey) {
		return nil, fmt.Errorf("got an invalid public key '%s'", pubkey)
	}

	relays, _ := result.Relays[pubkey]
	return &nostr.ProfilePointer{
		PublicKey: pubkey,
		Relays:    relays,
	}, nil
}

var httpClient = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func Fetch(ctx context.Context, fullname string) (resp WellKnownResponse, name string, err error) {
	name, domain, err := ParseIdentifier(fullname)
	if err != nil {
		return resp, name, fmt.Errorf("failed to parse '%s': %w", fullname, err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("https://%s/.well-known/nostr.json?name=%s", domain, name), nil)
	if err != nil {
		return resp, name, fmt.Errorf("failed to create a request: %w", err)
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return resp, name, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	var result WellKnownResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return resp, name, fmt.Errorf("failed to decode json response: %w", err)
	}

	return result, name, nil
}

func NormalizeIdentifier(fullname string) string {
	if strings.HasPrefix(fullname, "_@") {
		return fullname[2:]
	}

	return fullname
}
