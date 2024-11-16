// Copyright (c) 2024 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blake256

// ErrorKind identifies a kind of error.
type ErrorKind string

// These constants are used to identify a specific ErrorKind.
const (
	// ErrMalformedState indicates a serialized intermediate state is malformed
	// in some way such as not having at least the expected number of bytes.
	ErrMalformedState = ErrorKind("ErrMalformedState")

	// ErrMismatchedState indicates a serialized intermediate state is not for
	// the hash type that is attempting to restore it.  For example, it will be
	// returned when attempting to restore a BLAKE-256 intermediate state with
	// a BLAKE-224 hasher.
	ErrMismatchedState = ErrorKind("ErrMismatchedState")
)

// Error satisfies the error interface and prints human-readable errors.
func (e ErrorKind) Error() string {
	return string(e)
}

// Error identifies an error related to restoring an intermediate hashing state.
//
// It has full support for [errors.Is] and [errors.As], so the caller can
// ascertain the specific reason for the error by checking the underlying error.
type Error struct {
	Err         error
	Description string
}

// Error satisfies the error interface and prints human-readable errors.
func (e Error) Error() string {
	return e.Description
}

// Unwrap returns the underlying wrapped error.
func (e Error) Unwrap() error {
	return e.Err
}

// makeError creates an [Error] given a set of arguments.
func makeError(kind ErrorKind, desc string) Error {
	return Error{Err: kind, Description: desc}
}
