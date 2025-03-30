package gcp

import (
	"os"
	"sync"
)

var env *environment
var once sync.Once

type environment map[string]string

type Environment interface {
	Application() string
	GinMode() string
	DeploymentId() string
	Env() string
	Instance() string
	MemoryMb() string
	Runtime() string
	Service() string
	Version() string
	CloudProject() string
	Port() string
}

func NewEnvironment() Environment {
	once.Do(func() {
		env = &environment{
			"is_cloud_run_function": os.Getenv("K_SERVICE"),
			"gin_mode":              os.Getenv("GIN_MODE"),
			"gae_application":       os.Getenv("GAE_APPLICATION"),
			"gae_deployment_id":     os.Getenv("GAE_DEPLOYMENT_ID"),
			"gae_env":               os.Getenv("GAE_ENV"),
			"gae_instance":          os.Getenv("GAE_INSTANCE"),
			"gae_memory_mb":         os.Getenv("GAE_MEMORY_MD"),
			"gae_runtime":           os.Getenv("GAE_RUNTIME"),
			"gae_service":           os.Getenv("GAE_SERVICE"),
			"gae_version":           os.Getenv("GAE_VERSION"),
			"google_cloud_project":  os.Getenv("GOOGLE_CLOUD_PROJECT"),
			"port":                 os.Getenv("PORT"),
		}
	})
	return env
}

func (e *environment) IsCloudRunFunction() bool {
	return (*e)["is_cloud_run_function"] != ""
}

func (e *environment) Application() string {
	return (*e)["gae_application"]
}

func (env *environment) GinMode() string {
	return (*env)["gin_mode"]
}

func (e *environment) DeploymentId() string {
	return (*e)["gae_deployment_id"]
}

func (e *environment) Env() string {
	return (*e)["gae_env"]
}

func (e *environment) Instance() string {
	return (*e)["gae_instance"]
}

func (e *environment) MemoryMb() string {
	return (*e)["gae_memory_mb"]
}

func (e *environment) Runtime() string {
	return (*e)["gae_runtime"]
}

func (e *environment) Service() string {
	return (*e)["gae_service"]
}

func (e *environment) Version() string {
	return (*e)["gae_version"]
}

func (e *environment) CloudProject() string {
	return (*e)["google_cloud_project"]
}

func (e *environment) Port() string {
	return (*e)["port"]
}
