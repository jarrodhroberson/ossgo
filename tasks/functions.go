package tasks

import (
	"context"
	"fmt"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"

	errs "github.com/jarrodhroberson/ossgo/errors"
	"github.com/jarrodhroberson/ossgo/functions/must"
	"github.com/jarrodhroberson/ossgo/gcp"
	"github.com/rs/zerolog/log"
)

func CreateQueuePath(queueId string) string {
	projectId := must.Must(gcp.ProjectId())
	locationId := must.Must(gcp.Region())
	return fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectId, locationId, queueId)
}

func DisallowDuplicates(path string, name string) CreateTaskRequestOption {
	return func(ctr *cloudtaskspb.CreateTaskRequest) {
		// Task name must be formatted: "projects/<PROJECT_ID>/locations/<LOCATION_ID>/queues/<QUEUE_ID>/tasks/<TASK_ID>"
		ctr.Task.Name = fmt.Sprintf("%s/tasks/%s", path, name)
	}
}

func CreateTask(ctx context.Context, req *cloudtaskspb.CreateTaskRequest) (*cloudtaskspb.Task, error) {
	c, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer func(c *cloudtasks.Client) {
		err = c.Close()
		if err != nil {
			log.Error().Stack().Err(err).Msg(err.Error())
		}
	}(c)

	createdTask, err := client.CreateTask(ctx, req)
	if err != nil {
		if err.Error() == "AlreadyExists" {
			return nil, errs.NotCreatedError.Wrap(err, "task already exists %s", req.Task.Name)
		} else {
			return nil, fmt.Errorf("cloudtasks.CreateTask: %v", err)
		}
	}
	return createdTask, nil
}

func WaitForTaskCompletion(ctx context.Context, client *cloudtasks.Client, task *cloudtaskspb.Task) error {
	const maxRetries = 3                // Maximum number of retries.  Adjust as appropriate.
	const retryDelay = 10 * time.Second // Delay between retries. Adjust as appropriate.
	retries := 0

	for retries < maxRetries {
		retries++

		getTaskRequest := &cloudtaskspb.GetTaskRequest{
			Name:         task.Name,
			ResponseView: cloudtaskspb.Task_BASIC,
		}

		t, err := client.GetTask(ctx, getTaskRequest)
		if err != nil {
			return err // Returning the error directly
		}

		if t.LastAttempt.ResponseTime != nil && t.LastAttempt.ResponseStatus.Code == int32(200) {
			return nil // Task completed successfully
		}

		time.Sleep(retryDelay)
	}

	return fmt.Errorf("timed out waiting for task completion after %d retries", maxRetries)
}
