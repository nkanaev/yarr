package server

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var instapaperAddURL = "https://www.instapaper.com/api/add"

// InstapaperAdd saves a URL to Instapaper using the Simple API with HTTP Basic auth.
func InstapaperAdd(username, password, articleURL, title string) error {
	form := url.Values{}
	form.Set("url", articleURL)
	if title != "" {
		form.Set("title", title)
	}

	req, err := http.NewRequest("POST", instapaperAddURL, nil)
	if err != nil {
		return fmt.Errorf("instapaper: %w", err)
	}
	req.URL.RawQuery = form.Encode()
	req.SetBasicAuth(username, password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("instapaper: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("instapaper: unexpected status %d: %s", resp.StatusCode, string(body))
}
