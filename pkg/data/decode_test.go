package data

import (
	"encoding/hex"
	"testing"
)

func TestDecodeKey(t *testing.T) {
	type want struct {
		sessionID string
		dataType  uint16
	}
	type args struct {
		keyBytesHex string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"blockchain_address",
			args{
				keyBytesHex: "",
			},
			want{
				sessionID: "+2547789654",
				dataType:  DATA_PUBLIC_KEY,
			},
		},
	}
	for _, tt := range tests {
		t.Skip()

		t.Run(tt.name, func(t *testing.T) {
			keyBytes, err := hex.DecodeString(tt.args.keyBytesHex)
			if err != nil {
				t.Fatalf("failed to decode hex string: %v", err)
			}

			dataType, sessionID := DecodeKey(keyBytes)
			t.Logf("%s data_type: %d, session_id: %s", tt.name, dataType, sessionID)
		})
	}
}
