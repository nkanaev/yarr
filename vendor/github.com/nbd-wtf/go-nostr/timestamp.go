package nostr

import "time"

type Timestamp int64

func Now() Timestamp {
	return Timestamp(time.Now().Unix())
}

func (t Timestamp) Time() time.Time {
	return time.Unix(int64(t), 0)
}
