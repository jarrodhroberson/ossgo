package storage

import (
	"errors"
	"fmt"
	"net/http"
	"path"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/gin-gonic/gin"
	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

var expectedCloudEventError = errors.New("expected CloudEvent")
var expectedCloudStorageEventError = errors.New("Bad Request: expected Cloud Storage event")

func failedToParseCloudEventError(e *event.Event) error {
	return fmt.Errorf("failed to parse CloudEvent: %v", e)
}

func CloudStorageEvent(c *gin.Context) {
	ce, err := cloudevents.NewEventFromHTTPRequest(c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.Join(err, expectedCloudEventError))
		return
	}

	var so storagedata.StorageObjectData
	err = protojson.Unmarshal(ce.Data(), &so)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.Join(err, expectedCloudStorageEventError, failedToParseCloudEventError(ce)))
		return
	}

	log.Info().Msgf("Cloud Storage object changed: %s updated at %s", path.Join(so.GetBucket(), so.GetName()), so.Updated.AsTime().UTC())
}
