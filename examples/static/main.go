package main

import (
	"fmt"
	l "log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	srv "github.com/fupas/platform/pkg/http"
)

// the router instance
var svc *srv.Server

func shutdown(*echo.Echo) {
	l.Printf("Cleaning-up ...")
}

func setupRouter() *echo.Echo {
	// Create a new instance
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	// Serve static files from ./public
	e.Static("/", "public")

	return e
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	errorPage := fmt.Sprintf("public/%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}

func main() {
	svc := srv.NewServer(setupRouter, shutdown, customHTTPErrorHandler)
	svc.StartBlocking()
}
