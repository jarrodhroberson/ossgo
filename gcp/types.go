package gcp

import (
	"sync"
)

var environment *Environment
var once sync.Once

type Environment struct {
	gin_mode             string
	gae_application      string
	gae_deployment_id    string
	gae_env              string
	gae_instance         string
	gae_memory_mb        string
	gae_runtime          string
	gae_service          string
	gae_version          string
	google_cloud_project string
	port                 string
}

func (e *Environment) Application() string {
	return e.gae_application
}

func (env *Environment) GinMode() string {
	return env.gin_mode
}

func (e *Environment) DeploymentId() string {
	return e.gae_deployment_id
}

func (e *Environment) Env() string {
	return e.gae_env
}

func (e *Environment) Instance() string {
	return e.gae_instance
}

func (e *Environment) MemoryMb() string {
	return e.gae_memory_mb
}

func (e *Environment) Runtime() string {
	return e.gae_runtime
}

func (e *Environment) Service() string {
	return e.gae_service
}

func (e *Environment) Version() string {
	return e.gae_version
}

func (e *Environment) CloudProject() string {
	return e.google_cloud_project
}

func (e *Environment) Port() string {
	return e.port
}
