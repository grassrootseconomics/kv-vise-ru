package data

import (
	"encoding/binary"
)

// EncodeKey returns the byte representation of a urdt/ussd user data key.
// This key can be used to lookup the value on any storage backend implementation.
// Warning: This is a shortcut specifically for user data, it is not expected to work for all go-vise keys.
// TODO: Replace with imported data types from the common package once lib-gdbm dependency is removed.
func EncodeKey(sessionID string, dataType uint16) []byte {
	keyBytes := []byte(sessionID)
	keyBytes = append(keyBytes, dot)
	keyBytes = append(keyBytes, uint16ToBytes(dataType)...)
	keyBytes = append([]uint8{keyPrefix}, keyBytes...)

	return keyBytes
}

func uint16ToBytes(v uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, v)
	return bytes
}

func EncodeSessionID(sessionID string) []byte {
	sessionIDBytes := []byte(sessionID)
	sessionIDBytes = append(sessionIDBytes, dot)
	sessionIDBytes = append([]uint8{keyPrefix}, sessionIDBytes...)
	return sessionIDBytes
}
