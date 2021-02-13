package internal

import "flag"

var (
	Output       string
	DiscordToken string
	GoogleToken  string
	GoogleCSE    string
	Rulebooks    string
)

func init() {

	flag.StringVar(&Output, "logfile", "", "Logfile for debugging (default outputs to stdout")
	flag.StringVar(&DiscordToken, "discord-token", "", "API token for rulebot to connect to discord")
	flag.StringVar(&GoogleToken, "google-token", "", "API token for google searches")
	flag.StringVar(&GoogleCSE, "google-cse", "", "custom search engine for google searches")
	flag.StringVar(&Rulebooks, "rulebooks", ".", "Location of directory containing rule book images")

	flag.Parse()
}
