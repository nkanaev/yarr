package nostr

import "context"

type Keyer interface {
	Signer
	Cipher
}

// A Signer provides basic public key signing methods.
type Signer interface {
	GetPublicKey(context.Context) (string, error)
	SignEvent(context.Context, *Event) error
}

// A Cipher provides NIP-44 encryption and decryption methods.
type Cipher interface {
	Encrypt(ctx context.Context, plaintext string, recipientPublicKey string) (base64ciphertext string, err error)
	Decrypt(ctx context.Context, base64ciphertext string, senderPublicKey string) (plaintext string, err error)
}
