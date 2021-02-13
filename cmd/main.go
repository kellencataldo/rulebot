package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	dg "github.com/bwmarrin/discordgo"
	rb "github.com/kellencataldo/rulebot/internal/rulebot"
)

var (
	output       string
	discordToken string
	GoogleToken  string
	Rulebooks    string
)

func init() {

	flag.StringVar(&output, "logfile", "", "Logfile for debugging (default outputs to stdout")
	flag.StringVar(&discordToken, "discord-token", "", "API token for rulebot to connect to discord")
	flag.StringVar(&GoogleToken, "google-token", "", "API token for google searches")
	flag.StringVar(&GoogleCSE, "google-cse", "", "custom search engine for google searches")
	flag.StringVar(&Rulebooks, "rulebooks", ".", "Location of directory containing rule book images")

	flag.Parse()
}

func main() {

	if "" != output {

		f, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if nil != err {
			log.Fatalf("Unable to open log file, error recieved: %s\n", err)
			os.Exit(1)
		}

		log.SetOutput(f)
		defer f.Close()
	}

	if "" == discordToken {
		log.Fatal("No API token for discord specified")
		os.Exit(1)
	}

	if "" == GoogleToken {
		log.Fatal("No API token for google searches specified")
		os.Exit(1)
	}

	if "" == GoogleCSE {
		log.Fatal("No custom search engine for google specified")
		os.Exit(1)
	}

	log.Println("rulebot starting.")
	dSession, err := dg.New("Bot " + token)
	if nil != err {
		log.Fatalf("Unable to create discord session, error: %s\n", err)
		os.Exit(1)
	}

	dSession.Identify.Intents = dg.IntentsGuildMessages
	dSession.AddHandler(rb.MessageCreate)
	err = dSession.Open()
	if nil != err {
		log.Fatalf("Unable to open websocket to server, error: %s\n", err)
		os.Exit(1)
	}

	defer dSession.Close()
	log.Println("Connection established, now listening for messages")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	fmt.Println("Bot running, press CTRL-C to exit")
	<-sc

	log.Println("Exiting cleanly")
}
