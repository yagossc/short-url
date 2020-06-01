package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yagossc/short-url/query"
)

// Server is a wrapper around echo.Echo and the database connection
type Server struct {
	db  *query.Executor
	e   *echo.Echo
	url string
}

// NewServer returns a new server instance
func NewServer(db *query.Executor, e *echo.Echo, url string) *Server {
	s := &Server{db: db, e: e, url: url}
	return s
}

// Start starts the server instance.
func (s *Server) Start(address string) error {
	return s.e.Start(address)
}

// ServeHTTP implements `http.Handler` interface, which serves HTTP requests.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.e.ServeHTTP(w, r)
}

// AddRoute does the dynamic route injection and is what
// gives the API the expected url shortener behavior.
func (s *Server) AddRoute(id string) {
	newRoute := "/" + id
	s.e.GET(newRoute, s.redirect)
}
