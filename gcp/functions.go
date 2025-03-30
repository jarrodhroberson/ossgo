package gcp

import (
	"context"
	"strings"

	"cloud.google.com/go/compute/metadata"

	errs "github.com/jarrodhroberson/ossgo/errors"
	strs "github.com/jarrodhroberson/ossgo/strings"
)

// Region retrieves the current region of the Cloud Run instance from metadata
func Region() (string, error) {
	// Check if running on a GCE or Cloud Run env
	if !metadata.OnGCE() {
		return "", errs.MustNeverError.New("not running on a GCE/GAE or Cloud Run env")
	}
	ctx := context.Background()
	// Fetch the zone from the metadata server
	// This fetches a string like "us-central1-f"
	zone, err := metadata.GetWithContext(ctx, "instance/zone")
	if err != nil {
		return strs.NO_DATA, errs.NotFoundError.New("failed to retrieve zone from metadata: %v", err)
	}

	// Extract the region from the zone
	// Metadata provides the zone in the format "projects/<project-number>/zones/<zone-name>"
	parts := strings.Split(zone, "/")
	if len(parts) < 4 {
		return "", errs.ParseError.New("unexpected format for zone: %s", zone)
	}
	zoneName := parts[len(parts)-1]
	region := zoneName[:strings.LastIndex(zoneName, "-")]

	return region, nil
}

func ProjectId() (string, error) {
	// Check if the code is running within a Google Cloud env.
	if !metadata.OnGCE() {
		return strs.NO_DATA, errs.MustNeverError.New("not running on Google Compute Engine, Google App Engine, or Cloud Run env")
	}

	ctx := context.Background()
	// Retrieve the project ID from the metadata server.
	projectID, err := metadata.GetWithContext(ctx, "project/project-id")
	if err != nil {
		return strs.NO_DATA, errs.NotFoundError.New("failed to retrieve project ID from metadata: %w", err)
	}

	return projectID, nil
}
