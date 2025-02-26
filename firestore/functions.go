package firestore

import (
	"context"
	"errors"
	"fmt"
	"github.com/jarrodhroberson/ossgo/containers"
	"iter"
	"maps"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/compute/metadata"
	fs "cloud.google.com/go/firestore"
	"github.com/jarrodhroberson/ossgo/functions"
	"github.com/jarrodhroberson/ossgo/functions/must"
	"github.com/jarrodhroberson/ossgo/seq"
	"github.com/jarrodhroberson/ossgo/timestamp"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/jarrodhroberson/ossgo/errors"
)

func NewQuery(collection *fs.CollectionRef) Query {
	return &newQuery{collection: collection}
}

func DocRefIDKeyer() containers.Keyer[fs.DocumentRef] {
	return func(docRef *fs.DocumentRef) string {
		return docRef.ID
	}
}

func DocSnapShotKeyer() containers.Keyer[fs.DocumentSnapshot] {
	return func(dss *fs.DocumentSnapshot) string {
		return dss.Ref.ID
	}
}

func NewCollectionStore[T any](database DatabaseName, collection string, keyerFunc containers.Keyer[T]) *CollectionStore[T] {
	return &CollectionStore[T]{
		clientProvider: func() *fs.Client {
			return Must(Client(context.Background(), database))
		},
		collection: collection,
		keyer: keyerFunc,
	}
}

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
	if !CollectionExists(ctx, client, path) {
		return errs.NotFoundError.New("collection \"%s\" does not exist", path)
	}

	bulkwriter := client.BulkWriter(ctx)
	defer bulkwriter.End()

	errgp, ctx := errgroup.WithContext(context.Background())
	docIter := client.Collection(path).Select().OrderBy("created_at", fs.Asc).Documents(ctx)

	for docSSIter := range seq.Chunk(DocumentIteratorToSeq(docIter), 500) {
		errgp.Go(func() error {
			for docSS := range docSSIter {
				_, err := bulkwriter.Delete(docSS.Ref)
				if err != nil {
					return BulkWriterError.New("error deleting document \"%s\" in collection \"%s\"", docSS.Ref.ID, path)
				}
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
	projectId := functions.FirstNonEmpty(os.Getenv("GOOGLE_CLOUD_PROJECT"), must.Must(metadata.ProjectIDWithContext(ctx)))
	client, err := fs.NewClientWithDatabase(ctx, projectId, string(database))
	if err != nil {
		log.Fatal().Err(err).Msgf("error creating firestore client %s", err)
		return nil, err
	}
	return client, nil
}

func Count(ctx context.Context, query fs.Query) int64 {
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
		switch v.(type) {
		case string:
			updates = append(updates, fs.Update{Path: k, Value: v.(string)})
		case int, int8, int16, int32, int64:
			updates = append(updates, fs.Update{Path: k, Value: strconv.FormatInt(reflect.ValueOf(v).Int(), 10)})
		case uint, uint8, uint16, uint32, uint64:
			updates = append(updates, fs.Update{Path: k, Value: strconv.FormatUint(reflect.ValueOf(v).Uint(), 10)})
		case float32, float64:
			updates = append(updates, fs.Update{Path: k, Value: strconv.FormatFloat(reflect.ValueOf(v).Float(), 'f', -1, 64)})
		case bool:
			updates = append(updates, fs.Update{Path: k, Value: strconv.FormatBool(v.(bool))})
		case time.Time:
			updates = append(updates, fs.Update{Path: k, Value: timestamp.From(v.(time.Time)).String()})
		case timestamp.Timestamp:
			updates = append(updates, fs.Update{Path: k, Value: v.(timestamp.Timestamp).String()})
		default:
			updates = append(updates, fs.Update{Path: k, Value: string(must.MarshalJson(v))})
		}
	}
	return updates
}

func traverseFirestore(ctx context.Context, docRef fs.DocumentRef) (map[string]interface{}, error) {
	var tree map[string]interface{}

	// Load the document snapshot
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

func DocSnapShotToType[T any](dss *fs.DocumentSnapshot) (*T, error) {
	var d T
	err := dss.DataTo(&d)
	if err != nil {
		err = errs.MarshalError.Wrap(err, "error unmarshalling Firestore document with ID %s", dss.Ref.ID)
		return nil, err
	}
	return &d, nil
}

func DocSnapShotSeq2ToType[V any](it iter.Seq2[string, *fs.DocumentSnapshot]) iter.Seq2[string, *V] {
	return seq.Map2[string, *fs.DocumentSnapshot, string, *V](it, seq.PassThruFunc[string], func(v *fs.DocumentSnapshot) *V {
		return must.Must(DocSnapShotToType[V](v))
	})
}

func DocSnapShotSeqToType[R any](it iter.Seq[*fs.DocumentSnapshot]) iter.Seq[*R] {
	return seq.Map[*fs.DocumentSnapshot, *R](it, func(dss *fs.DocumentSnapshot) *R {
		var t R
		err := dss.DataTo(&t)
		if err != nil {
			log.Error().Err(err).Msgf("error unmarshalling Firestore document with ID %s", dss.Ref.ID)
			panic(err)
		}
		return &t
	})
}

// DocumentIteratorToSeq converts a firestore.Iterator to an iter.Seq.
// value is a pointer to the type V
func DocumentIteratorToSeq(dsi *fs.DocumentIterator) iter.Seq[*fs.DocumentSnapshot] {
	return func(yield func(ref *fs.DocumentSnapshot) bool) {
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
			if !yield(doc) {
				return
			}
		}
	}
}

// DocumentIteratorToSeq2 converts a firestore.Iterator to an iter.Seq2.
// doc.Ref.ID is used as the "key" or first value, second value is a pointer to the type V
func DocumentIteratorToSeq2(dsi *fs.DocumentIterator) iter.Seq2[string, *fs.DocumentSnapshot] {
	return seq.SeqToSeq2[string, *fs.DocumentSnapshot](DocumentIteratorToSeq(dsi), func(v *fs.DocumentSnapshot) string {
		return v.Ref.ID
	})
}
