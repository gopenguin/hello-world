#!/bin/bash

SRC_DIR=$1
BUILD_DIR=build
BASE_NAME=$2
BUILD_PARAMS="-a -tags netgo"

mkdir -p $BUILD_DIR

PAIRS=$(go tool dist list | grep -v android | grep -v darwin/arm)

for pair in $PAIRS; do
    export GOOS=$(echo "$pair" | cut -f1 -d/)
    export GOARCH=$(echo "$pair" | cut -f2 -d/)

    echo -n "Building $GOOS/$GOARCH ..."
    go build $BUILD_PARAMS -o $BUILD_DIR/$BASE_NAME-$GOOS-$GOARCH $SRC_DIR
    echo "done"
done

