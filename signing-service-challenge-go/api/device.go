package api

import (
	"encoding/base64"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateDeviceRequest represents the request body for creating a device
type CreateDeviceRequest struct {
	ID        string                    `json:"id,omitempty"`
	Algorithm domain.SignatureAlgorithm `json:"algorithm" binding:"required"`
	Label     string                    `json:"label,omitempty"`
}

// CreateDeviceResponse represents the response after creating a device
type CreateDeviceResponse struct {
	ID               string                    `json:"id"`
	Algorithm        domain.SignatureAlgorithm `json:"algorithm"`
	Label            string                    `json:"label,omitempty"`
	SignatureCounter int                       `json:"signature_counter"`
}

// SignTransactionRequest represents the request body for signing a transaction
type SignTransactionRequest struct {
	Data string `json:"data" binding:"required"`
}

// CreateDevice creates a new signature device
func (s *Server) CreateDevice(c *gin.Context) {
	var req CreateDeviceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors: []string{"Invalid request body: " + err.Error()},
		})
		return
	}

	// Validate algorithm
	if req.Algorithm != domain.AlgorithmRSA && req.Algorithm != domain.AlgorithmECDSA {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors: []string{"Algorithm must be either 'RSA' or 'ECDSA'"},
		})
		return
	}

	// Generate ID if not provided
	deviceID := req.ID
	if deviceID == "" {
		deviceID = uuid.New().String()
	}

	// Generate key pair based on algorithm
	var publicKey, privateKey interface{}
	var err error

	if req.Algorithm == domain.AlgorithmRSA {
		generator := &crypto.RSAGenerator{}
		keyPair, genErr := generator.Generate()
		if genErr != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Errors: []string{"Failed to generate RSA key pair: " + genErr.Error()},
			})
			return
		}
		publicKey = keyPair.Public
		privateKey = keyPair.Private
	} else {
		generator := &crypto.ECCGenerator{}
		keyPair, genErr := generator.Generate()
		if genErr != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Errors: []string{"Failed to generate ECDSA key pair: " + genErr.Error()},
			})
			return
		}
		publicKey = keyPair.Public
		privateKey = keyPair.Private
	}

	// Create device
	device := domain.NewDevice(deviceID, req.Algorithm, req.Label, publicKey, privateKey)

	// Store device
	if err = s.repository.Create(device); err != nil {
		if err == persistence.ErrDeviceAlreadyExists {
			c.JSON(http.StatusConflict, ErrorResponse{
				Errors: []string{"Device with this ID already exists"},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Errors: []string{"Failed to store device: " + err.Error()},
		})
		return
	}

	response := CreateDeviceResponse{
		ID:               device.ID,
		Algorithm:        device.Algorithm,
		Label:            device.Label,
		SignatureCounter: device.SignatureCounter,
	}

	c.JSON(http.StatusCreated, Response{Data: response})
}

// ListDevices returns all signature devices
func (s *Server) ListDevices(c *gin.Context) {
	devices, err := s.repository.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Errors: []string{"Failed to list devices: " + err.Error()},
		})
		return
	}

	response := make([]CreateDeviceResponse, len(devices))
	for i, device := range devices {
		response[i] = CreateDeviceResponse{
			ID:               device.ID,
			Algorithm:        device.Algorithm,
			Label:            device.Label,
			SignatureCounter: device.SignatureCounter,
		}
	}

	c.JSON(http.StatusOK, Response{Data: response})
}

// GetDevice returns a single signature device by ID
func (s *Server) GetDevice(c *gin.Context) {
	id := c.Param("id")

	device, err := s.repository.Get(id)
	if err != nil {
		if err == persistence.ErrDeviceNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Errors: []string{"Device not found"},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Errors: []string{"Failed to get device: " + err.Error()},
		})
		return
	}

	response := CreateDeviceResponse{
		ID:               device.ID,
		Algorithm:        device.Algorithm,
		Label:            device.Label,
		SignatureCounter: device.SignatureCounter,
	}

	c.JSON(http.StatusOK, Response{Data: response})
}

// SignTransaction signs transaction data with the specified device
func (s *Server) SignTransaction(c *gin.Context) {
	id := c.Param("id")

	var req SignTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors: []string{"Invalid request body: " + err.Error()},
		})
		return
	}

	// Get device
	device, err := s.repository.Get(id)
	if err != nil {
		if err == persistence.ErrDeviceNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Errors: []string{"Device not found"},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Errors: []string{"Failed to get device: " + err.Error()},
		})
		return
	}

	// Build secured data to sign
	securedData := device.GetSecuredDataToSign(req.Data)

	// Create appropriate signer
	var signer crypto.Signer
	if device.Algorithm == domain.AlgorithmRSA {
		privateKey, keyErr := device.GetRSAPrivateKey()
		if keyErr != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Errors: []string{"Failed to get RSA private key: " + keyErr.Error()},
			})
			return
		}
		signer = crypto.NewRSASigner(privateKey)
	} else {
		privateKey, keyErr := device.GetECDSAPrivateKey()
		if keyErr != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Errors: []string{"Failed to get ECDSA private key: " + keyErr.Error()},
			})
			return
		}
		signer = crypto.NewECDSASigner(privateKey)
	}

	// Sign the data
	signature, err := signer.Sign([]byte(securedData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Errors: []string{"Failed to sign data: " + err.Error()},
		})
		return
	}

	// Encode signature to base64
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	// Update device with new signature and increment counter
	device.IncrementCounter(signatureBase64)

	// Persist updated device
	if err = s.repository.Update(device); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Errors: []string{"Failed to update device: " + err.Error()},
		})
		return
	}

	response := domain.SignatureResponse{
		Signature:  signatureBase64,
		SignedData: securedData,
	}

	c.JSON(http.StatusOK, Response{Data: response})
}
