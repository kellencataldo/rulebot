#!/bin/bash

go build ../cmd/rulebot.go
retVal=$?
if [ $retVal -ne 0 ]; then
    exit
fi

ssh -i $SECRETS/aws-rulebot.pem -o StrictHostKeyChecking=no $AWS_INSTANCE "ps -A | grep rulebot | awk '{print \$1}' | xargs kill -2"

scp -i $SECRETS/aws-rulebot.pem rulebot $AWS_INSTANCE:/home/ec2-user/rulebot

ssh -i $SECRETS/aws-rulebot.pem -o StrictHostKeyChecking=no $AWS_INSTANCE "nohup /home/ec2-user/rulebot \
    -discord-token $DISCORD_TOKEN \
    -google-cse $GOOGLE_CSE_ID \
    -google-token $GOOGLE_TOKEN \
    -kellen $KELLEN_ID \
    -rulebooks /home/ec2-user/rulebooks \
    -cache /home/ec2-user/cache/cache.json \
    -logfile /home/ec2-user/logs/rulebot.log 1>/home/ec2-user/nohup.log 2>/home/ec2-user/nohup.log &"

rm rulebot
