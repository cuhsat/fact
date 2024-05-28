#!/usr/bin/env bash
GO="go"
GOFLAGS="build -race"
GOBIN="bin"
VERSION=$(git describe --tags --abbrev=0)
LDFLAGS="-X 'github.com/cuhsat/fact/internal/fact.Version=$VERSION'"

mkdir -p ${GOBIN}

echo "Build ${VERSION}"

for DIR in cmd/*/ ; do
    BIN=$(basename $DIR)

    echo "  $BIN"

    ${GO} ${GOFLAGS} -ldflags "$LDFLAGS" -o ${GOBIN}/$BIN $DIR/main.go
done
