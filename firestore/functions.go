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

	"github.com/jarrodhroberson/ossgo/containers"
	"github.com/jarrodhroberson/ossgo/gcp"
	slyces "github.com/jarrodhroberson/ossgo/slices"

	"cloud.google.com/go/compute/metadata"
	fs "cloud.google.com/go/firestore"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/joomcode/errorx"
	"github.com/rs/zerolog/log"

	"github.com/jarrodhroberson/ossgo/functions/must"
	"github.com/jarrodhroberson/ossgo/seq"
	strs "github.com/jarrodhroberson/ossgo/strings"
	"github.com/jarrodhroberson/ossgo/timestamp"

	errs "github.com/jarrodhroberson/ossgo/errors"
)

// NewProjection creates a new projection instance for querying specific fields from Firestore documents.
// The paths parameter accepts dot-separated paths that specify the fields to select.
// For example, "user.address.street" will be split into a FieldPath ["user", "address", "street"].
// Each field path is converted to a fs.FieldPath internally for use in Firestore queries.
func NewProjection(paths ...string) Projection {
	pathIter := slices.Values(paths)
	fieldPathIter := slyces.Transform[string, fs.FieldPath](pathIter, func(path string) fs.FieldPath {
		return fs.FieldPath(strings.Split(path, "."))
	})
	return Projection{
		fieldPaths: slices.Collect(fieldPathIter),
	}
}

// NewQuery creates a new Query instance for the specified Firestore collection.
func NewQuery(collection *fs.CollectionRef) Query {
	return &newQuery{collection: collection}
}

// DocRefIDKeyer returns a Keyer function that extracts the ID of a DocumentRef as a string.
// This is useful for creating maps or sets keyed by DocumentRef IDs.
func DocRefIDKeyer() containers.Keyer[fs.DocumentRef] {
	return func(docRef *fs.DocumentRef) string {
		return docRef.ID
	}
}

// DocSnapShotKeyer returns a Keyer function that extracts the ID of a DocumentSnapshot as a string.
func DocSnapShotKeyer() containers.Keyer[fs.DocumentSnapshot] {
	return func(dss *fs.DocumentSnapshot) string {
		return dss.Ref.ID
	}
}

// NewCollectionStore creates a new collectionStore for a given database, collection, and keyer function.
func NewCollectionStore[T any](database DatabaseName, collection string, keyerFunc containers.Keyer[T]) *collectionStore[T] {
	return &collectionStore[T]{
		clientProvider: func() *fs.Client {
			return must.Must(Client(context.Background(), database))
		},
		collection: collection,
		keyer:      keyerFunc,
	}
}

// IsNotFound checks if the given error is a Firestore "not found" error.
func IsNotFound(err error) bool {
	return err != nil && status.Code(err) == codes.NotFound
}

// IsAlreadyExists checks if the given error is a Firestore "already exists" error.
// It returns true if the error is not nil and its status code indicates the resource already exists.
func IsAlreadyExists(err error) bool {
	return err != nil && status.Code(err) == codes.AlreadyExists
}

// Exists checks if the given error is not a Firestore "not found" error.
func Exists(err error) bool {
	return !IsNotFound(err)
}

// DocRefExists checks if the given DocumentRef exists.
func DocRefExists(ctx context.Context, docRef *fs.DocumentRef) bool {
	docSS, err := docRef.Get(ctx)
	return !IsNotFound(err) && docSS != nil && docSS.Exists()
}

// DeleteCollection deletes all documents in a Firestore collection using a BulkWriter.
// It processes documents in parallel batches of size MAX_BULK_WRITE_SIZE.
//
// ctx is the context for the operation
// client is the Firestore client to use
// path is the path to the collection to delete
//
// Returns an error if any document deletion fails, nil on success
func DeleteCollection(ctx context.Context, client *fs.Client, path string) error {
	bw := client.BulkWriter(ctx)
	defer func() {
		bw.Flush()
		bw.End()
	}()

	errgp, ctx := errgroup.WithContext(context.Background())
	docIter := client.Collection(path).Select().OrderBy("created_at", fs.Asc).Documents(ctx)
	batches := seq.Chunk(DocumentIteratorToSeq(docIter), MAX_BULK_WRITE_SIZE)
	for batch := range batches {
		errgp.Go(func() error {
			for docSS := range batch {
				_, err := bw.Delete(docSS.Ref)
				if err != nil {
					return BulkWriterError.Wrap(err, "error deleting document \"%s\" in collection \"%s\"", docSS.Ref.ID, path)
				}
			}
			bw.Flush()
			return nil
		})
	}

	return errgp.Wait()
}

// Client creates a new Firestore client for the specified database.
func Client(ctx context.Context, database DatabaseName) (*fs.Client, error) {
	if strings.Trim(string(database), " ") == "" {
		return nil, errorx.IllegalArgument.New("DatabaseName can not be an empty string")
	}
	projectId := strs.FirstNonEmpty(os.Getenv("GOOGLE_CLOUD_PROJECT"), must.Must(gcp.ProjectId()), must.Must(metadata.ProjectIDWithContext(ctx)))
	if projectId == "" {
		return nil, errorx.IllegalArgument.New("projectId can not be an empty string")
	}
	client, err := fs.NewClientWithDatabase(ctx, projectId, string(database))
	if err != nil {
		log.Fatal().Err(err).Msgf("error creating firestore client %s", err)
		return nil, err
	}
	return client, nil
}

// Count returns the number of documents that match the given query.
func Count(ctx context.Context, query fs.Query) int64 {
	cq := query.NewAggregationQuery().WithCount("count")
	cqr, err := cq.Get(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("could not run aggregation query")
		return -1
	}
	value, ok := cqr["count"]
	if !ok {
		err = errs.MustNeverError.New("could not get \"count\" from results %s", strings.Join(slices.Collect(maps.Keys(cqr)), ","))
		log.Error().Stack().Err(err).Msg(err.Error())
		panic(errorx.Panic(err))
	}
	count, ok := value.(int64)
	if !ok {
		err := errs.MustNeverError.New("could not assert that \"%s\" was of type int64", "count")
		log.Error().Stack().Err(err).Msg(err.Error())
		panic(errorx.Panic(err))
	}
	return count
}

// MapToUpdates converts a map to a slice of Firestore Update structs.
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
			updates = append(updates, fs.Update{Path: k, Value: v.(timestamp.Timestamp)})
		default:
			updates = append(updates, fs.Update{Path: k, Value: string(must.MarshalJson(v))})
		}
	}
	return updates
}

// traverseFirestore recursively traverses a Firestore document and its subcollections.
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

// DocSnapShotToType unmarshals a Firestore DocumentSnapshot into a struct of type T.
func DocSnapShotToType[T any](dss *fs.DocumentSnapshot) (*T, error) {
	var d T
	m := make(map[string]interface{})
	err := dss.DataTo(&m)
	if err != nil {
		err = errs.MarshalError.Wrap(err, "error unmarshalling Firestore document with ID %s", dss.Ref.ID)
		return nil, err
	}
	must.UnmarshallMap(m, &d)
	return &d, nil
}

// DocSnapShotSeq2ToType converts a Seq2 of DocumentSnapshots to a Seq2 of type V.
func DocSnapShotSeq2ToType[V any](it iter.Seq2[string, *fs.DocumentSnapshot]) iter.Seq2[string, *V] {
	return seq.Map2[string, *fs.DocumentSnapshot, string, *V](it, seq.PassThruFunc[string], func(v *fs.DocumentSnapshot) *V {
		return must.Must(DocSnapShotToType[V](v))
	})
}

// DocSnapShotSeqToType converts a Seq of DocumentSnapshots to a Seq of type R.
// Deprecated: Use DocumentIterToTypeSeq which automatically wraps the fs.DocumentIterator
// and more efficient data processing. This delegates to an unexported implementation, so this can be deleted in a future release.
func DocSnapShotSeqToType[R any](it iter.Seq[*fs.DocumentSnapshot]) iter.Seq[*R] {
	return docSnapShotSeqToType[R](it)
}

// docSnapShotSeqToType converts a Seq of DocumentSnapshots to a Seq of type R.
// This is an unexported version of DocSnapShotSeqToType that will replace the exported version
// once it is no longer being used. The exported version is deprecated and will be removed in a future release.
// This function provides the same functionality but with better integration with the DocumentIterator wrapper functions.
func docSnapShotSeqToType[R any](it iter.Seq[*fs.DocumentSnapshot]) iter.Seq[*R] {
	return seq.Map[*fs.DocumentSnapshot, *R](it, func(dss *fs.DocumentSnapshot) *R {
		var t R
		m := make(map[string]interface{})
		err := dss.DataTo(&m)
		if err != nil {
			log.Error().Stack().Err(err).Msgf("error unmarshalling Firestore document with ID %s", dss.Ref.ID)

			panic(errorx.Panic(err))
		}
		must.UnmarshallMap(m, &t)
		return &t
	})
}

// DocumentIteratorToSeq converts a firestore.Iterator to an iter.Seq.
// value is a pointer to the type V
func DocumentIteratorToSeq(dsi *fs.DocumentIterator) iter.Seq[*fs.DocumentSnapshot] {
	return func(yield func(ref *fs.DocumentSnapshot) bool) {
		defer dsi.Stop()
		for {
			docSS, err := dsi.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				err = errs.MustNeverError.Wrap(err, "error iterating through Firestore documents")
				log.Error().Stack().Err(err).Msg(err.Error())
				break
			}
			if !yield(docSS) {
				return
			}
		}
	}
}

// DocumentIteratorToSeq2 converts a firestore.Iterator to an iter.Seq2.
// doc.Ref.ID is used as the "key" or first value, second value is a pointer to the type V
func DocumentIteratorToSeq2(dsi *fs.DocumentIterator) iter.Seq2[string, *fs.DocumentSnapshot] {
	return seq.ToSeq2[string, *fs.DocumentSnapshot](DocumentIteratorToSeq(dsi), func(v *fs.DocumentSnapshot) string {
		return v.Ref.ID
	})
}

// CountDuplicateDocumentIds finds and reports duplicate document IDs within a Firestore collection.
// It takes the Firestore client and the collection path as input.
// It returns an iter.Seq2[string,int] where keys are duplicate document Ids and values are the number of occurrences.
// If no duplicates are found, it returns an empty iter.Seq2[string,int.  Returns an error if one occurs.
func CountDuplicateDocumentIds(ctx context.Context, databaseName DatabaseName, collectionPath string) (iter.Seq2[string, int], error) {
	client := must.Must(Client(ctx, databaseName))
	defer func(client *fs.Client) {
		err := client.Close()
		if err != nil {
			log.Warn().Err(err).Msgf("error closing firestore client %s", err)
		}
	}(client)

	// Iterate over all documents in the collection.
	docRefIter := client.Collection(collectionPath).Documents(ctx)
	docRefSeq := DocumentIteratorToSeq(docRefIter)
	duplicates := make(map[string]int)
	for docRef := range docRefSeq {
		id := docRef.Ref.ID
		// Check if the document ID already exists.
		if _, ok := duplicates[id]; ok {
			duplicates[id] = duplicates[id] + 1
		} else {
			duplicates[id] = 1
		}
	}
	return containers.Seq2(duplicates), nil
}

// ClosingWhenDoneSeq wraps the provided iter.Seq and ensures that fs.Client.Close() is called
// after the last item is provided by the Seq.
func ClosingWhenDoneSeq[T any](seq iter.Seq[T], client *fs.Client) iter.Seq[T] {
	return func(yield func(item T) bool) {
		defer func() {
			if err := client.Close(); err != nil {
				log.Warn().Err(err).Msg("error closing Firestore client")
			}
		}()
		for item := range seq {
			if !yield(item) {
				return
			}
		}
	}
}

// ClosingWhenDoneSeq2 wraps the provided iter.Seq2 and ensures that fs.Client.Close() is called
// after the last item is provided by the Seq2.
func ClosingWhenDoneSeq2[K, V any](seq2 iter.Seq2[K, V], client *fs.Client) iter.Seq2[K, V] {
	return func(yield func(key K, value V) bool) {
		defer func() {
			if err := client.Close(); err != nil {
				log.Warn().Err(err).Msg("error closing Firestore client")
			}
		}()
		for k, v := range seq2 {
			if !yield(k, v) {
				return
			}
		}
	}
}

// CollectionIterToSeq converts a Firestore CollectionIterator to an iter.Seq of CollectionRefs.
// It handles the iteration and error handling internally, providing a simplified interface to process collections.
// Returns an iter.Seq that yields *fs.CollectionRef values.
func CollectionIterToSeq(ci *fs.CollectionIterator) iter.Seq[*fs.CollectionRef] {
	return func(yield func(ref *fs.CollectionRef) bool) {
		for {
			colRef, err := ci.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				err = errs.MustNeverError.Wrap(err, "error iterating through Firestore documents")
				log.Error().Stack().Err(err).Msg(err.Error())
				break
			}
			if !yield(colRef) {
				return
			}
		}
	}
}

// DocumentIterToTypeSeq converts a Firestore DocumentIterator to an iter.Seq of type T.
// It first converts the DocumentIterator to DocumentSnapshots using DocumentIteratorToSeq,
// then converts those snapshots to type T using DocSnapShotSeqToType.
// Type parameter T should be the target type to unmarshal the documents into.
// Returns an iter.Seq that yields pointers to values of type T.
func DocumentIterToTypeSeq[T any](di *fs.DocumentIterator) iter.Seq[*T] {
	return docSnapShotSeqToType[T](DocumentIteratorToSeq(di))
}
