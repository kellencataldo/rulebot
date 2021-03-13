#!/bin/bash

set -e

go build ../cmd/rulebot.go
retVal=$?
if [ $retVal -ne 0 ]; then
    exit
fi

ssh -i $SECRETS/aws-instance.pem -o StrictHostKeyChecking=no $AWS_INSTANCE "ps -A | grep rulebot | awk '{print \$1}' | xargs kill -2"

scp -i $SECRETS/aws-instance.pem rulebot $AWS_INSTANCE:/home/ec2-user/rulebot/rulebot

ssh -i $SECRETS/aws-instance.pem -o StrictHostKeyChecking=no $AWS_INSTANCE "nohup /home/ec2-user/rulebot/rulebot \
    -discord-token $DISCORD_TOKEN \
    -google-cse $GOOGLE_CSE_ID \
    -google-token $GOOGLE_TOKEN \
    -kellen $KELLEN_ID \
    -rulebooks /home/ec2-user/rulebot/rulebooks \
    -cache /home/ec2-user/rulebot/cache/cache.json \
    -logfile /home/ec2-user/rulebot/logs/rulebot.log 1>/home/ec2-user/rulebot/logs/nohup.log 2>/home/ec2-user/rulebot/logs/nohup.log &"

rm rulebot
