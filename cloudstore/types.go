package cloudstore

import (
	"fmt"
)

type Location struct {
	Bucket string
	Name   string
}

func (l Location) String() string {
	return fmt.Sprintf("/%s/%s", l.Bucket, l.Name)
}
