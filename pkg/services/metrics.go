package services

import (
	"context"

	"cloud.google.com/go/datastore"

	"github.com/fupas/commons/pkg/util"
	"github.com/fupas/platform/pkg/platform"
)

const (
	// DatastoreMetrics collection of metrics
	DatastoreMetrics string = "METRICS"
	// CountType is a int counter
	CountType = "COUNTER"
)

type (
	// Metric is a generic data structure to store metrics
	Metric struct {
		Name    string // unique name of the metric
		Label   string // additional context, e.g. an id, name/value pairs, comma separated labels etc
		Type    string // the type, e.g. count,
		Created int64
	}

	// Counter is a metric to collect integer values
	Counter struct {
		Metric
		Value int64
	}
)

// Count records a numeric counter value
func Count(ctx context.Context, name, label string, value int) error {
	m := Counter{}
	m.Name = name
	m.Label = label
	m.Type = CountType
	m.Created = util.Timestamp()
	m.Value = int64(value)

	key := datastore.IncompleteKey(DatastoreMetrics, nil)
	_, err := platform.DataStore().Put(ctx, key, &m)

	return err
}
