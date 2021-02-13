package internal

import (
	"context"
	"log"
	"net/http"

	gcs "google.golang.org/api/customsearch/v1"
	gapi "google.golang.org/api/googleapi/transport"
)

func populateWebpages(ctx context.Context, searchQuery string, linkDepth int64) ([]string, bool) {

	client := &http.Client{Transport: &gapi.APIKey{Key: GoogleToken}}

	svc, err := gcs.New(client)
	if err != nil {
		log.Fatalf("error building google service, error: %s\n", err)
		return []string{}, false
	}

	resp, err := svc.Cse.List().Cx(GoogleCSE).Num(linkDepth).Q(searchQuery).Do()
	if err != nil {
		log.Fatalf("error executing search: %s, error recieved: %s\n", searchQuery, err)
		return []string{}, false
	}

	links := make([]string, len(resp.Items))
	for i, result := range resp.Items {
		links[i] = result.Link
	}

	return links, true
}
