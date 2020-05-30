package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yagossc/short-url/app"
)

func (s *Server) base(c echo.Context) error {
	var response struct{}
	return c.JSON(http.StatusOK, response)
}

func (s *Server) shortener(c echo.Context) error {
	var url app.MapURL
	if err := c.Bind(&url); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, url)
}
