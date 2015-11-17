package main

import (
	"net/http"
	"os"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

// Handler
func hello(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World! I was here\n")
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	// Routes
	e.Get("/", hello)

	// Start server
	var port = os.Getenv("OPENSHIFT_GO_PORT")
	if(port == "") {
		port = "1323"
	}
	var ip = os.Getenv("OPENSHIFT_GO_IP")
	if(ip == "") {
		ip = "localhost"
	}
	
	e.Run(ip+":"+port)
}
