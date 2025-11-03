package api

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/gin-gonic/gin"
)

// Response is the generic API response container.
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress string
	repository    *persistence.InMemoryRepository
	router        *gin.Engine
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string) *Server {
	repository := persistence.NewInMemoryRepository()

	return &Server{
		listenAddress: listenAddress,
		repository:    repository,
		router:        gin.Default(),
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	v0 := s.router.Group("/api/v0")
	{
		// Health endpoint
		v0.GET("/health", s.Health)

		// Device endpoints
		v0.POST("/devices", s.CreateDevice)
		v0.GET("/devices", s.ListDevices)
		v0.GET("/devices/:id", s.GetDevice)

		// Signature endpoint
		v0.POST("/devices/:id/sign", s.SignTransaction)
	}

	return s.router.Run(s.listenAddress)
}
