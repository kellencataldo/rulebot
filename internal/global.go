package internal

import (
	"flag"
	"fmt"
	"sync"
)

type SourcePage struct {
	Rulebook string `json:"rulebook"`
	Page     int    `json:"page"`
}

type Cache map[string][]SourcePage

// pretty much anything that needs an init goes here.

var (
	Output       string
	DiscordToken string
	GoogleToken  string
	GoogleCSE    string
	Rulebooks    string
	CacheFile    string
	KellenTag    string

	TALK_TO_KELLEN string
	HELP_STRING    string

	SourceCache Cache
	cacheLock   = sync.RWMutex{}

	PrefixMap map[string]string = make(map[string]string)
)

func (c *Cache) Get(page string) ([]SourcePage, bool) {
	if nil == c {
		return []SourcePage{}, false
	}

	cacheLock.RLock()
	defer cacheLock.RUnlock()
	sources, ok := (*c)[page]
	return sources, ok
}

func (c *Cache) Put(page string, sources []SourcePage) {

	if nil != c {
		cacheLock.Lock()
		defer cacheLock.Unlock()
		(*c)[page] = sources
	}
}

func init() {

	flag.StringVar(&Output, "logfile", "", "Logfile for debugging (default outputs to stdout")
	flag.StringVar(&DiscordToken, "discord-token", "", "API token for rulebot to connect to discord")
	flag.StringVar(&GoogleToken, "google-token", "", "API token for google searches")
	flag.StringVar(&GoogleCSE, "google-cse", "", "Custom search engine for google searches")
	flag.StringVar(&Rulebooks, "rulebooks", ".", "Location of directory containing rule book images")
	flag.StringVar(&CacheFile, "cache", "cache.json", "Location of the cache file used by crawler")
	kellenID := flag.String("kellen", "Kellen", "How to tag Kellen in messages")
	flag.Parse()
	KellenTag = fmt.Sprintf("<@%s>", *kellenID)

	PrefixMap["Core Rulebook"] = "core"
	PrefixMap["Advanced Player's Guide"] = "apg"
	PrefixMap["Agents of Edgewatch Player's Guide"] = "aoepg"
	PrefixMap["Gamemastery Guide"] = "gamemaster guide"
	PrefixMap["Character Guide"] = "character guide"
	// hack below, since you can't guarantee map interation order, prefix must be unique.
	PrefixMap["Bestiary pg"] = "bestiary 1"
	PrefixMap["Bestiary 2"] = "bestiary 2"

	TALK_TO_KELLEN = fmt.Sprintf("Something went wrong processing the search, tell %s to check the logs", KellenTag)
	HELP_STRING = ">>> \nRulebot usage: \t!rulebot [options] search terms\n\nOptions are prefixed with a forward slash and must be a non-interrupted string (IE no spaces).\n" +
		"After the first non option string everything will be treated as a search term so options must come first!\n\n" +
		"Options are listed as follows:\n\t/ld=[number]\t\t(Link Depth, default 1) Use this option to specify the number of links the bot will traverse looking for a topic\n" +
		"\t/sd=[number]\t\t(Source Depth, default 3) Use this option to specify the number of source images the bot will post when it finds a topic\n" +
		"\t/core=[number]\t\t Use this option to post a specfic page from the core rule book\n" +
		"\t/apg=[number]\t\t Use this option to post a specific page from the advanced player's guide\n" +
		"\t/aoepg=[number]\t\t Use this option to post a specific page from the agents of edgewatch players guide\n" +
		"\t/tail=[number] This option will post the specified amount of numerically subsequent pages from whatever source page it finds\n" +
		"\t/help\t\t Just prints this help message, also running !rulebot with no arguments will do the same things\n\n" +
		"An example query is as follows: **!rulebot /ld=3 /sd=2 animal companions**\n" +
		"When given the above query, the bot will traverse three links (if it finds that many) and post two source images (if there are that many) from each topic from those links\n\n" +
		"Try to be as specific as possible with your searches\n" +
		fmt.Sprintf("If you find a bug tell %s about it\n", KellenTag)
}
