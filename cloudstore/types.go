package cloudstore

import (
	"fmt"
	"path"

	"github.com/jarrodhroberson/ossgo/timestamp"
)

type Location struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

func (l Location) String() string {
	return fmt.Sprintf("/%s/%s", l.Bucket, l.Name)
}

type ReadResult struct {
	BytesRead int64 `json:"bytes_read"`
	Error     error `json:"error"`
}

type WriteResult struct {
	BytesWritten int64 `json:"bytes_written"`
	Error        error `json:"error"`
}

type Metadata struct {
	Path          string               `json:"path"`
	Size          int64                `json:"size"`
	ContentType   string               `json:"content_type"`
	CreatedAt     *timestamp.Timestamp `json:"created_at"`
	LastUpdatedAt *timestamp.Timestamp `json:"last_updated_at"`
}

func (o Metadata) Name() string {
	return path.Base(o.Path)
}

func (o Metadata) Dir() string {
	return path.Dir(o.Path)
}

type Bucket struct {
	Name    string     `json:"name"`
	Objects []Metadata `json:"objects"`
}
