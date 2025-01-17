package gcp

import (
	"context"
	"errors"
	"os"
	"strings"

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
	region, err := metadata.GetWithContext(context.Background(), "instance/region")
	if region == "" {
		return "", errors.Join(err, EnvVariableNotFound.New("environment variable %s not found", "REGION"))
	} else {
		parts := strings.Split(region, "/")
		region = parts[len(parts)-1]
		return region, nil
	}
}

func ProjectId() (string, error) {
	projectId, err := metadata.ProjectIDWithContext(context.Background())
	if err != nil {
		return "", errors.Join(err, EnvVariableNotFound.New("env variable %s not found", "GOOGLE_CLOUD_PROJECT"))
	} else {
		return projectId, nil
	}
}

func NewEnvironment() *Environment {
	once.Do(func() {
		environment = &Environment{
			gin_mode:             os.Getenv("GIN_MODE"),
			gae_application:      os.Getenv("GAE_APPLICATION"),
			gae_deployment_id:    os.Getenv("GAE_DEPLOYMENT_ID"),
			gae_env:              os.Getenv("GAE_ENV"),
			gae_instance:         os.Getenv("GAE_INSTANCE"),
			gae_memory_mb:        os.Getenv("GAE_MEMORY_MD"),
			gae_runtime:          os.Getenv("GAE_RUNTIME"),
			gae_service:          os.Getenv("GAE_SERVICE"),
			gae_version:          os.Getenv("GAE_VERSION"),
			google_cloud_project: os.Getenv("GOOGLE_CLOUD_PROJECT"),
			port:                 os.Getenv("PORT"),
		}
	})
	return environment
}
