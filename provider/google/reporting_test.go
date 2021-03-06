package google

import (
	"context"
	"fmt"
	"testing"

	"github.com/fupas/commons/pkg/env"
	"github.com/fupas/platform"
)

func TestGoogleErrorReporting(t *testing.T) {
	projectID := env.GetString("PROJECT_ID", "")
	serviceName := env.GetString("SERVICE_NAME", "test")

	cl, err := platform.NewClient(context.TODO(), NewErrorReporting(context.TODO(), projectID, serviceName))
	if err != nil {
		t.Error(err)
	}

	if cl.ErrorReportingImpl == nil {
		t.Error(fmt.Errorf("ErrorReporting was not initialized"))
	}

	cl.ErrorReportingImpl(fmt.Errorf("Report this: TestGoogleErrorReporting"))
	//ReportError(fmt.Errorf("Report this: TestErrorReporting"))
}
