package internal

import (
	//	"errors"
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

type Options struct {
	LinkDepth   int
	SourceDepth int
	SearchQuery string
	Tail        int

	SpecificPageRB string
	SpecificPage   int
}

const (
	LINK_DEPTH_DEFAULT   = 1
	SOURCE_DEPTH_DEFAULT = 1
	TAIL_DEFAULT         = 0

	FILE_EXTENSION = ".png"
)

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

func populateOptions(content string) (Options, string) {

	terms := strings.Fields(content)
	opts := Options{LINK_DEPTH_DEFAULT, SOURCE_DEPTH_DEFAULT, "", TAIL_DEFAULT, "", -1}

	if len(terms) == 1 {
		log.Println("No arguments specified, returning help string")
		return opts, HELP_STRING
	}

	terms = terms[1:]
	for index, term := range terms {

		switch {

		case !strings.HasPrefix(term, "/"):

			// treat the rest of the terms as search terms. if there was a flag in there, thats on them.
			opts.SearchQuery = strings.Join(terms[index:], " ")
			break

		case strings.HasPrefix(term, "/help"):
			log.Println("Explicit HELP argument found, returning help string")
			return opts, HELP_STRING

		case strings.HasPrefix(term, "/ld="):
			var success bool
			opts.LinkDepth, success = parseIntOption(term)
			if !success {
				log.Printf("Unable to parse /tail option string: %s\n", term)
				return opts, "Malformed /tail option, run **!rulebot /help** to see proper formatting"
			}

		case strings.HasPrefix(term, "/sd="):
			var success bool
			opts.SourceDepth, success = parseIntOption(term)
			if !success {
				log.Printf("Unable to parse /sd option string: %s\n", term)
				return opts, "Malformed /sd option, run **!rulebot /help** to see proper formatting"
			}

		case strings.HasPrefix(term, "/tail="):
			var success bool
			opts.Tail, success = parseIntOption(term)
			if !success {
				log.Printf("Unable to parse /tail option string: %s\n", term)
				return opts, "Malformed /tail option, run **!rulebot /help** to see proper formatting"
			}

		case strings.HasPrefix(term, "/apg="):
			opts.SpecificPageRB = "apg"
			var success bool
			opts.SpecificPage, success = parseIntOption(term)
			if !success {
				log.Printf("Unable to parse /apg option string: %s\n", term)
				return opts, "Malformed /apg option, run **!rulebot /help** to see proper formatting"
			}

		case strings.HasPrefix(term, "/aoepg="):
			opts.SpecificPageRB = "aoepg"
			var success bool
			opts.SpecificPage, success = parseIntOption(term)
			if !success {
				log.Printf("Unable to parse /aoepg option string: %s\n", term)
				return opts, "Malformed /aoepg option, run **!rulebot /help** to see proper formatting"
			}

		case strings.HasPrefix(term, "/core="):
			opts.SpecificPageRB = "core"
			var success bool
			opts.SpecificPage, success = parseIntOption(term)
			if !success {
				log.Printf("Unable to parse /core option string: %s\n", term)
				return opts, "Malformed /core option, run **!rulebot /help** to see proper formatting"
			}

		default:
			// unknown option here, just bail
			log.Printf("Unknown option: %s, returning error\n", term)
			return opts, fmt.Sprintf("Unknown option specified: %s, type **!rulebot /help** for supported options", term)
		}
	}

	return opts, ""
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
	case "gamemaster guide" == rulebook:
		fallthrough
	case "character guide" == rulebook:
		return true
	default:
		return false
	}
}

func MessageCreate(session *dg.Session, message *dg.MessageCreate) {

	ctx := context.TODO()
	content := strings.ToLower(message.Content)
	if !strings.HasPrefix(content, "!rulebot") {
		return
	}

	log.Printf("handling message: %s\n", content)
	opts, responseString := populateOptions(content)

	log.Printf("Options for message: %+v\n", opts)

	// ehhh, i could be smarter about this lol
	if "" != responseString {
		log.Println("Error parsing options, or HELP specified")
		session.ChannelMessageSend(message.ChannelID, responseString)
		return
	} else if opts.SpecificPageRB != "" {

		for i := 0; i <= opts.Tail; i++ {

			filename := path.Join(Rulebooks, opts.SpecificPageRB+strconv.Itoa(opts.SpecificPage+i)+FILE_EXTENSION)
			log.Printf("Sending specific file: %s\n", filename)
			sendFile(filename, message.ChannelID, session)
		}

		return
	}

	webpages, success := populateWebpages(ctx, opts.SearchQuery, int64(opts.LinkDepth))
	if !success {
		session.ChannelMessageSend(message.ChannelID, TALK_TO_KELLEN)
		return
	} else if 0 == len(webpages) {
		session.ChannelMessageSend(message.ChannelID, "No results found, broaden your search")
		return
	}

	sources, ok := crawlLinks(ctx, webpages, opts.SourceDepth)
	if !ok {
		log.Fatalln("error occured while crawling links")
		session.ChannelMessageSend(message.ChannelID, TALK_TO_KELLEN)
		return
	} else if 0 == len(sources) {
		session.ChannelMessageSend(message.ChannelID, "No results found, broaden your search")
		return
	}

	for _, source := range sources {

		if isHiddenRulebook(source.Rulebook) {
			session.ChannelMessageSend(message.ChannelID, source.Rulebook+" pg. "+strconv.Itoa(source.Page))
			continue
		}

		for i := 0; i <= opts.Tail; i++ {
			filename := path.Join(Rulebooks, source.Rulebook+strconv.Itoa(source.Page+i)+FILE_EXTENSION)
			log.Printf("Sending discovered file: %s\n", filename)
			sendFile(filename, message.ChannelID, session)
		}
	}
}
