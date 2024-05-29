#!/usr/bin/env bash
BIN="bin"
GO="go"
GOFLAGS="build -v -race"
VERSION=$(git describe --tags --abbrev=0)
LDFLAGS="-X 'github.com/cuhsat/fact/internal/fact.Version=${VERSION}'"

mkdir -p ${BIN}

echo "Build ${VERSION}"

for DIR in cmd/*/ ; do
    CMD=$(basename ${DIR})

    echo "--- ${CMD}"

    ${GO} ${GOFLAGS} -ldflags "${LDFLAGS}" -o ${BIN}/${CMD} ${DIR}/main.go
done
