package secrets

import (
	"fmt"
)

const OAUTH_CLIENT_ID = "OAUTH_CLIENT_ID"
const OAUTH_CLIENT_SECRET = "OAUTH_CLIENT_SECRET"

type Path struct {
	ProjectNumber int
	Name          string
	Version       int
}

func (p Path) WithoutVersion() string {
	return fmt.Sprintf(pathToSecret, p.ProjectNumber, p.Name)
}

func (p Path) LatestVersion() string {
	return fmt.Sprintf(pathToLatestVersion, p.ProjectNumber, p.Name)
}

func (p Path) WithVersion() string {
	return fmt.Sprintf(pathToNumericVersion, p.ProjectNumber, p.Name, p.Version)
}

func (p Path) String() string {
	if p.Version == 0 {
		return p.LatestVersion()
	} else if p.Version < 0 {
		return p.WithoutVersion()
	} else {
		return p.WithVersion()
	}
}
