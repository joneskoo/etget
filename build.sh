#!/bin/bash

PACKAGE=github.com/joneskoo/etget
NAME=`basename $PACKAGE`

export GOARCH=amd64

for GOOS in linux windows darwin; do
    export GOOS
    binary=`basename $PACKAGE`-$GOOS-$GOARCH
    if [[ $GOOS == "windows" ]]; then
        binary=$binary.exe
    fi
    go build -o "$binary" $PACKAGE
done
