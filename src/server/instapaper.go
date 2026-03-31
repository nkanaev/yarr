package server

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	instapaperAccessTokenURL = "https://www.instapaper.com/api/1/oauth/access_token"
	instapaperAddBookmarkURL = "https://www.instapaper.com/api/1/bookmarks/add"
)

type InstapaperClient struct {
	ConsumerKey    string
	ConsumerSecret string
}

func percentEncode(s string) string {
	return url.QueryEscape(s)
}

func signatureBaseString(method, baseURL string, params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, percentEncode(k)+"="+percentEncode(params[k]))
	}
	paramStr := strings.Join(pairs, "&")

	return method + "&" + percentEncode(baseURL) + "&" + percentEncode(paramStr)
}

func hmacSHA1Sign(key, data string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func generateNonce() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func (c *InstapaperClient) signedRequest(method, endpoint string, extraParams map[string]string, oauthToken, oauthTokenSecret string) (*http.Response, error) {
	params := map[string]string{
		"oauth_consumer_key":     c.ConsumerKey,
		"oauth_nonce":            generateNonce(),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        fmt.Sprintf("%d", time.Now().Unix()),
		"oauth_version":          "1.0",
	}
	if oauthToken != "" {
		params["oauth_token"] = oauthToken
	}
	for k, v := range extraParams {
		params[k] = v
	}

	baseStr := signatureBaseString(method, endpoint, params)
	signingKey := percentEncode(c.ConsumerSecret) + "&" + percentEncode(oauthTokenSecret)
	params["oauth_signature"] = hmacSHA1Sign(signingKey, baseStr)

	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	req, err := http.NewRequest(method, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return http.DefaultClient.Do(req)
}

func (c *InstapaperClient) GetAccessToken(username, password string) (token, secret string, err error) {
	extra := map[string]string{
		"x_auth_mode":     "client_auth",
		"x_auth_username": username,
		"x_auth_password": password,
	}

	resp, err := c.signedRequest("POST", instapaperAccessTokenURL, extra, "", "")
	if err != nil {
		return "", "", fmt.Errorf("instapaper auth request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("instapaper auth read failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("instapaper auth rejected (status %d): %s", resp.StatusCode, string(body))
	}

	vals, err := url.ParseQuery(string(body))
	if err != nil {
		return "", "", fmt.Errorf("instapaper auth parse failed: %w", err)
	}

	return vals.Get("oauth_token"), vals.Get("oauth_token_secret"), nil
}

func (c *InstapaperClient) AddBookmark(oauthToken, oauthTokenSecret, articleURL, title string) error {
	extra := map[string]string{
		"url": articleURL,
	}
	if title != "" {
		extra["title"] = title
	}

	resp, err := c.signedRequest("POST", instapaperAddBookmarkURL, extra, oauthToken, oauthTokenSecret)
	if err != nil {
		return fmt.Errorf("instapaper add bookmark failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("instapaper add bookmark rejected (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
