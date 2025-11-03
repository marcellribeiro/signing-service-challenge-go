package persistence

import (
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

func TestInMemoryRepository_Create(t *testing.T) {
	tests := []struct {
		name      string
		device    *domain.Device
		wantError error
	}{
		{
			name: "success - create new device",
			device: &domain.Device{
				ID:        "device-1",
				Algorithm: domain.AlgorithmRSA,
				Label:     "Test Device",
			},
			wantError: nil,
		},
		{
			name: "error - duplicate device ID",
			device: &domain.Device{
				ID:        "duplicate-id",
				Algorithm: domain.AlgorithmRSA,
			},
			wantError: ErrDeviceAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInMemoryRepository()

			// For duplicate test, create first device
			if tt.wantError == ErrDeviceAlreadyExists {
				repo.Create(&domain.Device{ID: "duplicate-id"})
			}

			err := repo.Create(tt.device)

			if tt.wantError != nil {
				if err != tt.wantError {
					t.Errorf("expected error %v, got %v", tt.wantError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestInMemoryRepository_Get(t *testing.T) {
	tests := []struct {
		name      string
		deviceID  string
		setup     func(*InMemoryRepository)
		wantError error
	}{
		{
			name:     "success - get existing device",
			deviceID: "existing-device",
			setup: func(repo *InMemoryRepository) {
				repo.Create(&domain.Device{
					ID:        "existing-device",
					Algorithm: domain.AlgorithmRSA,
				})
			},
			wantError: nil,
		},
		{
			name:      "error - device not found",
			deviceID:  "non-existent",
			setup:     func(repo *InMemoryRepository) {},
			wantError: ErrDeviceNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInMemoryRepository()
			tt.setup(repo)

			device, err := repo.Get(tt.deviceID)

			if tt.wantError != nil {
				if err != tt.wantError {
					t.Errorf("expected error %v, got %v", tt.wantError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if device == nil {
					t.Error("expected device, got nil")
				}
				if device.ID != tt.deviceID {
					t.Errorf("expected device ID %q, got %q", tt.deviceID, device.ID)
				}
			}
		})
	}
}

func TestInMemoryRepository_List(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(*InMemoryRepository)
		expectedCount int
	}{
		{
			name:          "success - empty list",
			setup:         func(repo *InMemoryRepository) {},
			expectedCount: 0,
		},
		{
			name: "success - list multiple devices",
			setup: func(repo *InMemoryRepository) {
				repo.Create(&domain.Device{ID: "device-1"})
				repo.Create(&domain.Device{ID: "device-2"})
				repo.Create(&domain.Device{ID: "device-3"})
			},
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInMemoryRepository()
			tt.setup(repo)

			devices, err := repo.List()

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(devices) != tt.expectedCount {
				t.Errorf("expected %d devices, got %d", tt.expectedCount, len(devices))
			}
		})
	}
}

func TestInMemoryRepository_Update(t *testing.T) {
	tests := []struct {
		name      string
		device    *domain.Device
		setup     func(*InMemoryRepository)
		wantError error
	}{
		{
			name: "success - update existing device",
			device: &domain.Device{
				ID:               "device-1",
				SignatureCounter: 10,
			},
			setup: func(repo *InMemoryRepository) {
				repo.Create(&domain.Device{
					ID:               "device-1",
					SignatureCounter: 5,
				})
			},
			wantError: nil,
		},
		{
			name: "error - update non-existent device",
			device: &domain.Device{
				ID: "non-existent",
			},
			setup:     func(repo *InMemoryRepository) {},
			wantError: ErrDeviceNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInMemoryRepository()
			tt.setup(repo)

			err := repo.Update(tt.device)

			if tt.wantError != nil {
				if err != tt.wantError {
					t.Errorf("expected error %v, got %v", tt.wantError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				// Verify update worked
				device, _ := repo.Get(tt.device.ID)
				if device.SignatureCounter != tt.device.SignatureCounter {
					t.Errorf("expected counter %d, got %d", tt.device.SignatureCounter, device.SignatureCounter)
				}
			}
		})
	}
}
