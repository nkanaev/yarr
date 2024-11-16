//go:build libsecp256k1

package nostr

/*
#cgo LDFLAGS: -lsecp256k1
#include <secp256k1.h>
#include <secp256k1_schnorrsig.h>
#include <secp256k1_extrakeys.h>
*/
import "C"

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"unsafe"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
)

// CheckSignature checks if the signature is valid for the id
// (which is a hash of the serialized event content).
// returns an error if the signature itself is invalid.
func (evt Event) CheckSignature() (bool, error) {
	var pk [32]byte
	_, err := hex.Decode(pk[:], []byte(evt.PubKey))
	if err != nil {
		return false, fmt.Errorf("event pubkey '%s' is invalid hex: %w", evt.PubKey, err)
	}

	var sig [64]byte
	_, err = hex.Decode(sig[:], []byte(evt.Sig))
	if err != nil {
		return false, fmt.Errorf("event signature '%s' is invalid hex: %w", evt.Sig, err)
	}

	msg := sha256.Sum256(evt.Serialize())

	var xonly C.secp256k1_xonly_pubkey
	if C.secp256k1_xonly_pubkey_parse(globalSecp256k1Context, &xonly, (*C.uchar)(unsafe.Pointer(&pk[0]))) != 1 {
		return false, fmt.Errorf("failed to parse xonly pubkey")
	}

	res := C.secp256k1_schnorrsig_verify(globalSecp256k1Context, (*C.uchar)(unsafe.Pointer(&sig[0])), (*C.uchar)(unsafe.Pointer(&msg[0])), 32, &xonly)
	return res == 1, nil
}

// Sign signs an event with a given privateKey.
func (evt *Event) Sign(secretKey string, signOpts ...schnorr.SignOption) error {
	sk, err := hex.DecodeString(secretKey)
	if err != nil {
		return fmt.Errorf("Sign called with invalid secret key '%s': %w", secretKey, err)
	}

	if evt.Tags == nil {
		evt.Tags = make(Tags, 0)
	}

	var keypair C.secp256k1_keypair
	if C.secp256k1_keypair_create(globalSecp256k1Context, &keypair, (*C.uchar)(unsafe.Pointer(&sk[0]))) != 1 {
		return errors.New("failed to parse private key")
	}

	var xonly C.secp256k1_xonly_pubkey
	var pk [32]byte
	C.secp256k1_keypair_xonly_pub(globalSecp256k1Context, &xonly, nil, &keypair)
	C.secp256k1_xonly_pubkey_serialize(globalSecp256k1Context, (*C.uchar)(unsafe.Pointer(&pk[0])), &xonly)
	evt.PubKey = hex.EncodeToString(pk[:])

	h := sha256.Sum256(evt.Serialize())

	var sig [64]byte
	var random [32]byte
	rand.Read(random[:])
	if C.secp256k1_schnorrsig_sign32(globalSecp256k1Context, (*C.uchar)(unsafe.Pointer(&sig[0])), (*C.uchar)(unsafe.Pointer(&h[0])), &keypair, (*C.uchar)(unsafe.Pointer(&random[0]))) != 1 {
		return errors.New("failed to sign message")
	}

	evt.ID = hex.EncodeToString(h[:])
	evt.Sig = hex.EncodeToString(sig[:])

	return nil
}

var globalSecp256k1Context *C.secp256k1_context

func init() {
	globalSecp256k1Context = C.secp256k1_context_create(C.SECP256K1_CONTEXT_SIGN | C.SECP256K1_CONTEXT_VERIFY)
	if globalSecp256k1Context == nil {
		panic("failed to create secp256k1 context")
	}
}
