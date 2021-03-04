package observer

import (
	"context"
	"log"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
)

type (
	// Client holds all clients needed to access basic Google Cloud services
	Client struct {
		ErrorClient *errorreporting.Client
		LogClient   *logging.Client
		//Logger      *logging.Logger
	}
)

// NewClient creates a new client
func NewClient(ctx context.Context, projectID, serviceName string) (*Client, error) {
	c := Client{}

	// initialize logging
	lc, err := logging.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	c.LogClient = lc

	// initialize error reporting
	ec, err := errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName: serviceName,
		OnError: func(err error) {
			log.Printf("Could not log error: %v", err)
		},
	})
	if err != nil {
		return nil, err
	}
	c.ErrorClient = ec

	return &c, nil
}

// Close closes all clients to the Google Cloud services
func (c *Client) Close() {
	if c.ErrorClient != nil {
		c.ErrorClient.Close()
	}
	if c.LogClient != nil {
		c.LogClient.Close()
	}
}

// ReportError reports an error, what else?
func (c *Client) ReportError(err error) {
	c.ErrorClient.Report(errorreporting.Entry{Error: err})
}
