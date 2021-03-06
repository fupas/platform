package google

import (
	"context"
	"log"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
)

type (
	ErrorReportingClient struct {
		ErrorClient *errorreporting.Client
		LogClient   *logging.Client
	}
)

func NewErrorReporting(ctx context.Context, projectID, serviceName string) *ErrorReportingClient {
	c := ErrorReportingClient{}

	// initialize logging
	lc, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
	c.ErrorClient = ec

	return &c
}

func (c *ErrorReportingClient) ProviderType() string {
	return "ERROR"
}

func (c *ErrorReportingClient) Implementation(context.Context) interface{} {
	return c.ReportError
}

// ReportError reports an error, what else?
func (c *ErrorReportingClient) ReportError(e error) {
	c.ErrorClient.Report(errorreporting.Entry{Error: e})
}
