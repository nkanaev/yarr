package main

import (
	"strings"
	"testing"
	"time"
)

func TestResolveDuratoin(t *testing.T) {
	for _, tc := range [...]struct {
		raw      string
		expected time.Duration
		error    bool
	}{
		{raw: "", expected: 30 * time.Second},
		{raw: "120s", expected: 120 * time.Second},
		{raw: "2m", expected: 2 * time.Minute},
		{raw: "abc", error: true},
		{raw: "0", error: true},
		{raw: "-5s", error: true},
	} {
		t.Run(tc.raw, func(t *testing.T) {
			got, err := resolveDuration(tc.raw)
			if tc.error {
				if err == nil {
					t.Errorf("expected error for %q, got nil", tc.raw)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.expected {
				t.Errorf("got %v; want %v", got, tc.expected)
			}
		})
	}
}

func TestPasswordFromAuthfile(t *testing.T) {
	for _, tc := range [...]struct {
		authfile         string
		expectedUsername string
		expectedPassword string
		expectedError    bool
	}{
		{
			authfile:         "username:password",
			expectedUsername: "username",
			expectedPassword: "password",
			expectedError:    false,
		},
		{
			authfile:      "username-and-no-password",
			expectedError: true,
		},
		{
			authfile:         "username:password:with:columns",
			expectedUsername: "username",
			expectedPassword: "password:with:columns",
			expectedError:    false,
		},
		{
			authfile:      "",
			expectedError: true,
		},
		{
			authfile:      "   :   ",
			expectedError: true,
		},
	} {
		t.Run(tc.authfile, func(t *testing.T) {
			username, password, err := parseAuthfile(strings.NewReader(tc.authfile))
			if tc.expectedUsername != username {
				t.Errorf("expected username %q, got %q", tc.expectedUsername, username)
			}
			if tc.expectedPassword != password {
				t.Errorf("expected password %q, got %q", tc.expectedPassword, password)
			}
			if tc.expectedError && err == nil {
				t.Errorf("expected error, got nil")
			} else if !tc.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
