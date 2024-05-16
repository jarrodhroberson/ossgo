package firebase

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
)

func Must(client *firestore.Client, err error) *firestore.Client {
	if err != nil {
		log.Error().Err(err).Msgf("error creating firestore client %s", err)
		panic(err)
	} else {
		return client
	}
}

func Client(ctx context.Context, database string) (*firestore.Client, error) {
	if database == "" {
		client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
		if err != nil {
			log.Error().Err(err).Msgf("error creating firestore client %s", err)
			return nil, err
		}
		return client, nil
	} else {
		client, err := firestore.NewClientWithDatabase(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"), string(database))
		if err != nil {
			log.Fatal().Err(err).Msgf("error creating firestore client %s", err)
			return nil, err
		}
		return client, nil
	}
}
