package services

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/fupas/platform/pkg/platform"
)

const (
	// DatastoreJobs collection of job / cron metadata
	DatastoreJobs string = "JOBS"
)

type (
	// Job is the datastructure to store when a job was last run
	Job struct {
		Name    string
		Count   int
		LastRun int64
	}
)

// GetJobTimestamp returns the timestamp when a job was last executed
func GetJobTimestamp(ctx context.Context, name string) int64 {
	key := datastore.NameKey(DatastoreJobs, name, nil)

	var job = Job{}
	err := platform.DataStore().Get(ctx, key, &job)
	if err != nil {
		return 0
	}

	return job.LastRun
}

// UpdateJob updates the timestamp and count of the job metadata
func UpdateJob(ctx context.Context, name string, ts int64) error {
	key := datastore.NameKey(DatastoreJobs, name, nil)

	var job = Job{}
	err := platform.DataStore().Get(ctx, key, &job)
	if err != nil {
		job = Job{
			Name:    name,
			Count:   1,
			LastRun: ts,
		}
	} else {
		job.Count = job.Count + 1
		job.LastRun = ts
	}
	_, err = platform.DataStore().Put(ctx, key, &job)

	return err
}

/*

// ScheduleJob creates a background job
func ScheduleJob(ctx context.Context, q, req string) {
	t := taskqueue.NewPOSTTask(req, nil)
	_, err := taskqueue.Add(ctx, t, q)
	if err != nil {
		logger.Error(ctx, "schedule.jobs", err.Error())
	}
}

*/
