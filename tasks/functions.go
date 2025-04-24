package tasks

import (
	"context"
	"fmt"

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
