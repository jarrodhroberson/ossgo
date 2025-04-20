package tasks

import (
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
)

type CreateTaskRequestOption func(ctr *cloudtaskspb.CreateTaskRequest)

type QueuePathProvider func() string