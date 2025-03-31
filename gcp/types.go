package gcp

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
