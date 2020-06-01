package api

import (
	"fmt"

	"github.com/yagossc/short-url/store"
)

// Routes set up the server routes.
func (s *Server) Routes() {

	// Shortener route
	s.e.POST("/", s.shortener)

	// History routes
	s.e.GET("/history/", s.fullHistory)
	s.e.GET("/history/day", s.entriesLastDay)
	s.e.GET("/history/week", s.entriesLastWeek)

	s.loadDynamicRoutes()
}

func (s *Server) loadDynamicRoutes() {
	results, err := store.FindAllURL(s.db)
	if err != nil { // FIXME: properly handle this error
		fmt.Printf("error: %v\n", err)
	}

	// Inject dynamic routes
	for _, val := range results {
		s.AddRoute(val.Short)
	}
}
