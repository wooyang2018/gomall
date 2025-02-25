#!/bin/bash

ps -ef | grep -E '/tmp/go-build|go run \.' | awk '{print $2}' | xargs kill
