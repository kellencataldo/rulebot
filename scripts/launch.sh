#!/bin/bash

go build ../cmd/rulebot.go
retVal=$?
if [ $retVal -ne 0 ]; then
    exit
fi

./rulebot -discord-token $DISCORD_TOKEN -google-cse $GOOGLE_CSE_ID -google-token $GOOGLE_TOKEN -kellen $KELLEN_ID -rulebooks ~/Desktop/rulebooks

