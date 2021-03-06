package platform

import (
	"context"
	"fmt"
	"testing"
)

func TestNewDefaultClient(t *testing.T) {
	client, err := NewDefaultClient(context.TODO())

	if err != nil {
		t.Error(err)
	}
	if client == nil {
		t.Error(fmt.Errorf("Client did not initialize"))
	}

	if client.ErrorReportingImpl == nil {
		t.Error(fmt.Errorf("ErrorReporting was not initialized"))
	}
}

func TestNewStdoutErrorReporting(t *testing.T) {
	client, err := NewDefaultClient(context.TODO())

	if err != nil {
		t.Error(err)
	}

	client.ErrorReportingImpl(fmt.Errorf("Report this: TestNewStdoutErrorReporting"))
}

func TestErrorReporting(t *testing.T) {
	ReportError(fmt.Errorf("Report this: TestErrorReporting"))
}
