package storage

import (
	"testing"
)

func TestGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      string
		options map[string]interface{}
		wantErr bool
	}{
		{
			name:    "Invalid storage ID",
			id:      "invalid",
			options: nil,
			wantErr: true,
		},
		{
			name:    "Valid FS storage",
			id:      "fs",
			options: map[string]interface{}{"path": "/tmp"},
			wantErr: false,
		},
		{
			name:    "Valid S3 storage (no available backend)",
			id:      "s3",
			options: map[string]interface{}{"bucket": "test"},
			wantErr: true,
		},
		{
			name:    "Invalid options",
			id:      "fs",
			options: map[string]interface{}{"invalid": "option"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Get(tt.id, tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("Get() returned nil storage when no error expected")
			}
		})
	}
}
