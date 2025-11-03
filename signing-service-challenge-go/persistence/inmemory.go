package persistence

import (
	"errors"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

var (
	ErrDeviceNotFound      = errors.New("device not found")
	ErrDeviceAlreadyExists = errors.New("device already exists")
)

// InMemoryRepository implements an in-memory storage for signature devices
type InMemoryRepository struct {
	devices map[string]*domain.Device
	mu      sync.RWMutex
}

// NewInMemoryRepository creates a new in-memory repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		devices: make(map[string]*domain.Device),
	}
}

// Create stores a new device
func (r *InMemoryRepository) Create(device *domain.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.devices[device.ID]; exists {
		return ErrDeviceAlreadyExists
	}

	r.devices[device.ID] = device
	return nil
}

// Get retrieves a device by ID
func (r *InMemoryRepository) Get(id string) (*domain.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	device, exists := r.devices[id]
	if !exists {
		return nil, ErrDeviceNotFound
	}

	return device, nil
}

// List returns all devices
func (r *InMemoryRepository) List() ([]*domain.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	devices := make([]*domain.Device, 0, len(r.devices))
	for _, device := range r.devices {
		devices = append(devices, device)
	}

	return devices, nil
}

// Update updates an existing device
func (r *InMemoryRepository) Update(device *domain.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.devices[device.ID]; !exists {
		return ErrDeviceNotFound
	}

	r.devices[device.ID] = device
	return nil
}
