package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yagossc/short-url/api"
)

func main() {
	fmt.Printf("URL shortener\n")

	// server configuration
	e := echo.New()
	e.Use(middleware.Recover())

	// Create server
	s := api.NewServer(e)

	// API routes
	s.Routes()

	log.Fatal(s.Start(":8080"))
}
