package server

import (
	"testing"
)

func TestOAuthSignatureBaseString(t *testing.T) {
	params := map[string]string{
		"oauth_consumer_key":     "testkey",
		"oauth_nonce":            "testnonce",
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        "1234567890",
		"oauth_version":          "1.0",
		"x_auth_mode":            "client_auth",
		"x_auth_username":        "user@example.com",
		"x_auth_password":        "password123",
	}

	base := signatureBaseString("POST", "https://www.instapaper.com/api/1/oauth/access_token", params)

	if base == "" {
		t.Fatal("signature base string should not be empty")
	}
	if base[:4] != "POST" {
		t.Errorf("base string should start with POST, got %s", base[:4])
	}
}

func TestOAuthHMACSHA1(t *testing.T) {
	sig := hmacSHA1Sign("consumerSecret&", "base")
	if sig == "" {
		t.Fatal("signature should not be empty")
	}
}
