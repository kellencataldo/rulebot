package rulebot

import (

    "fmt"

    dg "github.com/bwmarrin/discordgo"
)


func MessageCreate(session* dg.Session, message* dg.MessageCreate) {
    fmt.Println(m.Content)
}
