package internal

import (
	//	"errors"
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

type Options struct {
	LinkDepth   int
	SourceDepth int
	SearchQuery string
}

const (
	LINK_DEPTH_DEFAULT   = 1
	SOURCE_DEPTH_DEFAULT = 1

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
	opts := Options{LINK_DEPTH_DEFAULT, SOURCE_DEPTH_DEFAULT, ""}

	if len(terms) == 1 {
		log.Println("No arguments specified, returning help string")
		return opts, HELP_STRING
	}

	terms = terms[1:]
	for index, term := range terms {

		if !strings.HasPrefix(term, "/") {

			// treat the rest of the terms as search terms. if there was a flag in there, thats on them.
			opts.SearchQuery = strings.Join(terms[index:], " ")
			break

		} else if strings.HasPrefix(term, "/HELP") {
			log.Println("Explicit HELP argument found, returning help string")
			return opts, HELP_STRING

		} else if strings.HasPrefix(term, "/LD=") {
			var success bool
			opts.LinkDepth, success = parseIntOption(term)
			if !success {
				log.Printf("Unable to parse /LD option string: %s\n", term)
				return opts, "Malformed /LD option, run **!rulebot /HELP** to see proper formatting"
			}

		} else if strings.HasPrefix(term, "/SD=") {
			var success bool
			opts.SourceDepth, success = parseIntOption(term)
			if !success {
				log.Printf("Unable to parse /SD option string: %s\n", term)
				return opts, "Malformed /SD option, run **!rulebot /HELP** to see proper formatting"
			}

		} else {
			// unknown option here, just bail
			log.Printf("Unknown option: %s, returning error\n", term)
			return opts, fmt.Sprintf("Unknown option specified: %s, type **!rulebot /HELP** for supported options", term)
		}
	}

	return opts, ""
}

func convertToFilenames(sources []SourcePage) []string {

	files := make([]string, 0, len(sources))
	for _, source := range sources {
		files = append(files, path.Join(Rulebooks, source.Rulebook+source.Page+FILE_EXTENSION))
	}

	return files
}

func MessageCreate(session *dg.Session, message *dg.MessageCreate) {

	if !strings.HasPrefix(message.Content, "!rulebot") {
		return
	}

	log.Printf("handling message: %s\n", message.Content)
	opts, responseString := populateOptions(message.Content)

	// ehhh, i could be smarter about this lol
	if "" != responseString {
		log.Println("Error parsing options, or HELP specified")
		session.ChannelMessageSend(message.ChannelID, responseString)
		return
	}

	log.Printf("Options for message: %+v\n", opts)

	ctx := context.TODO()
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
	}

	files := convertToFilenames(sources)
	for _, filename := range files {
		f, err := os.Open(filename)
		if nil != err {
			log.Printf("Error opening file for attachment: %s\n", err)
			session.ChannelMessageSend(message.ChannelID, TALK_TO_KELLEN)
			continue
		}

		func() {
			defer f.Close()
			_, err := session.ChannelFileSend(message.ChannelID, filename, f)
			if nil != err {
				log.Printf("Error sending file: %s\n", err)
				session.ChannelMessageSend(message.ChannelID, TALK_TO_KELLEN)
				return
			}

			log.Printf("Successfully sent file: %s\n", filename)
		}()
	}

	log.Println("Finished handling request")
}
