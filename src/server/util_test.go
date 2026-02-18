package server

import "testing"

func TestIsInternalFromURL(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"http://192.168.1.1:8080", true},
		{"http://10.0.0.5", true},
		{"http://172.16.0.1", true},
		{"http://172.31.255.255", true},
		{"http://172.32.0.1", false}, // outside private range
		{"http://127.0.0.1", true},
		{"http://127.0.0.1:7000", true},
		{"http://127.0.0.1:7000/secret", true},
		{"http://169.254.0.5", true},
		{"http://localhost", true}, // resolves to 127.0.0.1
		{"http://8.8.8.8", false},
		{"http://google.com", false}, // resolves to public IPs
		{"invalid-url", false},       // invalid format
		{"", false},                  // empty string
	}
	for _, test := range tests {
		result := isInternalFromURL(test.url)
		if result != test.expected {
			t.Errorf("isInternalFromURL(%q) = %v; want %v", test.url, result, test.expected)
		}
	}
}