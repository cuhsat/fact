#!/usr/bin/env bash
GO="go"
GOFLAGS="clean"
GOBIN="bin"

rm -rf ${GOBIN}
${GO} ${GOFLAGS}
