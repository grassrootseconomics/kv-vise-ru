package data

import (
	"testing"
)

// TODO: Import from urdt/ussd
func TestEncodeKey(t *testing.T) {
	type args struct {
		sessionID string
		dataType  uint16
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"blockchain_address",
			args{
				sessionID: "+254711987456",
				dataType:  DATA_PUBLIC_KEY,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeKey(tt.args.sessionID, tt.args.dataType)
			t.Logf("%s key: %x", tt.name, got)
		})
	}
}

func TestEncodeSessionID(t *testing.T) {
	sessionID := "+254711987654"
	got := EncodeSessionID(sessionID)
	t.Logf("session_id: %x", got)
}
