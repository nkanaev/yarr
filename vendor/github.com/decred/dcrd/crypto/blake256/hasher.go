// Copyright (c) 2024 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
//
// Main Go code originally written and optimized by Dave Collins May 2020.
// Additional cleanup and comments added July 2024.

// Package blake256 implements BLAKE-256 and BLAKE-224 with SSE2, SSE4.1, and
// AVX acceleration and zero allocations.
package blake256

import (
	"encoding/binary"
	"fmt"

	"github.com/decred/dcrd/crypto/blake256/internal/compress"
)

const (
	// BlockSize is the block size of the hash algorithm in bytes.
	BlockSize = 64

	// Size is the size of a BLAKE-256 hash in bytes.
	Size = 32

	// Size224 is the size of a BLAKE-224 hash in bytes.
	Size224 = 28

	// SavedStateSize is the number of bytes of a serialized intermediate state.
	SavedStateSize = 128
)

// pad provides an efficient means to pad a message.
var pad = [64]byte{0x80}

// hasher implements a zero-allocation rolling BLAKE checksum.  It can safely be
// copied at any point to save its internal state for use in additional
// processing later, without having to write the previously written data again.
//
// It contains the common logic between BLAKE-224 and BLAKE-256.
type hasher struct {
	state compress.State  // the current chain value and salt
	count uint64          // running total of message bits hashed
	buf   [BlockSize]byte // partial block data buffer
	nbuf  uint32          // number of bytes written to data buffer
}

// makeHasher returns an instance of a rolling hasher initialized with the
// provided chain value.
func makeHasher(cv [8]uint32) hasher {
	return hasher{state: compress.State{CV: cv}}
}

// reset resets the state of the rolling hash.
func (h *hasher) reset(iv [8]uint32) {
	h.state.CV = iv
	h.count = 0
	h.nbuf = 0
}

// initializeSalt initialize the hasher state with the provided salt.  Note that
// this must only be done when first creating the hasher state for correct
// results.
//
// It will panic if the provided salt is not 16 bytes.
func (h *hasher) initializeSalt(salt []byte) {
	if len(salt) != 16 {
		panic("salt length must be 16 bytes")
	}
	h.state.S[0] = binary.BigEndian.Uint32(salt)
	h.state.S[1] = binary.BigEndian.Uint32(salt[4:])
	h.state.S[2] = binary.BigEndian.Uint32(salt[8:])
	h.state.S[3] = binary.BigEndian.Uint32(salt[12:])
}

// write adds the given bytes to the rolling hash.
//
// NOTE: This method only returns an error in order to satisfy the [io.Writer]
// and [hash.Hash] interfaces.  However, it will never error, meaning the error
// will always be nil, so it is safe to ignore.
func (h *hasher) write(b []byte) (int, error) {
	// All bytes will be written.
	totalWritten := len(b)

	// When a partial block exists and adding the new data would meet or exceed
	// the size of a block, fill up the partial block and compress it.
	if h.nbuf > 0 && h.nbuf+uint32(len(b)) >= BlockSize {
		written := uint32(copy(h.buf[h.nbuf:], b))
		h.count += BlockSize << 3
		compress.Blocks(&h.state, h.buf[:], h.count)
		b = b[written:]
		h.nbuf = 0
	}

	// The previous section ensures there is no partial block data remaining.
	//
	// Use that fact to compress full blocks directly when the remaining number
	// of bytes to write will completely fill one or more additional blocks.
	//
	// It is perhaps also worth noting that this approach is used over having a
	// compression function that only accepts a single block because it provides
	// a rather significant speed advantage on inputs that are larger than the
	// size of a couple of blocks while only having a negligible impact on small
	// inputs.
	if len(b) >= BlockSize {
		h.count += BlockSize << 3
		compress.Blocks(&h.state, b, h.count)

		// Update the count of message bits hashed and slice of remaining
		// unwritten bytes to account for the total number of blocks compressed.
		bytesHashed := uint64(len(b) &^ (BlockSize - 1))
		h.count += (bytesHashed - BlockSize) << 3
		b = b[bytesHashed:]
	}

	// Write any remaining bytes to the next partial block.  Note the number of
	// remaining bytes is guaranteed to be less than the size of a full block
	// due to the previous sections.
	if len(b) > 0 {
		h.nbuf += uint32(copy(h.buf[h.nbuf:], b))
	}

	return totalWritten, nil
}

// writeByte adds the given byte to the rolling hash.
func (h *hasher) writeByte(b byte) {
	var buf [1]byte
	buf[0] = b
	h.write(buf[:])
}

// writeString adds the given string to the rolling hash.
func (h *hasher) writeString(s string) {
	h.write([]byte(s))
}

// writeUint16LE encodes the given unsigned 16-bit integer as a 2-byte
// little-endian byte sequence and adds it to the rolling hash.
func (h *hasher) writeUint16LE(v uint16) {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], v)
	h.write(buf[:])
}

// writeUint16BE encodes the given unsigned 16-bit integer as a 2-byte
// big-endian byte sequence and adds it to the rolling hash.
func (h *hasher) writeUint16BE(v uint16) {
	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], v)
	h.write(buf[:])
}

// writeUint32LE encodes the given unsigned 32-bit integer as a 4-byte
// little-endian byte sequence and adds it to the rolling hash.
func (h *hasher) writeUint32LE(v uint32) {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], v)
	h.write(buf[:])
}

// writeUint32BE encodes the given unsigned 32-bit integer as a 4-byte
// big-endian byte sequence and adds it to the rolling hash.
func (h *hasher) writeUint32BE(v uint32) {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], v)
	h.write(buf[:])
}

// writeUint64LE encodes the given unsigned 64-bit integer as an 8-byte
// little-endian byte sequence and adds it to the rolling hash.
func (h *hasher) writeUint64LE(v uint64) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], v)
	h.write(buf[:])
}

// writeUint64BE encodes the given unsigned 64-bit integer as an 8-byte
// big-endian byte sequence and adds it to the rolling hash.
func (h *hasher) writeUint64BE(v uint64) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], v)
	h.write(buf[:])
}

// finalize finalizes of the rolling hash by writing any remaining partial block
// data and appending the necessary padding.
//
// The hasher may no longer be used after invoking this method.  Callers always
// run finalize on a copy of the hasher so the original hasher state is not
// modified.
//
// The length preamble bit MUST be 0 (for BLAKE-224) or 1 (for BLAKE-256).
func (h *hasher) finalize(lenPreambleBit uint8) {
	// Hashing a message consists of padding the message to a multiple of the
	// block size and processing it block per block by the compression function.
	//
	// Padding the message consists of first extending the message so that its
	// bit length is congruent to 447 modulo 512 by appending a 1 bit followed
	// by enough 0s to reach the required congruence.  Then a length preamble
	// bit is added (1 for BLAKE-256, 0 for BLAKE-224) followed by the length
	// of original message encoded as a 64-bit unsigned big-endian integer.
	// This ensures the message length is a multiple of the block size since
	// 447+1+64 = 512.
	//
	// Note that a special case occurs when the final block contains no original
	// message bit.  In that case, the message bit counter provided to the
	// compression function is set to zero for that final block.  This
	// guarantees unique blocks.
	//
	// This implementation performs iterated hashing by compressing full blocks
	// as data is written and storing the resulting chain value, total number of
	// message bits compressed, and any remaining partial block data in the
	// state.
	//
	// Thus, finalization consists of writing any remaining partial block data
	// that hasn't already been compressed and padding the message out per the
	// above.
	//
	// Since this implementation only allows writing full 8-bit bytes at a time,
	// the following is optimized to only consider message bit lengths that are
	// multiples of 8.  Concretely, note that floor(447/8) = 55.  Therefore, as
	// long as the remaining partial block data is <= 55, only one compression
	// is needed.  Otherwise a second compression is needed.
	msgBitLen := h.count + uint64(h.nbuf)<<3
	switch {
	// Exactly one padding byte is needed.
	case h.nbuf == 55:
		h.buf[55] = 0x80 | lenPreambleBit
		binary.BigEndian.PutUint64(h.buf[56:], msgBitLen)
		compress.Blocks(&h.state, h.buf[:], msgBitLen)
		return

	// Appending the padding to the remaining partial block data will fit
	// without needing another block.
	case h.nbuf < 55:
		copy(h.buf[h.nbuf:55], pad[:])
		h.buf[55] = lenPreambleBit
		binary.BigEndian.PutUint64(h.buf[56:], msgBitLen)

		// Per the specification, the counter is set to zero for the final
		// compression when the final block contains no bits from the original
		// message.
		if h.nbuf == 0 {
			msgBitLen = 0
		}
		compress.Blocks(&h.state, h.buf[:], msgBitLen)
		return
	}

	// The partial block data plus the padding and message bit length exceed the
	// size of a block, so two compressions are needed where the second one is
	// a padding block (all zeros except for the final 8 bytes which house the
	// original message length encoded as a 64-bit unsigned big-endian integer).

	// Pad the remaining partial block data and compress it.
	copy(h.buf[h.nbuf:], pad[:])
	compress.Blocks(&h.state, h.buf[:], msgBitLen)

	// Create the final padding block and compress it.
	//
	// Note that since the padding block does not contain any bits from the
	// original message, the counter is set to zero when performing compression
	// per the specification.
	copy(h.buf[:], pad[1:56])
	h.buf[55] = lenPreambleBit
	binary.BigEndian.PutUint64(h.buf[56:], msgBitLen)
	compress.Blocks(&h.state, h.buf[:], 0)
}

// wordsToBytes224 converts an array of 8 32-bit unsigned big-endian words to an
// array of 28 bytes.  The final word is truncated.
func wordsToBytes224(cv [8]uint32) (out [28]byte) {
	binary.BigEndian.PutUint32(out[24:], cv[6])
	binary.BigEndian.PutUint32(out[20:], cv[5])
	binary.BigEndian.PutUint32(out[16:], cv[4])
	binary.BigEndian.PutUint32(out[12:], cv[3])
	binary.BigEndian.PutUint32(out[8:], cv[2])
	binary.BigEndian.PutUint32(out[4:], cv[1])
	binary.BigEndian.PutUint32(out[0:], cv[0])
	return out
}

// finalize224 finalizes of the rolling hash by writing any remaining partial
// block data and appending the necessary padding for BLAKE-224.
//
// The hasher may no longer be used after invoking this method.  Callers always
// run finalize on a copy of the hasher so the original hasher state is not
// modified.
func (h *hasher) finalize224() [Size224]byte {
	const lenPreambleBit = 0x00
	h.finalize(lenPreambleBit)
	return wordsToBytes224(h.state.CV)
}

// wordsToBytes256 converts an array of 8 32-bit unsigned big-endian words to an
// array of 32 bytes.
func wordsToBytes256(cv [8]uint32) (out [32]byte) {
	binary.BigEndian.PutUint32(out[28:], cv[7])
	binary.BigEndian.PutUint32(out[24:], cv[6])
	binary.BigEndian.PutUint32(out[20:], cv[5])
	binary.BigEndian.PutUint32(out[16:], cv[4])
	binary.BigEndian.PutUint32(out[12:], cv[3])
	binary.BigEndian.PutUint32(out[8:], cv[2])
	binary.BigEndian.PutUint32(out[4:], cv[1])
	binary.BigEndian.PutUint32(out[0:], cv[0])
	return out
}

// finalize256 finalizes of the rolling hash by writing any remaining partial
// block data and appending the necessary padding for BLAKE-256.
//
// The hasher may no longer be used after invoking this method.  Callers always
// run finalize on a copy of the hasher so the original hasher state is not
// modified.
func (h *hasher) finalize256() [Size]byte {
	const lenPreambleBit = 0x01
	h.finalize(lenPreambleBit)
	return wordsToBytes256(h.state.CV)
}

// putSavedState serializes the intermediate state directly into the passed byte
// slice.  The target slice MUST have at least [SavedStateSize] bytes available
// or it will panic.
func (h *hasher) putSavedState(target []byte, prefix uint32) {
	var offset uint32
	binary.BigEndian.PutUint32(target[offset:], prefix)
	offset += 4
	for _, cv := range h.state.CV {
		binary.BigEndian.PutUint32(target[offset:], cv)
		offset += 4
	}
	for _, s := range h.state.S {
		binary.BigEndian.PutUint32(target[offset:], s)
		offset += 4
	}
	binary.BigEndian.PutUint64(target[offset:], h.count)
	offset += 8
	offset += uint32(copy(target[offset:], h.buf[:]))
	binary.BigEndian.PutUint32(target[offset:], h.nbuf)
}

// saveState appends the current intermediate state of the rolling hash prefixed
// by the passed value to the provided slice and returns the resulting slice.
// It does not change the underlying hash state.
//
// The provided prefix is expected to either be [statePrefix224] or
// [statePrefix256] depending on which hash variant is being saved.
//
// As described by the [hasher] documentation, the hasher instance can simply be
// copied to achieve the same result much more efficiently when the caller is
// able to keep a copy.  Therefore, that approach should be preferred when
// possible.
//
// However, the ability to serialize the state is also provided to enable
// sharing it across process boundaries.
func (h *hasher) saveState(target []byte, prefix uint32) []byte {
	// Create a new array and append it to the target when there is not enough
	// space remaining in the slice.  Otherwise, write directly into it.
	//
	// Note that this could alternatively just grow the slice if needed and then
	// write directly into it unconditionally, but this approach is faster for
	// the two much more common cases of the caller providing a slice that is
	// already big enough or a nil slice.
	if needed := SavedStateSize - (cap(target) - len(target)); needed > 0 {
		var state [SavedStateSize]byte
		h.putSavedState(state[:], prefix)
		return append(target, state[:]...)
	}
	h.putSavedState(target[len(target):len(target)+SavedStateSize], prefix)
	return target[:len(target)+SavedStateSize]
}

// loadState restores the rolling hash to the provided serialized intermediate
// state.  See [hasher.saveState] for more details.
//
// The provided prefix is expected to either be [statePrefix224] or
// [statePrefix256] depending on which hash variant is being loaded.
//
// [ErrMalformedState] will be returned when the provided serialized state is
// not at least the required [SavedStateSize] number of bytes.
//
// [ErrMismatchedState] will be returned if the prefix in the serialized state
// does not match the given required prefix.
func (h *hasher) loadState(state []byte, requiredPrefix uint32) error {
	if len(state) < SavedStateSize {
		str := fmt.Sprintf("malformed intermediate state - must be at least "+
			"%d bytes", SavedStateSize)
		return makeError(ErrMalformedState, str)
	}
	var offset uint32
	if pre := binary.BigEndian.Uint32(state[offset:]); pre != requiredPrefix {
		hashType := "BLAKE-256"
		if requiredPrefix != statePrefix256 {
			hashType = "BLAKE-224"
		}
		str := fmt.Sprintf("the provided intermediate state is not for %s",
			hashType)
		return makeError(ErrMismatchedState, str)
	}
	offset += 4
	for i := range h.state.CV {
		h.state.CV[i] = binary.BigEndian.Uint32(state[offset:])
		offset += 4
	}
	for i := range h.state.S {
		h.state.S[i] = binary.BigEndian.Uint32(state[offset:])
		offset += 4
	}
	h.count = binary.BigEndian.Uint64(state[offset:])
	offset += 8
	offset += uint32(copy(h.buf[:], state[offset:]))
	h.nbuf = binary.BigEndian.Uint32(state[offset:])
	return nil
}
