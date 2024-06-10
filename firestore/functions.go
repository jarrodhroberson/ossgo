package firestore

import (
	"context"
	"errors"
	"os"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/firestore/apiv1/firestorepb"
	errs "github.com/jarrodhroberson/ossgo/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func IsDocumentNotFound(err error) bool {
	return status.Code(err) == codes.NotFound
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
	client, err := firestore.NewClientWithDatabase(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"), string(database))
	if err != nil {
		log.Fatal().Err(err).Msgf("error creating firestore client %s", err)
		return nil, err
	}
	return client, nil
}

func Count(ctx context.Context, query firestore.Query) int64 {
	cq := query.NewAggregationQuery().WithCount("val")
	cqr, err := cq.Get(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("could not run aggregation query")
		return -1
	}
	value, ok := cqr["val"]
	if !ok {
		log.Fatal().Err(err).Msg("could not get \"all\" alias for count")
		return -1
	}
	val := value.(*firestorepb.Value)
	return val.GetIntegerValue()
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
