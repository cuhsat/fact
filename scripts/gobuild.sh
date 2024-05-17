#!/usr/bin/env bash
GO="go"
GOFLAGS="build -race"
GOBIN="../bin"
VERSION=$(git describe --tags --abbrev=0)

mkdir -p ${GOBIN}

echo "Building version ${VERSION}"

for DIR in ../cmd/*/ ; do
    TOOL=$(basename $DIR)

    echo "$TOOL"

    ${GO} ${GOFLAGS} -ldflags="-X '$DIR/main.Version=${VERSION}'" -o ${GOBIN}/$TOOL $DIR/main.go
done
