package sendgrid

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	fs "github.com/jarrodhroberson/ossgo/firestore"
	"github.com/rs/zerolog/log"
)

// Region is the region where your Cloud Functions are deployed.
const Region string = "your-region"

// LeaseTime is the duration for which a function execution holds the lease to
// prevent duplicate emails.
const LeaseTime = 60 * time.Second

// flakySuccessRatio simulates the probability of success for the flaky service.
var flakySuccessRatio = 0.5

func initFirebaseApp(ctx context.Context, database string) *firestore.Client {
	if database == "" {
		client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
		if err != nil {
			log.Fatal().Err(err).Msgf("error creating firestore client %s", err)
			return nil
		}
		return client
	} else {
		client, err := firestore.NewClientWithDatabase(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"), string(database))
		if err != nil {
			log.Fatal().Err(err).Msgf("error creating firestore client %s", err)
			return nil
		}
		return client
	}
}

// shouldSend checks if the email should be sent based on the existence of a
// document in Firestore and a lease mechanism to prevent duplicates.
func shouldSend(ctx context.Context, eventID string) error {
	client := fs.Must(fs.Client(ctx, "sendgrid"))
	defer func(db *firestore.Client) {
		err := db.Close()
		if err != nil {
			log.Error().Err(err).Msg("could not close firebase client")
		}
	}(client)

	emailRef := client.Collection("sent").Doc(eventID)
	return client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		docSnapshot, err := tx.Get(emailRef)
		if err != nil {
			return err
		}
		if docSnapshot.Exists() {
			data := docSnapshot.Data()
			if data["sent"] == true {
				return nil
			}
			now := time.Now()
			lease, ok := data["lease"].(time.Time)
			if ok && now.Before(lease) {
				return fmt.Errorf("lease already taken, try later.")
			}
		}
		err = tx.Set(emailRef, map[string]interface{}{"lease": time.Now().Add(LeaseTime)})
		if err != nil {
			return err
		}
		return nil
	})
}

// markSent marks the email as sent in Firestore.
func markSent(ctx context.Context, eventID string) error {
	client := fs.Must(fs.Client(ctx, "sendgrid"))
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Error().Err(err).Msg("could not close firebase client")
		}
	}(client)
	//emailRef := client.Collection("sent").Doc(eventID)
	_, err := client.Collection("sent").Doc(eventID).Set(ctx, map[string]interface{}{"sent": true})
	return err
}

// handleNonIdempotentFirestoreFunction handles the non-idempotent Pub/Sub event
// for adding a document to Firestore and calling a flaky service.
//func handleNonIdempotentFirestoreFunction(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//	body, err := io.ReadAll(r.Body)
//	if err != nil {
//		fmt.Fprintf(w, "Error reading request body: %v", err)
//		return
//	}
//	content := map[string]interface{}{}
//	err = json.Unmarshal(body, &content)
//	if err != nil {
//		fmt.Fprintf(w, "Error unmarshalling request body: %v", err)
//		return
//	}
//
//	db, err := fs.Client(ctx, "sendgrid")
//	if err != nil {
//		fmt.Fprintf(w, "Error initializing Firebase app: %v", err)
//		return
//	}
//
//	_, _, err = db.Collection("contents").Add(ctx, content)
//	if err != nil {
//		fmt.Fprintf(w, "Error adding content to Firestore: %v", err)
//		return
//	}
//
//	// Call the flaky service
//	flakyServiceURL := fmt.Sprintf("https://%s-%s.cloudfunctins.net/flaky", gcp.Must(gcp.Region()), gcp.Must(gcp.ProjectId()))
//	resp, err := http.Post(flakyServiceURL, "application/json", body)
//	if err != nil {
//		fmt.Fprintf(w, "Error calling flaky service: %v", err)
//		return
//	}
//	defer resp.Body.Close()
//	body, err = io.ReadAll(resp.Body)
//	if err != nil {
//		fmt.Fprintf(w, "Error reading response from flaky service: %v", err)
//		return
//	}
//	//fmt.Fprintf
//}

//func shouldSendWithLease(ctx context.Context, db *firestore.Client, emailRef string) error {
//	return db.RunTransaction(ctx, func(tx *firestore.Transaction) error {
//		docSnapshot, err := tx.Get(emailRef)
//		if err != nil {
//			return false, fmt.Errorf("failed to get email document: %v", err)
//		}
//
//		if !docSnapshot.Exists() {
//			// Email not found, consider sending
//			return true, nil
//		}
//
//		data, err := docSnapshot.Data()
//		if err != nil {
//			return false, fmt.Errorf("failed to decode email document data: %v", err)
//		}
//
//		// Check if email already sent
//		if sent, ok := data["sent"].(bool); ok && sent {
//			return false, nil
//		}
//
//		// Check lease validity
//		now := time.Now()
//		if lease, ok := data["lease"].(time.Time); ok && now.Before(lease) {
//			return false, fmt.Errorf("lease already taken, try later.")
//		}
//
//		// Set lease and allow sending
//		tx.Set(emailRef, map[string]interface{}{"lease": now.Add(LeaseTime)})
//		return true, nil
//	})
//}
