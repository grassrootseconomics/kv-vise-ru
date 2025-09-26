package store

import (
	"context"
)

type (
	// Store defines the interface for key-value storage operations.
	Store interface {
		// GetSessionData retrieves all keys and values for a session ID.
		GetSessionData(context.Context, []byte) (map[uint16][]string, error)
		// Get retrieves the value for a specific key.
		GetAddress(context.Context, string) (string, error)
		// GetProfileDetails retrieves profile details for a session ID.
		GetProfileDetailsForSMS(context.Context, string) (*ProfileDetails, error)
		// ReverseAddress looks up a session ID from an EVM address.
		ReverseAddress(context.Context, string) (string, error)
	}

	// ProfileDetails contains the profile information for a session.
	ProfileDetails struct {
		PublicKey    string
		FirstName    string
		FamilyName   string
		LanguageCode string
		AccountAlias string
	}
)
