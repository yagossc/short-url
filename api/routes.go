package api

// Routes set up the server routes.
func (s *Server) Routes() {

	// Base route
	s.e.GET("/", s.base)

	// Shortener route
	s.e.POST("/", s.shortener)
}
