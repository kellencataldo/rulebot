package internal

import (
	//	"errors"
	_ "context"
	"fmt"
	"log"
	"strconv"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

type Options struct {
	LinkDepth   int
	SourceDepth int
	SearchTerms []string
}

const (
	LINK_DEPTH_DEFAULT   = 1
	SOURCE_DEPTH_DEFAULT = 3
	HELP_STRING          = ">>> Rulebot usage: !rulebot [options] search terms\nOptions are prefixed with a forward slash and must be a non-interrupted string (IE no spaces).\n" +
		"After the first non option string everything will be treated as a search term so options must come first!\n\n" +
		"Options are listed as follows:\n\t/LD=[number]\t\t(Link Depth, default 1) Use this option to specify the number of links the bot will traverse looking for a topic\n" +
		"\t/SD=[number]\t\t(Source Depth, default 3) Use this option to specify the number of source images the bot will post when it finds a topic\n" +
		"\t/HELP\t\t Just prints this help message, also running !rulebot with no arguments will do the same things\n\n" +
		"An example query is as follows: **!rulebot /LD=3 /SD=2 animal companions**\n" +
		"When given the above query, the bot will traverse three links (if it finds that many) and post two source images (if there are that many) from each topic from those links\n\n" +
		"Try to be as specific as possible with your searches\n" +
		"If you find a bug tell Kellen about it\n"
)

// all that is supported right now, but maybe more to come???
func parseIntOption(option string) (int, bool) {

	parsed := strings.Split(option, "=")

	if 2 == len(parsed) {
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
	opts := Options{LINK_DEPTH_DEFAULT, SOURCE_DEPTH_DEFAULT, []string{}}

	if len(terms) == 1 {
		log.Println("No arguments specified, returning help string")
		return opts, HELP_STRING
	}

	terms = terms[1:]
	fmt.Println(terms)

	for index, term := range terms {

		if !strings.HasPrefix(term, "/") {

			// treat the rest of the terms as search terms. if there was a flag in there, thats on them.
			opts.SearchTerms = terms[index:]
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
			return opts, fmt.Sprintf("Unknown option specified: %s, type **!ruleboy help** for supported options", term)
		}
	}

	return opts, ""
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
	// ctx := context.TODO()
	// webpages := populateWebpages(ctx, opts)
	// fmt.Println(webpages)

	// main.Rulebooks

	session.ChannelMessageSend(message.ChannelID, HELP_STRING)
}
