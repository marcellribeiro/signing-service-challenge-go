package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/gin-gonic/gin"
)

func setupTestServer() *Server {
	gin.SetMode(gin.TestMode)
	return NewServer(":8080")
}

func TestCreateDevice(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
	}{
		{
			name: "success - create RSA device",
			requestBody: CreateDeviceRequest{
				Algorithm: domain.AlgorithmRSA,
				Label:     "Test Device",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "success - create ECDSA device",
			requestBody: CreateDeviceRequest{
				Algorithm: domain.AlgorithmECDSA,
				Label:     "ECDSA Device",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "error - invalid algorithm",
			requestBody:    CreateDeviceRequest{Algorithm: "INVALID"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "error - missing algorithm",
			requestBody:    map[string]string{"label": "test"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v0/devices", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			server.CreateDevice(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestListDevices(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(*Server)
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "success - empty list",
			setup:          func(s *Server) {},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "success - list with devices",
			setup: func(s *Server) {
				s.repository.Create(domain.NewDevice("dev-1", domain.AlgorithmRSA, "Device 1", nil, nil))
				s.repository.Create(domain.NewDevice("dev-2", domain.AlgorithmECDSA, "Device 2", nil, nil))
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer()
			tt.setup(server)

			req := httptest.NewRequest(http.MethodGet, "/api/v0/devices", nil)
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			server.ListDevices(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response Response
			json.Unmarshal(w.Body.Bytes(), &response)
			devices := response.Data.([]interface{})
			if len(devices) != tt.expectedCount {
				t.Errorf("expected %d devices, got %d", tt.expectedCount, len(devices))
			}
		})
	}
}

func TestGetDevice(t *testing.T) {
	tests := []struct {
		name           string
		deviceID       string
		setup          func(*Server)
		expectedStatus int
	}{
		{
			name:     "success - get existing device",
			deviceID: "existing-device",
			setup: func(s *Server) {
				s.repository.Create(domain.NewDevice("existing-device", domain.AlgorithmRSA, "Test", nil, nil))
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "error - device not found",
			deviceID:       "non-existent",
			setup:          func(s *Server) {},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer()
			tt.setup(server)

			req := httptest.NewRequest(http.MethodGet, "/api/v0/devices/"+tt.deviceID, nil)
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: tt.deviceID}}

			server.GetDevice(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestSignTransaction(t *testing.T) {
	tests := []struct {
		name           string
		deviceID       string
		requestBody    interface{}
		setup          func(*Server)
		expectedStatus int
	}{
		{
			name:     "success - sign with RSA device",
			deviceID: "rsa-device",
			requestBody: SignTransactionRequest{
				Data: "test transaction",
			},
			setup: func(s *Server) {
				gen := &crypto.RSAGenerator{}
				kp, _ := gen.Generate()
				s.repository.Create(domain.NewDevice("rsa-device", domain.AlgorithmRSA, "RSA", kp.Public, kp.Private))
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "error - device not found",
			deviceID: "non-existent",
			requestBody: SignTransactionRequest{
				Data: "test",
			},
			setup:          func(s *Server) {},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "error - missing data",
			deviceID:       "any-device",
			requestBody:    map[string]string{},
			setup:          func(s *Server) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer()
			tt.setup(server)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v0/devices/"+tt.deviceID+"/sign", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: tt.deviceID}}

			server.SignTransaction(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
