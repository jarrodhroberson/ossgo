package firestore

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"maps"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/firestore"
	"github.com/jarrodhroberson/ossgo/functions"
	"github.com/jarrodhroberson/ossgo/functions/must"
	"github.com/jarrodhroberson/ossgo/timestamp"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/jarrodhroberson/ossgo/errors"
)

func NewCollectionRepository[T any](database DatabaseName, collection string, keyer func(t *T) string) *CollectionRepository[T] {
	return &CollectionRepository[T]{
		clientProvider: func() *firestore.Client {
			return Must(Client(context.Background(), database))
		},
		collection: collection,
		keyer:      keyer,
	}
}

func IsNotFound(err error) bool {
	return status.Code(err) == codes.NotFound
}

func Exists(err error) bool {
	return !IsNotFound(err)
}

func CollectionExists(ctx context.Context, client *firestore.Client, path string) bool {
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

func DeleteCollection(ctx context.Context, client *firestore.Client, path string) error {
	if !CollectionExists(ctx, client, path) {
		return errs.NotFoundError.New("collection \"%s\" does not exist", path)
	}

	bulkwriter := client.BulkWriter(ctx)
	defer bulkwriter.End()

	errgp, ctx := errgroup.WithContext(context.Background())
	docIter := client.Collection(path).Select().OrderBy("id", firestore.Asc).Documents(ctx)
	type idOnly struct {
		Id string `json:"id"`
	}
	for record := range DocumentIteratorToSeq[idOnly](docIter) {
		errgp.Go(func() error {
			numDeleted := 0
			for {
				doc, err := docIter.Next()
				if err != nil {
					if errors.Is(err, iterator.Done) {
						break
					} else {
						return BulkWriterError.New("error deleting collection at \"%s\"", path)
					}
				}

				_, err = bulkwriter.Delete(client.Collection(path).Doc(record.Id))
				if err != nil {
					return BulkWriterError.New("error deleting document \"%s\" in collection \"%s\"", doc.Ref.ID, path)
				}
				numDeleted++
			}

			if numDeleted == 0 {
				return nil
			}
			bulkwriter.Flush()
			return nil
		})
	}
	if err := errgp.Wait(); err != nil {
		return err
	}
	return nil
}

func Must(client *firestore.Client, err error) *firestore.Client {
	if err != nil {
		log.Error().Err(err).Msgf("error creating firestore client %s", err)
		panic(err)
	} else {
		return client
	}
}

func Client(ctx context.Context, database DatabaseName) (*firestore.Client, error) {
	if strings.Trim(string(database), " ") == "" {
		return nil, errorx.IllegalArgument.New("DatabaseName can not be an empty string")
	}
	projectId := functions.FirstNonEmpty(os.Getenv("GOOGLE_CLOUD_PROJECT"), must.Must(metadata.ProjectIDWithContext(ctx)))
	client, err := firestore.NewClientWithDatabase(ctx, projectId, string(database))
	if err != nil {
		log.Fatal().Err(err).Msgf("error creating firestore client %s", err)
		return nil, err
	}
	return client, nil
}

func Count(ctx context.Context, query firestore.Query) int64 {
	cq := query.NewAggregationQuery().WithCount("count")
	cqr, err := cq.Get(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("could not run aggregation query")
		return -1
	}
	value, ok := cqr["val"]
	if !ok {
		err = errs.MustNeverError.New("could not get \"count\" from results %s", strings.Join(slices.Collect(maps.Keys(cqr)), ","))
		log.Error().Err(err).Msg(err.Error())
		panic(err)
	}
	count, ok := value.(int64)
	if !ok {
		err := errs.MustNeverError.New("could not assert that \"%s\" was of type int64", "count")
		log.Error().Err(err).Msg(err.Error())
		panic(err)
	}
	return count
}

func GetAs[T any](ctx context.Context, database DatabaseName, path string, t *T) error {
	client := Must(Client(ctx, database))
	defer func(client *firestore.Client) {
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

func MapToUpdates(m map[string]interface{}) []firestore.Update {
	updates := make([]firestore.Update, 0, len(m))
	for k, v := range m {
		switch v.(type) {
		case string:
			updates = append(updates, firestore.Update{Path: k, Value: v.(string)})
		case int, int8, int16, int32, int64:
			updates = append(updates, firestore.Update{Path: k, Value: strconv.FormatInt(reflect.ValueOf(v).Int(), 10)})
		case uint, uint8, uint16, uint32, uint64:
			updates = append(updates, firestore.Update{Path: k, Value: strconv.FormatUint(reflect.ValueOf(v).Uint(), 10)})
		case float32, float64:
			updates = append(updates, firestore.Update{Path: k, Value: strconv.FormatFloat(reflect.ValueOf(v).Float(), 'f', -1, 64)})
		case bool:
			updates = append(updates, firestore.Update{Path: k, Value: strconv.FormatBool(v.(bool))})
		case time.Time:
			updates = append(updates, firestore.Update{Path: k, Value: timestamp.From(v.(time.Time)).String()})
		case timestamp.Timestamp:
			updates = append(updates, firestore.Update{Path: k, Value: v.(timestamp.Timestamp).String()})
		default:
			updates = append(updates, firestore.Update{Path: k, Value: string(must.MarshalJson(v))})
		}
	}
	return updates
}

func traverseFirestore(ctx context.Context, docRef firestore.DocumentRef) (map[string]interface{}, error) {
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

// DocumentIteratorToSeq2 converts a firestore.Iterator to an iter.Seq2.
// doc.Ref.ID is used as the "key" or first value, second value is a pointer to the type V
func DocumentIteratorToSeq2[V any](dsi *firestore.DocumentIterator) iter.Seq2[string, *V] {
	return func(yield func(string, *V) bool) {
		defer dsi.Stop()
		for {
			doc, err := dsi.Next()
			if errors.Is(err, iterator.Done) {
				return
			}
			if err != nil {
				log.Error().Err(err).Msg("error iterating through Firestore documents")
				return
			}

			var b V
			err = doc.DataTo(&b)
			if err != nil {
				log.Error().Err(err).Msgf("error unmarshalling Firestore document with ID %s", doc.Ref.ID)
				continue
			}

			if !yield(doc.Ref.ID, &b) {
				return
			}
		}
	}
}

// DocumentIteratorToSeq converts a firestore.Iterator to an iter.Seq.
// value is a pointer to the type V
func DocumentIteratorToSeq[V any](dsi *firestore.DocumentIterator) iter.Seq[*V] {
	return func(yield func(*V) bool) {
		defer dsi.Stop()
		for {
			doc, err := dsi.Next()
			if errors.Is(err, iterator.Done) {
				return
			}
			if err != nil {
				log.Error().Err(err).Msg("error iterating through Firestore documents")
				return
			}

			var b V
			err = doc.DataTo(&b)
			if err != nil {
				log.Error().Err(err).Msgf("error unmarshalling Firestore document with ID %s", doc.Ref.ID)
				return
			}

			if !yield(&b) {
				return
			}
		}
	}
}
