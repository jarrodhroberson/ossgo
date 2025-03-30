package gcp

import (
	"encoding/json"
	"os"
	"sync"
)

var env *environment
var once sync.Once

func NewEnvironment() Environment {
	once.Do(func() {
		env = &environment{
			isCloudRunFunction: os.Getenv("K_SERVICE") != "",
			ginMode:              os.Getenv("GIN_MODE"),
			gaeApplication:       os.Getenv("GAE_APPLICATION"),
			gaeDeploymentId:     os.Getenv("GAE_DEPLOYMENT_ID"),
			gaeEnv:               os.Getenv("GAE_ENV"),
			gaeInstance:          os.Getenv("GAE_INSTANCE"),
			gaeMemoryMb:         os.Getenv("GAE_MEMORY_MD"),
			gaeRuntime:           os.Getenv("GAE_RUNTIME"),
			gaeService:           os.Getenv("GAE_SERVICE"),
			gaeVersion:           os.Getenv("GAE_VERSION"),
			googleCloudProject:  os.Getenv("GOOGLE_CLOUD_PROJECT"),
			port:                 os.Getenv("PORT"),
		}
	})
	return env
}

type environment struct {
	isCloudRunFunction bool
	ginMode              string
	gaeApplication       string
	gaeDeploymentId     string
	gaeEnv               string
	gaeInstance          string
	gaeMemoryMb         string
	gaeRuntime           string
	gaeService           string
	gaeVersion           string
	googleCloudProject  string
	port                 string
}

func (e *environment) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		IsCloudRunFunction bool   `json:"is_cloud_run_function"`
		GinMode            string `json:"gin_mode"`
		GaeApplication     string `json:"gae_application"`
		GaeDeploymentId   string `json:"gae_deployment_id"`
		GaeEnv             string `json:"gae_env"`
		GaeInstance        string `json:"gae_instance"`
		GaeMemoryMb        string `json:"gae_memory_mb"`
		GaeRuntime         string `json:"gae_runtime"`
		GaeService         string `json:"gae_service"`
		GaeVersion         string `json:"gae_version"`
		GoogleCloudProject string `json:"google_cloud_project"`
		Port               string `json:"port"`
	}{
		IsCloudRunFunction: e.isCloudRunFunction,
		GinMode:            e.ginMode,
		GaeApplication:     e.gaeApplication,
		GaeDeploymentId:   e.gaeDeploymentId,
		GaeEnv:             e.gaeEnv,
		GaeInstance:        e.gaeInstance,
		GaeMemoryMb:        e.gaeMemoryMb,
		GaeRuntime:         e.gaeRuntime,
		GaeService:         e.gaeService,
		GaeVersion:         e.gaeVersion,
		GoogleCloudProject: e.googleCloudProject,
		Port:               e.port,
	})
}

func (e *environment) IsCloudRunFunction() bool {
	return e.isCloudRunFunction
}

func (e *environment) Application() string {
	return e.gaeApplication
}

func (env *environment) GinMode() string {
	return env.ginMode
}

func (e *environment) DeploymentId() string {
	return e.gaeDeploymentId
}

func (e *environment) Env() string {
	return e.gaeEnv
}

func (e *environment) Instance() string {
	return e.gaeInstance
}

func (e *environment) MemoryMb() string {
	return e.gaeMemoryMb
}

func (e *environment) Runtime() string {
	return e.gaeRuntime
}

func (e *environment) Service() string {
	return e.gaeService
}

func (e *environment) Version() string {
	return e.gaeVersion
}

func (e *environment) CloudProject() string {
	return e.googleCloudProject
}

func (e *environment) Port() string {
	return e.port
}

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
