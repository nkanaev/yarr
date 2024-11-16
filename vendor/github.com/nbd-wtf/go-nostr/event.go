package nostr

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/mailru/easyjson"
)

type Event struct {
	ID        string
	PubKey    string
	CreatedAt Timestamp
	Kind      int
	Tags      Tags
	Content   string
	Sig       string

	// anything here will be mashed together with the main event object when serializing
	extra map[string]any
}

// Event Stringer interface, just returns the raw JSON as a string.
func (evt Event) String() string {
	j, _ := easyjson.Marshal(evt)
	return string(j)
}

// GetID serializes and returns the event ID as a string.
func (evt *Event) GetID() string {
	h := sha256.Sum256(evt.Serialize())
	return hex.EncodeToString(h[:])
}

// CheckID checks if the implied ID matches the given ID
func (evt *Event) CheckID() bool {
	ser := evt.Serialize()
	h := sha256.Sum256(ser)

	const hextable = "0123456789abcdef"

	for i := 0; i < 32; i++ {
		b := hextable[h[i]>>4]
		if b != evt.ID[i*2] {
			return false
		}

		b = hextable[h[i]&0x0f]
		if b != evt.ID[i*2+1] {
			return false
		}
	}

	return true
}

// Serialize outputs a byte array that can be hashed/signed to identify/authenticate.
// JSON encoding as defined in RFC4627.
func (evt *Event) Serialize() []byte {
	// the serialization process is just putting everything into a JSON array
	// so the order is kept. See NIP-01
	dst := make([]byte, 0)

	// the header portion is easy to serialize
	// [0,"pubkey",created_at,kind,[
	dst = append(dst, []byte(
		fmt.Sprintf(
			"[0,\"%s\",%d,%d,",
			evt.PubKey,
			evt.CreatedAt,
			evt.Kind,
		))...)

	// tags
	dst = evt.Tags.marshalTo(dst)
	dst = append(dst, ',')

	// content needs to be escaped in general as it is user generated.
	dst = escapeString(dst, evt.Content)
	dst = append(dst, ']')

	return dst
}
