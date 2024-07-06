package firestore

import (
	"context"
	"errors"
	"os"
	"strings"

	fs "cloud.google.com/go/firestore"
	fspb "cloud.google.com/go/firestore/apiv1/firestorepb"
	errs "github.com/jarrodhroberson/ossgo/errors"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func IsDocumentNotFound(err error) bool {
	return status.Code(err) == codes.NotFound
}

func DeleteCollection(ctx context.Context, client *fs.Client, path string) error {
	colRef := client.Collection(path)
	bulkwriter := client.BulkWriter(ctx)
	defer bulkwriter.End()

	for {
		iter := colRef.Limit(500).Documents(ctx)
		numDeleted := 0
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			_, err = bulkwriter.Delete(doc.Ref)
			if err != nil {
				return err
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
