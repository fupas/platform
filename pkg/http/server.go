package http

import (
	"context"
	"fmt"
	l "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fupas/commons/pkg/env"
	"github.com/labstack/echo/v4"
)

const (
	// ShutdownDelay is the time to wait for all request, go-routines etc complete
	ShutdownDelay = 10 // seconds
)

type (
	// RouterFunc creates a mux
	RouterFunc func() *echo.Echo
	// ShutdownFunc is called before the server stops
	ShutdownFunc func(*echo.Echo)

	// Server is an interface for the HTTP server
	Server interface {
		StartBlocking()
		Stop()
	}

	server struct {
		mux              *echo.Echo
		shutdown         ShutdownFunc
		errorHandlerImpl echo.HTTPErrorHandler
	}
)

// NewServer returns a new HTTP server
func NewServer(router RouterFunc, shutdown ShutdownFunc, errorHandler echo.HTTPErrorHandler) Server {
	return &server{
		mux:              router(),
		shutdown:         shutdown,
		errorHandlerImpl: errorHandler,
	}
}

// Stop forces a shutdown
func (s *server) Stop() {
	// all the implementation specific shoutdown code to clean-up
	s.shutdown(s.mux)

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownDelay*time.Second)
	defer cancel()
	if err := s.mux.Shutdown(ctx); err != nil {
		s.mux.Logger.Fatal(err)
		fmt.Println(err)
	}
	l.Printf("Exiting now ...")
}

// StartBlocking starts a new server in the main process
func (s *server) StartBlocking() {
	// setup shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		s.Stop()
	}()

	// add the central error handler
	if s.errorHandlerImpl != nil {
		s.mux.HTTPErrorHandler = s.errorHandlerImpl
	}

	// start the server
	if err := s.mux.Start(env.GetString("PORT", ":8080")); err != nil && err != http.ErrServerClosed {
		s.mux.Logger.Fatal("Error shutting down the server")
		fmt.Println(err)
	}
}
