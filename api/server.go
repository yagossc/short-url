package api

import (
	"github.com/labstack/echo/v4"
)

// Server is a wrapper around echo.Echo and the database connection
type Server struct {
	e *echo.Echo
}

// NewServer returns a new server instance
func NewServer(e *echo.Echo) *Server {
	s := &Server{e: e}

	return s
}

// Start starts the server instance.
func (s *Server) Start(address string) error {
	return s.e.Start(address)
}
