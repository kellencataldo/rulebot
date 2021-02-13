package internal

import (
	"fmt"
	"log"
	"net/http"

	gcs "google.golang.org/api/customsearch/v1"
	gapi "google.golang.org/api/googleapi/transport"
)

const (
	apiKey = "some-api-key"
	cx     = "some-custom-search-engine-id"
	query  = "some-custom-query"
)

func findLinks() {
	client := &http.Client{Transport: &transport.APIKey{Key: apiKey}}

	svc, err := customsearch.New(client)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := svc.Cse.List().Cx(cx).Q(query).Do()
	if err != nil {
		log.Fatal(err)
	}

	for i, result := range resp.Items {
		fmt.Printf("#%d: %s\n", i+1, result.Title)
		fmt.Printf("\t%s\n", result.Snippet)
	}
}
