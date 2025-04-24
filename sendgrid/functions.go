package sendgrid

import (
	"context"
	"fmt"
	"net"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jarrodhroberson/ossgo/functions/must"

	"cloud.google.com/go/firestore"

	fs "github.com/jarrodhroberson/ossgo/firestore"
	"github.com/rs/zerolog/log"
)

// IsEmailAddressFormatValid checks if a string is a valid RFC 822 email address, with caveats.
// Returns nil error if valid, and an error with a reason if invalid.
// RFC 822 is a very old standard. This attempts to cover the core requirements, but is not a perfect implementation.
// For complete validation, use a dedicated email validation library/package, or a more thorough approach
// such as parsing the email using the `mail` package and validating parts.
//
// This uses `net/mail.ParseAddress`, which is a more robust base, and then applies a regex
// and checks length and content.
func IsEmailAddressFormatValid(emailAddress string) (bool, error) {
	// Quick check to ensure input is a string and not empty.
	if len(strings.TrimSpace(emailAddress)) == 0 {
		return false, fmt.Errorf("emailAddress is empty: %s", emailAddress)
	}

	// Use net/mail.ParseAddress for basic validation of the address structure.
	// This handles the quoting and other complexities, and is the most important part.
	_, err := mail.ParseAddress(emailAddress)
	if err != nil {
		return false, fmt.Errorf("invalid address structure for emailAddress '%s': %w", emailAddress, err)
	}

	// Further checks (some are redundant but provide extra validation)

	// Regular expression to enforce basic format (very basic, covers most common issues)
	// Covers the "user@domain.tld" format generally.
	// This is a SIMPLIFIED regex. Full RFC 822 validation is EXTREMELY complex.
	// Using more advanced package and library solutions is *always* the right answer if it's important.
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`) // common format

	if !re.MatchString(emailAddress) {
		return false, fmt.Errorf("invalid emailAddress format (regex) for emailAddress '%s': regex: %s", emailAddress, re.String())
	}

	// Additional content and length checks (more defensive/redundant)
	if len(emailAddress) > 254 { // RFC 5321 allows a maximum of 256 characters for an emailAddress address
		return false, fmt.Errorf("emailAddress length exceeds maximum allowed (254 characters) for emailAddress '%s'", emailAddress)
	}

	parts := strings.Split(emailAddress, "@")
	if len(parts) != 2 {
		return false, fmt.Errorf("emailAddress must contain exactly one @ symbol for emailAddress '%s'", emailAddress)
	}

	localPart := parts[0]
	domainPart := parts[1]

	if len(localPart) > 64 { // RFC 5321 allows a maximum of 64 characters for the local-part
		return false, fmt.Errorf("local part '%s' exceeds maximum length (64 characters) for emailAddress '%s'", localPart, emailAddress)
	}

	if len(domainPart) > 255 { //  RFC 5321 allows a maximum of 255 characters for the domain part
		return false, fmt.Errorf("domain part '%s' exceeds maximum length (255 characters) for emailAddress '%s'", domainPart, emailAddress)
	}

	//  Additional domain-specific validation (can be significantly more involved)
	// This checks that the domain is either a valid hostname OR an IP address
	isValid, err := isValidDomain(domainPart)
	if !isValid {
		return false, fmt.Errorf("invalid domain for emailAddress '%s': %w", emailAddress, err) // Wrap the error
	}

	return true, nil // If all checks pass, it is considered valid (with the noted caveats)
}

// isValidDomain checks if a domain is valid either as a hostname or an IP address
func isValidDomain(domain string) (bool, error) {
	// Check for a valid IP address first.  net.ParseIP is a good check.
	if net.ParseIP(domain) != nil {
		return true, nil // It's a valid IP address.
	}

	// Check if it's a valid hostname using net.LookupHost (less reliable, but a basic check)
	if len(domain) > 0 && domain[len(domain)-1] == '.' {
		domain = domain[:len(domain)-1] // Remove trailing dot, if present
	}
	_, err := net.LookupHost(domain)
	if err == nil {
		return true, nil // Hostname lookup successful (basic indication it *could* be valid)
	} else {
		return false, fmt.Errorf("invalid domain '%s': %w", domain, err) // Return the specific domain and the lookup error.
	}
}

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
	client := must.Must(fs.Client(ctx, "sendgrid"))
	defer func(db *firestore.Client) {
		err := db.Close()
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not close firebase client")
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
	client := must.Must(fs.Client(ctx, "sendgrid"))
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not close firebase client")
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
//		docSnapshot, err := tx.Load(emailRef)
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
