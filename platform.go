package platform

import (
	"context"
	"fmt"
	"log"

	"github.com/fupas/platform/pkg/errorreporting"
)

const (
	DatastoreProvider     = "DATASTORE"
	StorageProvider       = "STORAGE"
	ErroReportingProvider = "ERROR"
	LoggingProvider       = "LOGGING"
	MetricsProvider       = "METRICS"
)

type (
	Client struct {
		// Datastore
		// Storage
		// ErrorReporting provider
		ErrorReportingImpl errorreporting.ErrorReportingFunc
		// Logging
		// Metrics
	}

	ProviderImplementation interface {
		ProviderType() string
		Implementation(context.Context) interface{}
	}
)

var globalClient *Client

func init() {
	// create a default instance
	cl, err := NewDefaultClient(context.TODO())
	if err != nil {
		log.Fatal(err) // no point to continue if already this fails ...
	}
	RegisterGlobally(cl)
}

// NewClient initializes a new platform client and initializes it with the given provider implemenations
func NewClient(ctx context.Context, with ...ProviderImplementation) (*Client, error) {
	if len(with) != 0 {
		cl := &Client{}

		for _, p := range with {
			if p != nil {
				t := p.ProviderType()
				switch t {
				case DatastoreProvider:
					fmt.Println(t)
					break
				case StorageProvider:
					fmt.Println(t)
					break
				case ErroReportingProvider:
					cl.ErrorReportingImpl = p.Implementation(ctx).(func(error))
					break
				case LoggingProvider:
					fmt.Println(t)
					break
				case MetricsProvider:
					fmt.Println(t)
					break
				default:
					fmt.Println("Unknown")
				}
			}
		}
		return cl, nil
	}
	return NewDefaultClient(ctx)
}

// NewDefaultClient initializes a new platform client and initializes it with default providers
func NewDefaultClient(ctx context.Context, with ...ProviderImplementation) (*Client, error) {
	return NewClient(ctx, errorreporting.NewStdoutErrorReporting())
}

// RegisterGlobally replaces the current global platform client with a new one
func RegisterGlobally(cl *Client) *Client {
	old := globalClient
	globalClient = cl
	return old
}

func ReportError(err error) {
	globalClient.ErrorReportingImpl(err)
}
