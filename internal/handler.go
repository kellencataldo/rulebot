package internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

const (
	LINK_DEPTH_DEFAULT   = 1
	SOURCE_DEPTH_DEFAULT = 1
	TAIL_DEFAULT         = 0

	FILE_EXTENSION = ".png"
)

type (
	Query interface {
		performQuery(context.Context) ([]SourcePage, error)
		getTail() int
	}

	SpecificQuery struct {
		SourcePage
		Tail int
	}

	CrawlerQuery struct {
		LinkDepth   int
		SourceDepth int
		SearchQuery string
		Tail        int
	}
)

func (sq SpecificQuery) getTail() int {
	return sq.Tail
}

func (sq SpecificQuery) performQuery(ctx context.Context) ([]SourcePage, error) {
	return []SourcePage{sq.SourcePage}, nil
}

func (cq CrawlerQuery) getTail() int {
	return cq.Tail
}

func (cq CrawlerQuery) performQuery(ctx context.Context) ([]SourcePage, error) {

	webpages, success := populateWebpages(ctx, cq.SearchQuery, int64(cq.LinkDepth))
	if !success {
		return []SourcePage{}, errors.New(TALK_TO_KELLEN)
	} else if 0 == len(webpages) {
		return []SourcePage{}, errors.New("No results found, broaden your search")
	}

	sources, ok := crawlLinks(ctx, webpages, cq.SourceDepth)
	if !ok {
		log.Println("Error occured while crawling links")
		return sources, errors.New(TALK_TO_KELLEN)
	} else if 0 == len(sources) {
		return sources, errors.New("No results found, broaden your search")
	}

	return sources, nil
}

// all that is supported right now, but maybe more to come???
func parseIntOption(option string) (int, bool) {

	parsed := strings.Split(option, "=")

	if 2 != len(parsed) {
		log.Printf("Unable to parse option: %s, into identifier and value pair", option)
		return 0, false
	}

	value, err := strconv.Atoi(parsed[1])
	if nil != err {
		log.Printf("Unable to convert value: %s to integer\n", parsed[1])
		return 0, false
	}

	return value, true
}

func populateOptions(content string) (Query, error) {

	terms := strings.Fields(content)
	query := CrawlerQuery{LINK_DEPTH_DEFAULT, SOURCE_DEPTH_DEFAULT, "", TAIL_DEFAULT}

	if len(terms) == 1 {
		log.Println("No arguments specified, returning help string")
		return query, errors.New(HELP_STRING)
	}

	terms = terms[1:]

ScannerLoop:
	for index, term := range terms {

		switch {

		case !strings.HasPrefix(term, "/"):
			// treat the rest of the terms as search terms. if there was a flag in there, thats on them.
			query.SearchQuery = strings.Join(terms[index:], " ")
			break ScannerLoop

		case strings.HasPrefix(term, "/help"):
			log.Println("Explicit HELP argument found, returning help string")
			return query, errors.New(HELP_STRING)

		case strings.HasPrefix(term, "/ld="):
			var success bool
			query.LinkDepth, success = parseIntOption(term)
			if !success {
				log.Printf("Unable to parse /tail option string: %s\n", term)
				return query, errors.New("Malformed /tail option, run **!rulebot /help** to see proper formatting")
			}

		case strings.HasPrefix(term, "/sd="):
			var success bool
			query.SourceDepth, success = parseIntOption(term)
			if !success {
				log.Printf("Unable to parse /sd option string: %s\n", term)
				return query, errors.New("Malformed /sd option, run **!rulebot /help** to see proper formatting")
			}

		case strings.HasPrefix(term, "/tail="):
			var success bool
			query.Tail, success = parseIntOption(term)
			if !success {
				log.Printf("Unable to parse /tail option string: %s\n", term)
				return query, errors.New("Malformed /tail option, run **!rulebot /help** to see proper formatting")
			}

		case strings.HasPrefix(term, "/apg="):
			var success bool
			// If a tail option was already encountered, add it here.
			pageQuery := SpecificQuery{SourcePage{Rulebook: "apg"}, query.Tail}
			if pageQuery.SourcePage.Page, success = parseIntOption(term); !success {
				log.Printf("Unable to parse /apg option string: %s\n", term)
				return pageQuery, errors.New("Malformed /apg option, run **!rulebot /help** to see proper formatting")
			}
			return pageQuery, nil

		case strings.HasPrefix(term, "/aoepg="):
			var success bool
			// If a tail option was already encountered, add it here.
			pageQuery := SpecificQuery{SourcePage{Rulebook: "aoepg"}, query.Tail}
			if pageQuery.SourcePage.Page, success = parseIntOption(term); !success {
				log.Printf("Unable to parse /aoepg option string: %s\n", term)
				return query, errors.New("Malformed /aoepg option, run **!rulebot /help** to see proper formatting")
			}
			return pageQuery, nil

		case strings.HasPrefix(term, "/core="):
			var success bool
			// If a tail option was already encountered, add it here.
			pageQuery := SpecificQuery{SourcePage{Rulebook: "core"}, query.Tail}
			if pageQuery.SourcePage.Page, success = parseIntOption(term); !success {
				log.Printf("Unable to parse /core option string: %s\n", term)
				return query, errors.New("Malformed /core option, run **!rulebot /help** to see proper formatting")
			}
			return pageQuery, nil

		default:
			// unknown option here, just bail
			log.Printf("Unknown option: %s, returning error\n", term)
			return query, fmt.Errorf("Unknown option specified: %s, type **!rulebot /help** for supported options", term)
		}
	}

	log.Printf("Option parse successful, options for message: %+v\n", query)
	return query, nil
}

func sendFile(filename, channelID string, session *dg.Session) {

	f, err := os.Open(filename)
	if nil != err {
		log.Printf("Error opening file for attachment: %s\n", err)

		if os.IsNotExist(err) {
			session.ChannelMessageSend(channelID, fmt.Sprintf("Can't find file: %s", filepath.Base(filename)))
			return
		}

		session.ChannelMessageSend(channelID, TALK_TO_KELLEN)
		return
	}

	defer f.Close()
	_, err = session.ChannelFileSend(channelID, filename, f)
	if nil != err {
		log.Printf("Error sending file: %s\n", err)
		session.ChannelMessageSend(channelID, TALK_TO_KELLEN)
		return
	}

	log.Printf("Successfully sent file: %s\n", filename)
}

func isHiddenRulebook(rulebook string) bool {

	switch {
	case "bestiary 1" == rulebook:
		fallthrough
	case "bestiary 2" == rulebook:
		fallthrough
	case "gamemaster guide" == rulebook:
		fallthrough
	case "character guide" == rulebook:
		return true
	default:
		return false
	}
}

func MessageCreate(session *dg.Session, message *dg.MessageCreate) {

	ctx := context.Background()
	content := strings.ToLower(message.Content)
	if !strings.HasPrefix(content, "!") {
		return
	}

	if !strings.HasPrefix(content, "!rulebot") {
		// remove this if another bot is added to the channel lol
		msg := fmt.Sprintf("It's \"!rulebot\" not \"%s\"", strings.Fields(content)[0])
		session.ChannelMessageSend(message.ChannelID, msg)
		return
	}

	log.Printf("handling message: %s\n", content)
	query, err := populateOptions(content)

	if err != nil {
		log.Println("Error parsing options, or HELP specified")
		session.ChannelMessageSend(message.ChannelID, err.Error())
		return
	}

	sources, err := query.performQuery(ctx)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, err.Error())
	}

	// second layer of de-dup (mostly for /tail queries)
	sent := make(map[string]bool)
	for _, source := range sources {
		if isHiddenRulebook(source.Rulebook) {
			session.ChannelMessageSend(message.ChannelID, source.Rulebook+" pg. "+strconv.Itoa(source.Page))
			continue
		}

		for i := 0; i <= query.getTail(); i++ {
			filename := path.Join(Rulebooks, source.Rulebook+strconv.Itoa(source.Page+i)+FILE_EXTENSION)
			if sent[filename] {
				continue
			}

			sent[filename] = true
			log.Printf("Sending discovered file: %s\n", filename)
			sendFile(filename, message.ChannelID, session)
		}
	}
}
