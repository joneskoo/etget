#!/bin/bash
# Work around go test -coverprofile not being able to test multiple packages
set -e

rm -f coverage.txt profile.out

echo 'mode: set' > coverage.txt
go list ./... \
    | xargs -I% bash -c 'go test -coverprofile=profile.out % && tail -n +2 profile.out >> coverage.txt || true'
rm -f profile.out

go tool cover -func coverage.txt

goveralls -coverprofile=coverage.txt -service=travis-ci
