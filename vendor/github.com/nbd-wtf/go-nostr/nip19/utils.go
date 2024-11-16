package nip19

import (
	"bytes"
)

const (
	TLVDefault uint8 = 0
	TLVRelay   uint8 = 1
	TLVAuthor  uint8 = 2
	TLVKind    uint8 = 3
)

func readTLVEntry(data []byte) (typ uint8, value []byte) {
	if len(data) < 2 {
		return 0, nil
	}

	typ = data[0]
	length := int(data[1])
	value = data[2 : 2+length]
	return
}

func writeTLVEntry(buf *bytes.Buffer, typ uint8, value []byte) {
	length := len(value)
	buf.WriteByte(typ)
	buf.WriteByte(uint8(length))
	buf.Write(value)
}
