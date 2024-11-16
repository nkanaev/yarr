// Copyright (c) 2024 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
//
// Automatic selection originally written by Dave Collins July 2024.

//go:build !amd64 || purego

package compress

// Blocks performs BLAKE-224 and BLAKE-256 block compression using pure Go.  It
// will compress as many full blocks as are available in the provided message.
//
// The parameters are:
//
//	state: the block compression state with the chain value and salt
//	msg: the padded message to compress (must be at least 64 bytes)
//	counter: the total number of message bits hashed so far
//
// It will panic if the provided message block does not have at least 64 bytes.
//
// The chain value in the passed state is updated in place.
func Blocks(state *State, msg []byte, counter uint64) {
	blocksGeneric(state, msg, counter)
}
