package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	dg "github.com/bwmarrin/discordgo"
	rb "github.com/kellencataldo/rulebot/internal"
)

func main() {

	if "" != rb.Output {

		f, err := os.OpenFile(rb.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if nil != err {
			log.Fatalf("Unable to open log file, error recieved: %s\n", err)
		}

		log.SetOutput(f)
		defer f.Close()
	}

	if "" == rb.DiscordToken {
		log.Fatal("No API token for discord specified")
	}

	if "" == rb.GoogleToken {
		log.Fatal("No API token for google searches specified")
	}

	if "" == rb.GoogleCSE {
		log.Fatal("No custom search engine for google specified")
	}

	rb.SourceCache = make(map[string][]rb.SourcePage)
	if _, err := os.Stat(rb.CacheFile); err == nil {
		// cache exists
		log.Println("Cache file discovered, attempting to load")
		file, err := ioutil.ReadFile(rb.CacheFile)
		if err != nil {
			log.Printf("Error reading cache file: %s\n", err)
		} else {
			err = json.Unmarshal([]byte(file), &rb.SourceCache)
			if err != nil {
				log.Printf("Error serializing json from cache file to cache map: %s\n", err)
			} else {
				log.Println("Successfully serialized json from cache file to cache map")
			}
		}
	} else {
		log.Println("No cache file found, will be created at shutdown")
	}

	log.Println("rulebot starting.")
	dSession, err := dg.New("Bot " + rb.DiscordToken)
	if nil != err {
		log.Fatalf("Unable to create discord session, error: %s\n", err)
	}

	dSession.Identify.Intents = dg.IntentsGuildMessages
	dSession.AddHandler(rb.MessageCreate)
	err = dSession.Open()
	if nil != err {
		log.Fatalf("Unable to open websocket to server, error: %s\n", err)
	}

	defer dSession.Close()
	log.Println("Connection established, now listening for messages")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	fmt.Println("Bot running, press CTRL-C to exit")
	<-sc

	log.Println("Exit signal received, storing cache file")

	rawBytes, err := json.MarshalIndent(rb.SourceCache, "", "\t")
	if nil != err {
		log.Fatalf("Error marshalling cache map to bytes: %s\n", err)
	}

	err = ioutil.WriteFile(rb.CacheFile, rawBytes, 0644)
	if nil != err {
		log.Fatalf("Error writing json to cache file: %s\n", err)
	}

	log.Println("Clean exit")
	os.Exit(0)
}
