package tasks

import (
	"fmt"

	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"

	"github.com/jarrodhroberson/ossgo/functions/must"
	"github.com/jarrodhroberson/ossgo/gcp"
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
