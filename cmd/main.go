package main

import (
    "log"
    "flag"
    "os"
    "fmt"

    rb "rulebot"

    dg "github.com/bwmarrin/discordgo"
)


func main() {

    var output string
    var token string

    flag.StringVar(&output, "logfile", "", "Logfile for debugging (default outputs to stdout")
    flag.StringVar(&token, "token", "", "API token for rulebot")

    flag.Parse()

    if "" != output {

        f, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if nil != err {
            log.Fatalf("Unable to open log file, error recieved: %s\n", err)
            os.Exit(1)
        }

        log.SetOutput(f)
        defer f.Close()
    }

    log.Println("rulebot starting.")

    dSession, err := dg.New("Bot " + token)
    if nil != err {
        log.Fatalf("Unable to create discord session, error: %s\n", err)
        os.Exit(1)
    }


    dSession.Identify.Intents = dSession.IntentsGuildMessages
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

