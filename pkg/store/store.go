package store

import (
	"context"
)

// Store defines the interface for key-value storage operations.
type Store interface {
	// GetSessionData retrieves all keys and values for a session ID.
	GetSessionData(context.Context, []byte) (map[uint16][]string, error)
	// Get address retrieves the value for a specific session ID.
	GetAddress(context.Context, string) (string, error)
}
