package crypto

import (
	"testing"
)

func TestRSAGenerator_Generate(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
	}{
		{
			name:      "success - generate RSA key pair",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := &RSAGenerator{}
			keyPair, err := generator.Generate()

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if keyPair == nil {
					t.Error("expected key pair, got nil")
				}
				if keyPair.Public == nil {
					t.Error("expected public key, got nil")
				}
				if keyPair.Private == nil {
					t.Error("expected private key, got nil")
				}
			}
		})
	}
}

func TestECCGenerator_Generate(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
	}{
		{
			name:      "success - generate ECC key pair",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := &ECCGenerator{}
			keyPair, err := generator.Generate()

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if keyPair == nil {
					t.Error("expected key pair, got nil")
				}
				if keyPair.Public == nil {
					t.Error("expected public key, got nil")
				}
				if keyPair.Private == nil {
					t.Error("expected private key, got nil")
				}
			}
		})
	}
}
