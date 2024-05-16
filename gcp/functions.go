package gcp

import (
	"errors"

	"cloud.google.com/go/compute/metadata"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog/log"
)

var namespace = errorx.NewNamespace("gcp")
var EnvVariableNotFound = errorx.NewType(namespace, "Env Variable Not Found", errorx.NotFound())

func Must(s string, err error) string {
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		panic(err)
	} else {
		return s
	}
}

func Region() (string, error) {
	region, err := metadata.Zone()
	if region == "" {
		return "", errors.Join(err, EnvVariableNotFound.New("environment variable %s not found", "REGION"))
	} else {
		return region, nil
	}
}

func ProjectId() (string, error) {
	projectId, err := metadata.ProjectID()
	if err != nil {
		return "", errors.Join(err, EnvVariableNotFound.New("env variable %s not found", "GOOGLE_CLOUD_PROJECT"))
	} else {
		return projectId, nil
	}
}
