package domain

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
)

// NewDevice creates a new signature device
func NewDevice(id string, algorithm SignatureAlgorithm, label string, publicKey, privateKey interface{}) *Device {
	return &Device{
		ID:               id,
		Algorithm:        algorithm,
		Label:            label,
		SignatureCounter: 0,
		PublicKey:        publicKey,
		PrivateKey:       privateKey,
		LastSignature:    "",
	}
}

// GetSecuredDataToSign builds the data string to be signed according to the format:
// <signature_counter>_<data_to_be_signed>_<last_signature_base64_encoded>
func (d *Device) GetSecuredDataToSign(dataToBeSigned string) string {
	d.mu.Lock()
	defer d.mu.Unlock()

	lastSig := d.LastSignature
	if lastSig == "" {
		// Base case: use base64 encoded device ID
		lastSig = base64.StdEncoding.EncodeToString([]byte(d.ID))
	}

	return fmt.Sprintf("%d_%s_%s", d.SignatureCounter, dataToBeSigned, lastSig)
}

// IncrementCounter increments the signature counter and updates the last signature
// This method is thread-safe
func (d *Device) IncrementCounter(newSignature string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.SignatureCounter++
	d.LastSignature = newSignature
}

// GetRSAPrivateKey returns the private key as *rsa.PrivateKey
func (d *Device) GetRSAPrivateKey() (*rsa.PrivateKey, error) {
	if d.Algorithm != AlgorithmRSA {
		return nil, fmt.Errorf("device algorithm is not RSA")
	}
	key, ok := d.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not of type *rsa.PrivateKey")
	}
	return key, nil
}

// GetECDSAPrivateKey returns the private key as *ecdsa.PrivateKey
func (d *Device) GetECDSAPrivateKey() (*ecdsa.PrivateKey, error) {
	if d.Algorithm != AlgorithmECDSA {
		return nil, fmt.Errorf("device algorithm is not ECDSA")
	}
	key, ok := d.PrivateKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not of type *ecdsa.PrivateKey")
	}
	return key, nil
}

// GetRSAPublicKey returns the public key as *rsa.PublicKey
func (d *Device) GetRSAPublicKey() (*rsa.PublicKey, error) {
	if d.Algorithm != AlgorithmRSA {
		return nil, fmt.Errorf("device algorithm is not RSA")
	}
	key, ok := d.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not of type *rsa.PublicKey")
	}
	return key, nil
}

// GetECDSAPublicKey returns the public key as *ecdsa.PublicKey
func (d *Device) GetECDSAPublicKey() (*ecdsa.PublicKey, error) {
	if d.Algorithm != AlgorithmECDSA {
		return nil, fmt.Errorf("device algorithm is not ECDSA")
	}
	key, ok := d.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not of type *ecdsa.PublicKey")
	}
	return key, nil
}
