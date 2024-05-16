package secrets

import (
	"fmt"
	"time"

	"github.com/jarrodhroberson/ossgo/functions/must"
	"github.com/kofalt/go-memoize"
)

var pathMemoizer *memoize.Memoizer

func init() {
	//pathMemoizer = memoize.NewMemoizer(90*time.Second, 10*time.Minute)
	pathMemoizer = memoize.NewMemoizer(3*time.Second, 10*time.Second)
}

func DumpCache() {
	fmt.Println(string(must.MarshalJson(pathMemoizer.Storage.Items())))
}

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

func (p *Path) String() string {
	return must.Call(pathMemoizer, fmt.Sprintf("%p", p), func() (string, error) {
		if p.Version == 0 {
			return p.LatestVersion(), nil
		} else if p.Version < 0 {
			return p.WithoutVersion(), nil
		} else {
			return p.WithVersion(), nil
		}
	})
}
