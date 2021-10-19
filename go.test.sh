#!/usr/bin/env bash

set -e
echo "" > coverage.txt

go mod download github.com/mats9693/utils

for d in $(go list ./... | grep -v vendor); do
    go test -race -coverprofile=profile.out -covermode=atomic "$d"
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done