package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yagossc/short-url/app"
	"github.com/yagossc/short-url/history"
	"github.com/yagossc/short-url/store"
)

// FIXME: these functions are a little repetitive, remember to 'DRY'.
func (s *Server) entriesLastDay(c echo.Context) error {
	var short app.Short

	if err := c.Bind(&short); err != nil {
		return err
	}

	if short.URL == "" {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}

	splitted := strings.Split(short.URL, "/")
	shortID := splitted[len(splitted)-1]

	entries, err := store.FindAllReqByShort(s.db, shortID)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}

	ocurrences := history.GetEntriesInInvertval(entries, 24)

	return c.JSON(http.StatusOK, ocurrences)
}

func (s *Server) entriesLastWeek(c echo.Context) error {
	var short app.Short

	if err := c.Bind(&short); err != nil {
		return err
	}

	if short.URL == "" {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}

	splitted := strings.Split(short.URL, "/")
	shortID := splitted[len(splitted)-1]

	entries, err := store.FindAllReqByShort(s.db, shortID)
	if err != nil {
		return err
	}

	ocurrences := history.GetEntriesInInvertval(entries, 24*7)

	return c.JSON(http.StatusOK, ocurrences)
}

func (s *Server) fullHistory(c echo.Context) error {
	var short app.Short

	if err := c.Bind(&short); err != nil {
		return err
	}

	if short.URL == "" {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}

	splitted := strings.Split(short.URL, "/")
	shortID := splitted[len(splitted)-1]

	entries, err := store.FindAllReqByShort(s.db, shortID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, len(entries))
}
