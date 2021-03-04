package main

import (
	"fmt"
	l "log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/fupas/commons/pkg/env"
	svc "github.com/fupas/platform/pkg/http"
)

// ShutdownDelay is the delay before exiting the process
const ShutdownDelay = 10

// the router instance
var mux *echo.Echo
var staticFileLocation string = env.GetString("STATIC_FILE_LOCATION", "public")

func setup() *echo.Echo {
	// Create a new router instance
	e := echo.New()

	// add and configure the middlewares
	e.Use(middleware.Recover())
	// TODO: add configure e.Use(middleware.Gzip())

	// TODO: add/configure e.Use(middleware.Logger())
	// TODO: e.Logger.SetLevel(log.INFO)

	// add the routes last
	e.Static("/", staticFileLocation) // serve static files from e.g. ./public

	return e
}

func shutdown(*echo.Echo) {
	// TODO: implement your own stuff here

	l.Printf("Exiting now ...")
}

func init() {
	// TODO: initialize everything global here

	l.Printf("Initializing ...")
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	errorPage := fmt.Sprintf("%s/%d.html", staticFileLocation, code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}

func main() {
	service := svc.NewServer(setup, shutdown, customHTTPErrorHandler)
	service.StartBlocking()
}
