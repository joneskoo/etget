#!/bin/bash
# Work around go test -coverprofile not being able to test multiple packages
set -e

rm -f coverage.txt profile.out

echo 'mode: atomic' > coverage.txt
go list ./... \
    | xargs -I% sh -c 'go test -covermode=atomic -coverprofile=profile.out % && tail -n +2 profile.out >> coverage.txt'
rm -f profile.out
