package api

// Routes set up the server routes.
func (s *Server) Routes() {

	// Shortener route
	s.e.POST("/", s.shortener)
}
