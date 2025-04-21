package tasks

import (
	"context"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"

	"github.com/rs/zerolog/log"
)

var client *cloudtasks.Client

func init() {
	c, err := cloudtasks.NewClient(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("could not create cloudtasks client")
	}
	client = c
}

type CreateTaskRequestOption func(ctr *cloudtaskspb.CreateTaskRequest)

type QueuePathProvider func() string
