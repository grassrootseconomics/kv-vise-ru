package data

import (
	"encoding/binary"
)

const ()

// DecodeKey specifically only decodes user data keys stored as bytes into its respective session ID and data type
// TODO: Replace return data type with imported data types from the common package once lib-gdbm dependency is removed.
// Note: 0x2e was added herehttps://holbrook.no/src/go-vise/file/db/db.go.html#l147, so we discard the last 3 bytes
func DecodeKey(key []byte) (uint16, string) {
	if key[0] != keyPrefix {
		return 0, ""
	}

	return binary.BigEndian.Uint16(key[len(key)-2:]), string(key[1 : len(key)-3])
}

// DecodeValue returns the utf-8 string representation of the value stored in the storage backend
func DecodeValue(v []byte) string {
	return string(v)
}
