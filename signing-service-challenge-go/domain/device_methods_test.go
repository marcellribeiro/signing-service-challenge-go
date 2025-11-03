package domain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"strings"
	"sync"
	"testing"
)

func TestGetSecuredDataToSign(t *testing.T) {
	tests := []struct {
		name             string
		device           *Device
		dataToBeSigned   string
		expectedContains []string
	}{
		{
			name: "success - first signature uses device ID",
			device: &Device{
				ID:               "test-device-id",
				SignatureCounter: 0,
				LastSignature:    "",
			},
			dataToBeSigned: "transaction data",
			expectedContains: []string{
				"0_",
				"transaction data",
				base64.StdEncoding.EncodeToString([]byte("test-device-id")),
			},
		},
		{
			name: "success - subsequent signature uses last signature",
			device: &Device{
				ID:               "test-device-id",
				SignatureCounter: 5,
				LastSignature:    "previousSignature",
			},
			dataToBeSigned: "new transaction",
			expectedContains: []string{
				"5_",
				"new transaction",
				"previousSignature",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.device.GetSecuredDataToSign(tt.dataToBeSigned)

			for _, expected := range tt.expectedContains {
				if !strings.Contains(result, expected) {
					t.Errorf("expected result to contain %q, got %q", expected, result)
				}
			}
		})
	}
}

func TestIncrementCounter(t *testing.T) {
	tests := []struct {
		name              string
		initialCounter    int
		newSignature      string
		expectedCounter   int
		expectedSignature string
	}{
		{
			name:              "success - increment from zero",
			initialCounter:    0,
			newSignature:      "signature1",
			expectedCounter:   1,
			expectedSignature: "signature1",
		},
		{
			name:              "success - increment from non-zero",
			initialCounter:    10,
			newSignature:      "signature11",
			expectedCounter:   11,
			expectedSignature: "signature11",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			device := &Device{
				SignatureCounter: tt.initialCounter,
			}

			device.IncrementCounter(tt.newSignature)

			if device.SignatureCounter != tt.expectedCounter {
				t.Errorf("expected counter %d, got %d", tt.expectedCounter, device.SignatureCounter)
			}

			if device.LastSignature != tt.expectedSignature {
				t.Errorf("expected signature %q, got %q", tt.expectedSignature, device.LastSignature)
			}
		})
	}
}

func TestIncrementCounter_Concurrency(t *testing.T) {
	device := &Device{
		SignatureCounter: 0,
	}

	var wg sync.WaitGroup
	iterations := 100

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			device.IncrementCounter("sig")
		}(i)
	}

	wg.Wait()

	if device.SignatureCounter != iterations {
		t.Errorf("expected counter %d, got %d (thread-safety issue)", iterations, device.SignatureCounter)
	}
}

func TestGetRSAPrivateKey(t *testing.T) {
	rsaKey, _ := rsa.GenerateKey(rand.Reader, 512)
	ecdsaKey, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)

	tests := []struct {
		name      string
		device    *Device
		wantError bool
	}{
		{
			name: "success - RSA device",
			device: &Device{
				Algorithm:  AlgorithmRSA,
				PrivateKey: rsaKey,
			},
			wantError: false,
		},
		{
			name: "error - ECDSA device",
			device: &Device{
				Algorithm:  AlgorithmECDSA,
				PrivateKey: ecdsaKey,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := tt.device.GetRSAPrivateKey()

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if key == nil {
					t.Error("expected key, got nil")
				}
			}
		})
	}
}

func TestGetECDSAPrivateKey(t *testing.T) {
	rsaKey, _ := rsa.GenerateKey(rand.Reader, 512)
	ecdsaKey, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)

	tests := []struct {
		name      string
		device    *Device
		wantError bool
	}{
		{
			name: "success - ECDSA device",
			device: &Device{
				Algorithm:  AlgorithmECDSA,
				PrivateKey: ecdsaKey,
			},
			wantError: false,
		},
		{
			name: "error - RSA device",
			device: &Device{
				Algorithm:  AlgorithmRSA,
				PrivateKey: rsaKey,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := tt.device.GetECDSAPrivateKey()

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if key == nil {
					t.Error("expected key, got nil")
				}
			}
		})
	}
}

func TestGetRSAPublicKey(t *testing.T) {
	rsaKey, _ := rsa.GenerateKey(rand.Reader, 512)
	ecdsaKey, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)

	tests := []struct {
		name      string
		device    *Device
		wantError bool
	}{
		{
			name: "success - RSA device",
			device: &Device{
				Algorithm: AlgorithmRSA,
				PublicKey: &rsaKey.PublicKey,
			},
			wantError: false,
		},
		{
			name: "error - ECDSA device",
			device: &Device{
				Algorithm: AlgorithmECDSA,
				PublicKey: &ecdsaKey.PublicKey,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := tt.device.GetRSAPublicKey()

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if key == nil {
					t.Error("expected key, got nil")
				}
			}
		})
	}
}

func TestGetECDSAPublicKey(t *testing.T) {
	rsaKey, _ := rsa.GenerateKey(rand.Reader, 512)
	ecdsaKey, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)

	tests := []struct {
		name      string
		device    *Device
		wantError bool
	}{
		{
			name: "success - ECDSA device",
			device: &Device{
				Algorithm: AlgorithmECDSA,
				PublicKey: &ecdsaKey.PublicKey,
			},
			wantError: false,
		},
		{
			name: "error - RSA device",
			device: &Device{
				Algorithm: AlgorithmRSA,
				PublicKey: &rsaKey.PublicKey,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := tt.device.GetECDSAPublicKey()

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if key == nil {
					t.Error("expected key, got nil")
				}
			}
		})
	}
}
