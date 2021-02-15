package internal

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	gq "github.com/PuerkitoBio/goquery"
	_ "golang.org/x/net/html"
)

func filterRawSources(rawSources []string) []SourcePage {

	sources := make([]SourcePage, 0, len(rawSources))
	for _, rawSource := range rawSources {

		converted := false
		for rawPrefix, filePrefix := range PrefixMap {
			if strings.HasPrefix(rawSource, rawPrefix) {
				tempSplit := strings.Fields(rawSource)
				pageNum, _ := strconv.Atoi(tempSplit[len(tempSplit)-1])
				sources = append(sources, SourcePage{filePrefix, pageNum})
				converted = true
				break
			}
		}

		if converted {
			log.Printf("raw source %s converted to source: %+v\n", rawSource, sources[len(sources)-1])
		} else {
			log.Printf("Not conversion found for raw source: %s\n", rawSource)
		}
	}

	return sources
}

func findSources(ctx context.Context, webpage string, sourceDepth int, ch chan<- SourcePage, wg *sync.WaitGroup) {

	defer wg.Done()

	if pages, ok := SourceCache.Get(webpage); ok {
		log.Printf("Cached source for webpage: %s found\n", webpage)
		for index, source := range pages {
			if index == sourceDepth {
				break
			}
			ch <- source
		}
		return
	}

	resp, err := http.Get(webpage)
	if nil != err {
		log.Printf("goroutine crawling: %s encountered error: %s\n", webpage, err)
		return
	} else if http.StatusOK != resp.StatusCode {
		log.Printf("goroutine crawling: %s got response %s from URL\n", webpage, resp.Status)
		return
	}

	defer resp.Body.Close()
	doc, err := gq.NewDocumentFromReader(resp.Body)
	if nil != err {
		log.Printf("Error parsing html document to reader: %s\n", err)
		return
	}

	rawSources := make([]string, 0)
	doc.Find("i").Each(func(_ int, selection *gq.Selection) {
		rawSources = append(rawSources, selection.Text())
	})

	results := filterRawSources(rawSources)
	SourceCache.Put(webpage, results)
	for index, source := range results {
		if index == sourceDepth {
			break
		}

		ch <- source
	}
}

func crawlLinks(ctx context.Context, webpages []string, sourceDepth int) ([]SourcePage, bool) {

	log.Printf("starting %d goroutines, each finding %d sources\n", len(webpages), sourceDepth)

	ch := make(chan SourcePage, len(webpages)*sourceDepth)
	var wg sync.WaitGroup
	for _, webpage := range webpages {
		wg.Add(1)
		go findSources(ctx, webpage, sourceDepth, ch, &wg)
	}

	searchDone := make(chan struct{})
	go func() {
		defer close(searchDone)
		wg.Wait()
	}()

	select {
	// dont need this now... but maybe 🤔
	case <-ctx.Done():
		log.Println("Context canceled while waiting for search")
		return []SourcePage{}, false
	case <-searchDone:
		log.Println("search completed successfully, processing items")
	}

	sourceSet := make(map[SourcePage]bool)
	close(ch)
	for source := range ch {
		sourceSet[source] = true
	}

	results := make([]SourcePage, 0, len(sourceSet))
	for source, _ := range sourceSet {
		results = append(results, source)
	}

	return results, true
}
