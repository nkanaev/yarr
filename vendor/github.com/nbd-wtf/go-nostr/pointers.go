package nostr

type ProfilePointer struct {
	PublicKey string   `json:"pubkey"`
	Relays    []string `json:"relays,omitempty"`
}

type EventPointer struct {
	ID     string   `json:"id"`
	Relays []string `json:"relays,omitempty"`
	Author string   `json:"author,omitempty"`
	Kind   int      `json:"kind,omitempty"`
}

type EntityPointer struct {
	PublicKey  string   `json:"pubkey"`
	Kind       int      `json:"kind,omitempty"`
	Identifier string   `json:"identifier,omitempty"`
	Relays     []string `json:"relays,omitempty"`
}
