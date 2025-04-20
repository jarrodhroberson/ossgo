package tasks

import (
	"fmt"

	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
)

func DisallowDuplicates(ctr *cloudtaskspb.CreateTaskRequest, path string, name string) CreateTaskRequestOption {
	return func(ctr *cloudtaskspb.CreateTaskRequest) {
		// Task name must be formatted: "projects/<PROJECT_ID>/locations/<LOCATION_ID>/queues/<QUEUE_ID>/tasks/<TASK_ID>"
		ctr.Task.Name = fmt.Sprintf("%s/tasks/%s", path, ctr.Task.Name)
	}
}
