package api

// Routes set up the server routes.
func (s *Server) Routes() {

	// Shortener route
	s.e.POST("/", s.shortener)

	// History routes
	s.e.GET("/history/", s.fullHistory)
	s.e.GET("/history/day", s.entriesLastDay)
	s.e.GET("/history/week", s.entriesLastWeek)
}
