package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	gcs "google.golang.org/api/customsearch/v1"
	gapi "google.golang.org/api/googleapi/transport"
)

func populateWebpages(ctx context.Context, opts Options) []string {

	searchQuery := strings.Join(opts.SearchTerms, "")

	client := &http.Client{Transport: &gapi.APIKey{Key: GoogleToken}}

	svc, err := gcs.New(client)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := svc.Cse.List().Cx(GoogleCSE).Q(searchQuery).Do()
	if err != nil {
		log.Fatal(err)
	}

	for i, result := range resp.Items {
		fmt.Printf("#%d: %s\n", i+1, result.Title)
		fmt.Printf("\t%s\n", result.Snippet)
	}

	return []string{"cool"}
}
