package domain

import (
	"sync"
)

type SignatureAlgorithm string

const (
	AlgorithmRSA   SignatureAlgorithm = "RSA"
	AlgorithmECDSA SignatureAlgorithm = "ECDSA"
)

type Device struct {
	ID               string             `json:"id"`
	Algorithm        SignatureAlgorithm `json:"algorithm"`
	Label            string             `json:"label,omitempty"`
	SignatureCounter int                `json:"signature_counter"`
	PublicKey        interface{}        `json:"-"`                        // Can be *rsa.PublicKey or *ecdsa.PublicKey
	PrivateKey       interface{}        `json:"-"`                        // Can be *rsa.PrivateKey or *ecdsa.PrivateKey
	LastSignature    string             `json:"last_signature,omitempty"` // base64 encoded
	mu               sync.Mutex         `json:"-"`                        // Mutex to ensure thread-safe counter increment
}

// SignatureResponse represents the response returned after signing data
type SignatureResponse struct {
	Signature  string `json:"signature"`   // base64 encoded signature
	SignedData string `json:"signed_data"` // the secured data that was signed
}
