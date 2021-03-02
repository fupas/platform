package services

import (
	"context"
	"encoding/json"
	"fmt"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"

	"github.com/fupas/commons/pkg/env"
	"github.com/fupas/observer/pkg/observer"
)

/*

For an example of a handler see

https://cloud.google.com/tasks/docs/creating-appengine-handlers

and

https://github.com/GoogleCloudPlatform/golang-samples/blob/master/appengine/go11x/tasks/handle_task/handle_task.go

*/

// CreateTask is used to schedule a background task using the default queue.
// The payload can be any struct and will be marshalled into a json string.
func CreateTask(ctx context.Context, handler string, payload interface{}) (*taskspb.Task, error) {

	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		observer.ReportError(err)
		return nil, err
	}
	defer client.Close()

	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", env.GetString("PROJECT_ID", ""), env.GetString("LOCATION_ID", ""), env.GetString("DEFAULT_QUEUE", ""))

	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: &taskspb.AppEngineHttpRequest{
					HttpMethod:  taskspb.HttpMethod_POST,
					RelativeUri: handler,
				},
			},
		},
	}

	if payload != nil {
		// marshal the payload
		b, err := json.Marshal(payload)
		if err != nil {
			observer.ReportError(err)
			return nil, err
		}
		req.Task.GetAppEngineHttpRequest().Body = b
	}

	task, err := client.CreateTask(ctx, req)
	if err != nil {
		observer.ReportError(err)
		return nil, err
	}

	return task, nil
}

// CreateSimpleTask is used to schedule a background task using the default queue.
// The payload is a simple string, i.e. no marshalling happens.
func CreateSimpleTask(ctx context.Context, handler, payload string) (*taskspb.Task, error) {

	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		observer.ReportError(err)
		return nil, err
	}
	defer client.Close()

	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", env.GetString("PROJECT_ID", ""), env.GetString("LOCATION_ID", ""), env.GetString("DEFAULT_QUEUE", ""))

	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: &taskspb.AppEngineHttpRequest{
					HttpMethod:  taskspb.HttpMethod_POST,
					RelativeUri: handler,
				},
			},
		},
	}

	req.Task.GetAppEngineHttpRequest().Body = []byte(payload)

	task, err := client.CreateTask(ctx, req)
	if err != nil {
		observer.ReportError(err)
		return nil, err
	}

	return task, nil
}
