#!/usr/bin/env sh

go install github.com/cosmtrek/air@latest

air

## Ignore below
# install curl
# apk --no-cache add curl

# install live reload
# curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# air