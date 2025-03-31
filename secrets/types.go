package secrets

import (
	"fmt"
)

// Path represents a path to a secret in Secret Manager.
//
// A path can be represented in three different ways:
// - Without a version: projects/{project_number}/secrets/{secret_name}
// - With a version: projects/{project_number}/secrets/{secret_name}/versions/{version_number}
// - With the latest version: projects/{project_number}/secrets/{secret_name}/versions/latest
type Path struct {
	ProjectNumber int
	Name          string
	Version       int
}

// WithoutVersion returns the path to the secret without a version.
// Example: projects/1234567890/secrets/my-secret
func (p Path) WithoutVersion() string {
	return fmt.Sprintf(pathToSecret, p.ProjectNumber, p.Name)
}

// LatestVersion returns the path to the latest version of the secret.
// Example: projects/1234567890/secrets/my-secret/versions/latest
func (p Path) LatestVersion() string {
	return fmt.Sprintf(pathToLatestVersion, p.ProjectNumber, p.Name)
}

func (p Path) WithVersion() string {
	return fmt.Sprintf(pathToNumericVersion, p.ProjectNumber, p.Name, p.Version)
}

// String returns the string representation of the path.
// If the version is 0, it returns the latest version.
// If the version is negative, it returns the path without a version.
// Otherwise, it returns the path with the version.
func (p Path) String() string {
	if p.Version == 0 {
		return p.LatestVersion()
	} else if p.Version < 0 {
		return p.WithoutVersion()
	} else {
		return p.WithVersion()
	}
}
