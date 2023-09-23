package main

import (
	"strings"
	"testing"
)

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
