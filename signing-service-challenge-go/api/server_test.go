package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealth(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "success - health check returns pass",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer()

			req := httptest.NewRequest(http.MethodGet, "/api/v0/health", nil)
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			server.Health(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response Response
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			healthData := response.Data.(map[string]interface{})
			if healthData["status"] != "pass" {
				t.Errorf("expected status 'pass', got %v", healthData["status"])
			}
			if healthData["version"] != "v0" {
				t.Errorf("expected version 'v0', got %v", healthData["version"])
			}
		})
	}
}
