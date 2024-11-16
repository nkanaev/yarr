// Copyright (c) 2024 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
//
// Main Go code originally written and optimized by Dave Collins May 2020.
// Additional cleanup and comments added July 2024.

package compress

import (
	"math/bits"
)

const (
	// blockSize is the block size of the hash algorithm in bytes.
	blockSize = 64

	// blockSizeLog2 is the base-2 log of the block size.  It is used to
	// efficiently perform integer division by the block size.
	blockSizeLog2 = 6
)

// State houses the current chain value and salt used during block compression.
// It is housed in a separate struct in order to reduce the number of parameters
// to help prevent register spillage on platforms with a limited number of
// registers for better performance.
type State struct {
	CV [8]uint32 // the current chain value
	S  [4]uint32 // salt (zero by default)
}

// g is the quarter round function that each round applies to the 4x4 internal
// state in the compression function.
func g(a, b, c, d, mx, my, cx, cy uint32) (uint32, uint32, uint32, uint32) {
	a += b + (mx ^ cx)
	d = bits.RotateLeft32(d^a, -16)
	c += d
	b = bits.RotateLeft32(b^c, -12)
	a += b + (my ^ cy)
	d = bits.RotateLeft32(d^a, -8)
	c += d
	b = bits.RotateLeft32(b^c, -7)
	return a, b, c, d
}

// blocksGeneric performs BLAKE-224 and BLAKE-256 block compression using pure
// Go.  It will compress as many full blocks as are available in the provided
// message.
//
// The parameters are:
//
//	state: the block compression state with the chain value and salt
//	msg: the padded message to compress (must be at least 64 bytes)
//	counter: the total number of message bits hashed so far
//
// It will panic if the provided message block does not have at least 64 bytes.
//
// The chain value in the provided state is updated in place.
func blocksGeneric(state *State, msg []byte, counter uint64) {
	// The compression func initializes the 16-word state matrix as follows:
	//
	// h0..h7 is the input chaining value.
	//
	// s0..s3 is the salt.
	//
	// c0..c7 are the first 8 words of the BLAKE constants defined by the
	// specification (the first digits of π).
	//
	// t0 and t1 are the lower and higher order words of the 64-bit counter.
	//
	// |v0  v1  v2  v3|   |  h0     h1     h2     h3 |
	// |v4  v5  v6  v7|   |  h4     h5     h6     h7 |
	// |v8  v9  va  vb| = |s0^c0  s1^c1  s2^c2  s3^c3|
	// |vc  vd  ve  vf|   |t0^c4  t0^c5  t1^c6  t1^c7|
	//
	// Each round consists of 8 applications of the G function as follows:
	// G0(v0,v4,v8,vc)  G1(v1,v5,v9,vd)  G2(v2,v6,va,ve)  G3(v3,v7,vb,vf)
	// G4(v0,v5,va,vf)  G5(v1,v6,vb,vc)  G6(v2,v7,v8,vd)  G7(v3,v4,v9,ve)
	//
	// In other words, the G function is applied to each column of the 4x4 state
	// and then to each of the diagonals.
	//
	// In addition, at each round, the message words are permuted according to
	// the following schedule where each round is mod 10.  In other words round
	// 11 uses the same permutation as round 1, round 12 the same as round 2,
	// etc:
	//
	// round 1:  0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15
	// round 2:  14 10 4  8  9  15 13 6  1  12 0  2  11 7  5  3
	// round 3:  11 8  12 0  5  2  15 13 10 14 3  6  7  1  9  4
	// round 4:  7  9  3  1  13 12 11 14 2  6  5  10 4  0  15 8
	// round 5:  9  0  5  7  2  4  10 15 14 1  11 12 6  8  3  13
	// round 6:  2  12 6  10 0  11 8  3  4  13 7  5  15 14 1  9
	// round 7:  12 5  1  15 14 13 4  10 0  7  6  3  9  2  8  11
	// round 8:  13 11 7  14 12 1  3  9  5  0  15 4  8  6  2  10
	// round 9:  6  15 14 9  11 3  0  8  12 2  13 7  1  4  10 5
	// round 10: 10 2  8  4  7  6  1  5  15 11 9  14 3  12 13 0

	const (
		c0, c1, c2, c3 = 0x243f6a88, 0x85a308d3, 0x13198a2e, 0x03707344
		c4, c5, c6, c7 = 0xa4093822, 0x299f31d0, 0x082efa98, 0xec4e6c89
		c8, c9, ca, cb = 0x452821e6, 0x38d01377, 0xbe5466cf, 0x34e90c6c
		cc, cd, ce, cf = 0xc0ac29b7, 0xc97c50dd, 0x3f84d5b5, 0xb5470917
	)

	// Ideally these hints wouldn't be necessary, but they do make quite a big
	// difference in avoiding a bunch of additional bounds checks as determined
	// by both reviewing the resulting compiled asm as well as benchmarks.
	h := &state.CV
	s := &state.S
	_ = h[7] // Bounds check hint to compiler.
	_ = s[3] // Bounds check hint to compiler.
	for numBlocks := len(msg) >> blockSizeLog2; numBlocks > 0; numBlocks-- {
		_ = msg[63] // Bounds check hint to compiler.

		// Convert the provided message of at least 64 bytes to an array of 16
		// 32-bit unsigned big-endian words.
		//
		// Note that this conversion is intentionally not in a separate function
		// because the Go compiler unfortunately sees a function that does this
		// as too complex to inline.
		//
		// It also avoids using binary.BigEndian.Uint32 for these even though it
		// is more readable for increased performance in the critical path.
		//
		// Finally, this is optimized to favor arm since there currently is only
		// a pure Go implementation for that arch while amd64 will be using the
		// vector accelerated versions in practice.  It's only slightly slower
		// on amd64 versus various rearrangements that favor it instead.
		m0 := uint32(msg[3]) | uint32(msg[2])<<8 | uint32(msg[1])<<16 | uint32(msg[0])<<24
		m1 := uint32(msg[7]) | uint32(msg[6])<<8 | uint32(msg[5])<<16 | uint32(msg[4])<<24
		m2 := uint32(msg[11]) | uint32(msg[10])<<8 | uint32(msg[9])<<16 | uint32(msg[8])<<24
		m3 := uint32(msg[15]) | uint32(msg[14])<<8 | uint32(msg[13])<<16 | uint32(msg[12])<<24
		m4 := uint32(msg[19]) | uint32(msg[18])<<8 | uint32(msg[17])<<16 | uint32(msg[16])<<24
		m5 := uint32(msg[23]) | uint32(msg[22])<<8 | uint32(msg[21])<<16 | uint32(msg[20])<<24
		m6 := uint32(msg[27]) | uint32(msg[26])<<8 | uint32(msg[25])<<16 | uint32(msg[24])<<24
		m7 := uint32(msg[31]) | uint32(msg[30])<<8 | uint32(msg[29])<<16 | uint32(msg[28])<<24
		m8 := uint32(msg[35]) | uint32(msg[34])<<8 | uint32(msg[33])<<16 | uint32(msg[32])<<24
		m9 := uint32(msg[39]) | uint32(msg[38])<<8 | uint32(msg[37])<<16 | uint32(msg[36])<<24
		m10 := uint32(msg[43]) | uint32(msg[42])<<8 | uint32(msg[41])<<16 | uint32(msg[40])<<24
		m11 := uint32(msg[47]) | uint32(msg[46])<<8 | uint32(msg[45])<<16 | uint32(msg[44])<<24
		m12 := uint32(msg[51]) | uint32(msg[50])<<8 | uint32(msg[49])<<16 | uint32(msg[48])<<24
		m13 := uint32(msg[55]) | uint32(msg[54])<<8 | uint32(msg[53])<<16 | uint32(msg[52])<<24
		m14 := uint32(msg[59]) | uint32(msg[58])<<8 | uint32(msg[57])<<16 | uint32(msg[56])<<24
		m15 := uint32(msg[63]) | uint32(msg[62])<<8 | uint32(msg[61])<<16 | uint32(msg[60])<<24

		// Round 1 plus state matrix initialization.
		t0, t1 := uint32(counter), uint32(counter>>32)
		v0, v4, v8, vc := g(h[0], h[4], s[0]^c0, t0^c4, m0, m1, c1, c0)
		v1, v5, v9, vd := g(h[1], h[5], s[1]^c1, t0^c5, m2, m3, c3, c2)
		v2, v6, va, ve := g(h[2], h[6], s[2]^c2, t1^c6, m4, m5, c5, c4)
		v3, v7, vb, vf := g(h[3], h[7], s[3]^c3, t1^c7, m6, m7, c7, c6)
		v0, v5, va, vf = g(v0, v5, va, vf, m8, m9, c9, c8)
		v1, v6, vb, vc = g(v1, v6, vb, vc, m10, m11, cb, ca)
		v2, v7, v8, vd = g(v2, v7, v8, vd, m12, m13, cd, cc)
		v3, v4, v9, ve = g(v3, v4, v9, ve, m14, m15, cf, ce)

		// Round 2 with message word and constants permutation.
		v0, v4, v8, vc = g(v0, v4, v8, vc, m14, m10, ca, ce)
		v1, v5, v9, vd = g(v1, v5, v9, vd, m4, m8, c8, c4)
		v2, v6, va, ve = g(v2, v6, va, ve, m9, m15, cf, c9)
		v3, v7, vb, vf = g(v3, v7, vb, vf, m13, m6, c6, cd)
		v0, v5, va, vf = g(v0, v5, va, vf, m1, m12, cc, c1)
		v1, v6, vb, vc = g(v1, v6, vb, vc, m0, m2, c2, c0)
		v2, v7, v8, vd = g(v2, v7, v8, vd, m11, m7, c7, cb)
		v3, v4, v9, ve = g(v3, v4, v9, ve, m5, m3, c3, c5)

		// Round 3 with message word and constants permutation.
		v0, v4, v8, vc = g(v0, v4, v8, vc, m11, m8, c8, cb)
		v1, v5, v9, vd = g(v1, v5, v9, vd, m12, m0, c0, cc)
		v2, v6, va, ve = g(v2, v6, va, ve, m5, m2, c2, c5)
		v3, v7, vb, vf = g(v3, v7, vb, vf, m15, m13, cd, cf)
		v0, v5, va, vf = g(v0, v5, va, vf, m10, m14, ce, ca)
		v1, v6, vb, vc = g(v1, v6, vb, vc, m3, m6, c6, c3)
		v2, v7, v8, vd = g(v2, v7, v8, vd, m7, m1, c1, c7)
		v3, v4, v9, ve = g(v3, v4, v9, ve, m9, m4, c4, c9)

		// Round 4 with message word and constants permutation.
		v0, v4, v8, vc = g(v0, v4, v8, vc, m7, m9, c9, c7)
		v1, v5, v9, vd = g(v1, v5, v9, vd, m3, m1, c1, c3)
		v2, v6, va, ve = g(v2, v6, va, ve, m13, m12, cc, cd)
		v3, v7, vb, vf = g(v3, v7, vb, vf, m11, m14, ce, cb)
		v0, v5, va, vf = g(v0, v5, va, vf, m2, m6, c6, c2)
		v1, v6, vb, vc = g(v1, v6, vb, vc, m5, m10, ca, c5)
		v2, v7, v8, vd = g(v2, v7, v8, vd, m4, m0, c0, c4)
		v3, v4, v9, ve = g(v3, v4, v9, ve, m15, m8, c8, cf)

		// Round 5 with message word and constants permutation.
		v0, v4, v8, vc = g(v0, v4, v8, vc, m9, m0, c0, c9)
		v1, v5, v9, vd = g(v1, v5, v9, vd, m5, m7, c7, c5)
		v2, v6, va, ve = g(v2, v6, va, ve, m2, m4, c4, c2)
		v3, v7, vb, vf = g(v3, v7, vb, vf, m10, m15, cf, ca)
		v0, v5, va, vf = g(v0, v5, va, vf, m14, m1, c1, ce)
		v1, v6, vb, vc = g(v1, v6, vb, vc, m11, m12, cc, cb)
		v2, v7, v8, vd = g(v2, v7, v8, vd, m6, m8, c8, c6)
		v3, v4, v9, ve = g(v3, v4, v9, ve, m3, m13, cd, c3)

		// Round 6 with message word and constants permutation.
		v0, v4, v8, vc = g(v0, v4, v8, vc, m2, m12, cc, c2)
		v1, v5, v9, vd = g(v1, v5, v9, vd, m6, m10, ca, c6)
		v2, v6, va, ve = g(v2, v6, va, ve, m0, m11, cb, c0)
		v3, v7, vb, vf = g(v3, v7, vb, vf, m8, m3, c3, c8)
		v0, v5, va, vf = g(v0, v5, va, vf, m4, m13, cd, c4)
		v1, v6, vb, vc = g(v1, v6, vb, vc, m7, m5, c5, c7)
		v2, v7, v8, vd = g(v2, v7, v8, vd, m15, m14, ce, cf)
		v3, v4, v9, ve = g(v3, v4, v9, ve, m1, m9, c9, c1)

		// Round 7 with message word and constants permutation.
		v0, v4, v8, vc = g(v0, v4, v8, vc, m12, m5, c5, cc)
		v1, v5, v9, vd = g(v1, v5, v9, vd, m1, m15, cf, c1)
		v2, v6, va, ve = g(v2, v6, va, ve, m14, m13, cd, ce)
		v3, v7, vb, vf = g(v3, v7, vb, vf, m4, m10, ca, c4)
		v0, v5, va, vf = g(v0, v5, va, vf, m0, m7, c7, c0)
		v1, v6, vb, vc = g(v1, v6, vb, vc, m6, m3, c3, c6)
		v2, v7, v8, vd = g(v2, v7, v8, vd, m9, m2, c2, c9)
		v3, v4, v9, ve = g(v3, v4, v9, ve, m8, m11, cb, c8)

		// Round 8 with message word and constants permutation.
		v0, v4, v8, vc = g(v0, v4, v8, vc, m13, m11, cb, cd)
		v1, v5, v9, vd = g(v1, v5, v9, vd, m7, m14, ce, c7)
		v2, v6, va, ve = g(v2, v6, va, ve, m12, m1, c1, cc)
		v3, v7, vb, vf = g(v3, v7, vb, vf, m3, m9, c9, c3)
		v0, v5, va, vf = g(v0, v5, va, vf, m5, m0, c0, c5)
		v1, v6, vb, vc = g(v1, v6, vb, vc, m15, m4, c4, cf)
		v2, v7, v8, vd = g(v2, v7, v8, vd, m8, m6, c6, c8)
		v3, v4, v9, ve = g(v3, v4, v9, ve, m2, m10, ca, c2)

		// Round 9 with message word and constants permutation.
		v0, v4, v8, vc = g(v0, v4, v8, vc, m6, m15, cf, c6)
		v1, v5, v9, vd = g(v1, v5, v9, vd, m14, m9, c9, ce)
		v2, v6, va, ve = g(v2, v6, va, ve, m11, m3, c3, cb)
		v3, v7, vb, vf = g(v3, v7, vb, vf, m0, m8, c8, c0)
		v0, v5, va, vf = g(v0, v5, va, vf, m12, m2, c2, cc)
		v1, v6, vb, vc = g(v1, v6, vb, vc, m13, m7, c7, cd)
		v2, v7, v8, vd = g(v2, v7, v8, vd, m1, m4, c4, c1)
		v3, v4, v9, ve = g(v3, v4, v9, ve, m10, m5, c5, ca)

		// Round 10 with message word and constants permutation.
		v0, v4, v8, vc = g(v0, v4, v8, vc, m10, m2, c2, ca)
		v1, v5, v9, vd = g(v1, v5, v9, vd, m8, m4, c4, c8)
		v2, v6, va, ve = g(v2, v6, va, ve, m7, m6, c6, c7)
		v3, v7, vb, vf = g(v3, v7, vb, vf, m1, m5, c5, c1)
		v0, v5, va, vf = g(v0, v5, va, vf, m15, m11, cb, cf)
		v1, v6, vb, vc = g(v1, v6, vb, vc, m9, m14, ce, c9)
		v2, v7, v8, vd = g(v2, v7, v8, vd, m3, m12, cc, c3)
		v3, v4, v9, ve = g(v3, v4, v9, ve, m13, m0, c0, cd)

		// Round 11 with message word and constants permutation.
		v0, v4, v8, vc = g(v0, v4, v8, vc, m0, m1, c1, c0)
		v1, v5, v9, vd = g(v1, v5, v9, vd, m2, m3, c3, c2)
		v2, v6, va, ve = g(v2, v6, va, ve, m4, m5, c5, c4)
		v3, v7, vb, vf = g(v3, v7, vb, vf, m6, m7, c7, c6)
		v0, v5, va, vf = g(v0, v5, va, vf, m8, m9, c9, c8)
		v1, v6, vb, vc = g(v1, v6, vb, vc, m10, m11, cb, ca)
		v2, v7, v8, vd = g(v2, v7, v8, vd, m12, m13, cd, cc)
		v3, v4, v9, ve = g(v3, v4, v9, ve, m14, m15, cf, ce)

		// Round 12 with message word and constants permutation.
		v0, v4, v8, vc = g(v0, v4, v8, vc, m14, m10, ca, ce)
		v1, v5, v9, vd = g(v1, v5, v9, vd, m4, m8, c8, c4)
		v2, v6, va, ve = g(v2, v6, va, ve, m9, m15, cf, c9)
		v3, v7, vb, vf = g(v3, v7, vb, vf, m13, m6, c6, cd)
		v0, v5, va, vf = g(v0, v5, va, vf, m1, m12, cc, c1)
		v1, v6, vb, vc = g(v1, v6, vb, vc, m0, m2, c2, c0)
		v2, v7, v8, vd = g(v2, v7, v8, vd, m11, m7, c7, cb)
		v3, v4, v9, ve = g(v3, v4, v9, ve, m5, m3, c3, c5)

		// Round 13 with message word and constants permutation.
		v0, v4, v8, vc = g(v0, v4, v8, vc, m11, m8, c8, cb)
		v1, v5, v9, vd = g(v1, v5, v9, vd, m12, m0, c0, cc)
		v2, v6, va, ve = g(v2, v6, va, ve, m5, m2, c2, c5)
		v3, v7, vb, vf = g(v3, v7, vb, vf, m15, m13, cd, cf)
		v0, v5, va, vf = g(v0, v5, va, vf, m10, m14, ce, ca)
		v1, v6, vb, vc = g(v1, v6, vb, vc, m3, m6, c6, c3)
		v2, v7, v8, vd = g(v2, v7, v8, vd, m7, m1, c1, c7)
		v3, v4, v9, ve = g(v3, v4, v9, ve, m9, m4, c4, c9)

		// Round 14 with message word and constants permutation.
		v0, v4, v8, vc = g(v0, v4, v8, vc, m7, m9, c9, c7)
		v1, v5, v9, vd = g(v1, v5, v9, vd, m3, m1, c1, c3)
		v2, v6, va, ve = g(v2, v6, va, ve, m13, m12, cc, cd)
		v3, v7, vb, vf = g(v3, v7, vb, vf, m11, m14, ce, cb)
		v0, v5, va, vf = g(v0, v5, va, vf, m2, m6, c6, c2)
		v1, v6, vb, vc = g(v1, v6, vb, vc, m5, m10, ca, c5)
		v2, v7, v8, vd = g(v2, v7, v8, vd, m4, m0, c0, c4)
		v3, v4, v9, ve = g(v3, v4, v9, ve, m15, m8, c8, cf)

		// Finally the output is defined as:
		//
		// h'0 = h0^s0^v0^v8
		// h'1 = h1^s1^v1^v9
		// h'2 = h2^s2^v2^va
		// h'3 = h3^s3^v3^vb
		// h'4 = h4^s0^v4^vc
		// h'5 = h5^s1^v5^vd
		// h'6 = h6^s2^v6^ve
		// h'7 = h7^s3^v7^vf
		h[0] ^= s[0] ^ v0 ^ v8
		h[1] ^= s[1] ^ v1 ^ v9
		h[2] ^= s[2] ^ v2 ^ va
		h[3] ^= s[3] ^ v3 ^ vb
		h[4] ^= s[0] ^ v4 ^ vc
		h[5] ^= s[1] ^ v5 ^ vd
		h[6] ^= s[2] ^ v6 ^ ve
		h[7] ^= s[3] ^ v7 ^ vf

		// Move to the next message and increase the message bits counter
		// accordingly.
		msg = msg[blockSize:]
		counter += blockSize << 3
	}
}
