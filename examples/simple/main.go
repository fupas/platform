package main

import (
	"fmt"
	l "log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	srv "github.com/fupas/platform/pkg/http"
)

// StatusResponse is the response struct that will be serialized and sent back
type StatusResponse struct {
	Status string `json:"status"`
	User   string `json:"user"`
}

// the router instance
var svc *srv.Server

func shutdown(*echo.Echo) {
	l.Printf("Cleaning-up ...")
}

// UserGetHandler is a request handler
func UserGetHandler(e echo.Context) error {
	// Create response object
	body := &StatusResponse{
		Status: "Hello world from echo!",
		User:   e.Param("user"),
	}

	return e.JSON(http.StatusOK, body)
}

// EchoRequestHandler returns request details
func EchoRequestHandler(e echo.Context) error {

	remoteAddr := e.Request().RemoteAddr
	userAgent := e.Request().UserAgent()

	dh := e.Request().Host
	dl := dh + e.Request().RequestURI
	dp := e.Request().URL.Path

	s := fmt.Sprintf("RemoteAdd='%s', UserAgent='%s', dl='%s', dh='%s', dp='%s'", remoteAddr, userAgent, dl, dh, dp)

	//analytics.TrackEvent(c.Request, "client", "category", "action", "label", 1)

	return e.String(http.StatusOK, s)
}

func setupRouter() *echo.Echo {
	// Create a new instance
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	// Add endpoint route for /users/<username>
	e.GET("/users/:user", UserGetHandler)

	// return the request parameters
	e.GET("/echo", EchoRequestHandler)

	return e
}

func main() {
	svc := srv.NewServer(setupRouter, shutdown, nil)
	svc.StartBlocking()
}
