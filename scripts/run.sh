#!/bin/bash

svcName=${1}

if [ -d "app/${svcName}" ]; then
  cd app/${svcName} && GO_ENV=dev go run .
fi
