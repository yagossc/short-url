package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yagossc/short-url/app"
	"github.com/yagossc/short-url/store"
)

func (s *Server) base(c echo.Context) error {
	currentPath := strings.TrimPrefix(c.Path(), "/")
	// fmt.Printf("Path:%s\n", currentPath)

	response, err := store.FindURLByShort(s.db, currentPath)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	if response == nil {
		return c.JSON(http.StatusNotFound, "Not Found")
	}

	return c.JSON(http.StatusOK, response)
}

func (s *Server) shortener(c echo.Context) error {
	var url app.LongURL

	if err := c.Bind(&url); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, url)
}
