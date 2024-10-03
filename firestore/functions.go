package firestore

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	fs "cloud.google.com/go/firestore"
	fspb "cloud.google.com/go/firestore/apiv1/firestorepb"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/jarrodhroberson/ossgo/errors"
)

func IsNotFound(err error) bool {
	return status.Code(err) == codes.NotFound
}

func Exists(err error) bool {
	return !IsNotFound(err)
}

func CollectionExists(ctx context.Context, client *fs.Client, path string) bool {
	iter := client.Collections(ctx)
	for {
		colRef, err := iter.Next()
		if err == iterator.Done {
			return false
		}
		if err != nil {
			panic(err)
		}
		if colRef.Path == path {
			return true
		}
	}
}

func DeleteCollection(ctx context.Context, client *fs.Client, path string) error {
	colRef := client.Collection(path)
	if !CollectionExists(ctx, client, path) {
		return errs.NotFoundError.New("collection \"%s\" does not exist", path)
	}

	bulkwriter := client.BulkWriter(ctx)
	defer bulkwriter.End()

	for {
		iter := colRef.Limit(500).Documents(ctx)
		numDeleted := 0
		for {
			doc, err := iter.Next()
			if err != nil {
				if err == iterator.Done {
					break
				} else {
					return BulkWriterError.New("error deleting collection at \"%s\"", path)
				}
			}

			_, err = bulkwriter.Delete(doc.Ref)
			if err != nil {
				return BulkWriterError.New("error deleting document \"%s\" in collection \"%s\"", doc.Ref.ID, path)
			}
			numDeleted++
		}

		if numDeleted == 0 {
			break
		}

		bulkwriter.Flush()
	}
	return nil
}

func Must(client *fs.Client, err error) *fs.Client {
	if err != nil {
		log.Error().Err(err).Msgf("error creating firestore client %s", err)
		panic(err)
	} else {
		return client
	}
}

func Client(ctx context.Context, database DatabaseName) (*fs.Client, error) {
	if strings.Trim(string(database), " ") == "" {
		return nil, errorx.IllegalArgument.New("DatabaseName can not be an empty string")
	}
	client, err := fs.NewClientWithDatabase(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"), string(database))
	if err != nil {
		log.Fatal().Err(err).Msgf("error creating firestore client %s", err)
		return nil, err
	}
	return client, nil
}

func Count(ctx context.Context, query fs.Query) int64 {
	cq := query.NewAggregationQuery().WithCount("val")
	cqr, err := cq.Get(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("could not run aggregation query")
		return -1
	}
	value, ok := cqr["val"]
	if !ok {
		log.Fatal().Err(err).Msg("could not get \"val\" alias for count")
		return -1
	}
	return value.(*fspb.Value).GetIntegerValue()
}

func GetAs[T any](ctx context.Context, database DatabaseName, path string, t *T) error {
	client := Must(Client(ctx, database))
	defer func(client *fs.Client) {
		err := client.Close()
		if err != nil {
			log.Error().Err(err).Msgf("error closing firestore client %s", err)
		}
	}(client)
	doc, err := client.Doc(path).Get(ctx)
	if err != nil {
		err = errors.Join(errs.NotFoundError.New("could not find document %s", path), err)
		return err
	}
	return doc.DataTo(t)
}

func MapToUpdates(m map[string]interface{}) []fs.Update {
	updates := make([]fs.Update, 0, len(m))
	for k, v := range m {
		updates = append(updates, fs.Update{Path: k, Value: v})
	}
	return updates
}

func traverseFirestore(ctx context.Context, docRef fs.DocumentRef) (map[string]interface{}, error) {
	var tree map[string]interface{}

	// Get the document snapshot
	colIter := docRef.Collections(ctx)
	for {
		colRef, err := colIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		tree[colRef.Path] = make([]interface{}, 0)
	}

	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Extract document data and collections
	data := docSnap.Data()
	for k, v := range data {
		switch v.(type) {
		case map[string]interface{}:
			tree[k] = v.(map[string]interface{})
		case []interface{}:
			tree[k] = v.([]interface{})
		case string, float64, bool:
			tree[k] = v
		case int64:
			tree[k] = strconv.Itoa(int(v.(int64)))
		case uint64:
			tree[k] = strconv.Itoa(int(v.(uint64)))
		case int:
			tree[k] = strconv.Itoa(v.(int))
		case uint:
			tree[k] = strconv.Itoa(int(v.(uint)))
		default:
			return nil, fmt.Errorf("unsupported type %T for key %s", v, k)
		}
	}

	return tree, nil
}
