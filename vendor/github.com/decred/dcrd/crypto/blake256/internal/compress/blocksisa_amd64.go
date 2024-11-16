// Copyright (c) 2024 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
//
// Automatic selection originally written by Dave Collins July 2024.

//go:build !purego

package compress

// Blocks performs BLAKE-224 and BLAKE-256 block compression using processor
// specific vector extensions when available.  It will compress as many full
// blocks as are available in the provided message.
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
	switch {
	case hasAVX:
		blocksAVX(state, msg, counter)
	case hasSSE41:
		blocksSSE41(state, msg, counter)
	case hasSSE2:
		blocksSSE2(state, msg, counter)
	default:
		blocksGeneric(state, msg, counter)
	}
}
