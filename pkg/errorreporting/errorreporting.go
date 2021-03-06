package errorreporting

import (
	"context"
	"fmt"
)

type (
	ErrorReportingFunc func(error)

	StdoutErrorReporting struct{}
)

func NewStdoutErrorReporting() *StdoutErrorReporting {
	return &StdoutErrorReporting{}
}

func (c *StdoutErrorReporting) ProviderType() string {
	return "ERROR"
}

func (c *StdoutErrorReporting) Implementation(context.Context) interface{} {
	return ReportErrorToStdout
}

func ReportErrorToStdout(e error) {
	fmt.Println(e)
}
