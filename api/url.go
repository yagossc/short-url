package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yagossc/short-url/app"
	"github.com/yagossc/short-url/shortener"
	"github.com/yagossc/short-url/store"
)

func (s *Server) redirect(c echo.Context) error {
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

	return c.Redirect(http.StatusMovedPermanently, response.Long)
}

func (s *Server) shortener(c echo.Context) error {
	var long app.Long

	if err := c.Bind(&long); err != nil {
		return err
	}

	if long.URL == "" {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}

	short := shortener.GetShortURL(time.Now().UnixNano())

	var mappedURL app.MapURL
	mappedURL.Short = short
	mappedURL.Long = long.URL

	shortened, err := store.InsertURL(s.db, &mappedURL)
	if err != nil {
		// fmt.Printf("%v\n", err)
		return err
	}

	s.AddRoute(shortened)

	return c.JSON(http.StatusCreated, s.url+"/"+shortened)
}
