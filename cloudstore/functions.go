package cloudstore

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"

	gcs "cloud.google.com/go/storage"
	"github.com/jarrodhroberson/ossgo/timestamp"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
)

var validBucketNamePattern = "^[a-z0-9]{1}[a-z0-9-_]{1,62}[a-z0-9]{1}$"
var validBucketNameRegEx = regexp.MustCompile(validBucketNamePattern)

func isValidBucketName(bucket string) bool {
	return validBucketNameRegEx.MatchString(bucket)
}

func newClient(ctx context.Context) *gcs.Client {
	client, err := gcs.NewClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}
	return client
}

func AttrToObject(attrs *gcs.ObjectAttrs) Metadata {
	return Metadata{
		Path:          attrs.Name,
		Size:          attrs.Size,
		ContentType:   attrs.ContentType,
		CreatedAt:     timestamp.From(attrs.Created),
		LastUpdatedAt: timestamp.From(attrs.Updated),
	}
}

func BucketWithObjects(bucket string) (*Bucket, error) {
	ctx := context.Background()
	client := newClient(ctx)
	defer client.Close()
	bh := client.Bucket(bucket)
	if !BucketExists(ctx, bh) {
		return nil, fmt.Errorf("bucket %s does not exist", bucket)
	}
	it := bh.Objects(ctx, nil)
	objects := make([]Metadata, 0, it.PageInfo().Remaining())
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		objects = append(objects, AttrToObject(attrs))
	}
	return &Bucket{
		Name:    bucket,
		Objects: objects,
	}, nil
}

func BucketExists(ctx context.Context, bh *gcs.BucketHandle) bool {
	_, err := bh.Attrs(ctx)
	return err != gcs.ErrBucketNotExist
}

func ListAllBuckets(ctx context.Context) []string {
	client := newClient(ctx)
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	it := client.Buckets(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	buckets := make([]string, 0, it.PageInfo().Remaining())
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return make([]string, 0, 0)
		}
		buckets = append(buckets, battrs.Name)
	}
	return buckets
}

func CreateBucket(ctx context.Context, bucket string) (*gcs.BucketHandle, error) {
	if !isValidBucketName(bucket) {
		return nil, fmt.Errorf("%s is not a valid bucket name; did not match %s", bucket, validBucketNamePattern)
	}
	client := newClient(ctx)
	defer client.Close()

	bh := client.Bucket(bucket)
	if BucketExists(ctx, bh) {
		return nil, fmt.Errorf("bucket %s already exists", bucket)
	}
	err := bh.Create(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"), nil)
	if err != nil {
		return nil, err
	}
	return bh, nil
}

func DeleteBucket(ctx context.Context, bucket string) error {
	client := newClient(ctx)
	defer client.Close()
	return client.Bucket(bucket).Delete(ctx)
}

func ReadObject(ctx context.Context, bucket string, path string, dst io.Writer) ReadResult {
	client := newClient(ctx)
	defer client.Close()
	bh := client.Bucket(bucket)
	src, err := bh.Object(path).NewReader(ctx)
	defer src.Close()
	if err != nil {
		return ReadResult{
			BytesRead: 0,
			Error:     err,
		}
	}
	i, err := io.Copy(dst, src)
	return ReadResult{
		BytesRead: i,
		Error:     err,
	}
}

func WriteObject(ctx context.Context, bucket string, path string, src io.Reader) WriteResult {
	client := newClient(ctx)
	defer client.Close()
	bh := client.Bucket(bucket)
	dst := bh.Object(path).NewWriter(ctx)
	defer dst.Close()
	i, err := io.Copy(dst, src)
	return WriteResult{
		BytesWritten: i,
		Error:        err,
	}
}

func DeleteObject(ctx context.Context, bh *gcs.BucketHandle, path string) error {
	client := newClient(ctx)
	defer client.Close()

	return bh.Object(path).Delete(ctx)
}
