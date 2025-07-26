#!/usr/bin/env bash
BIN="bin"

GO="go"
GOFLAGS="clean"

rm -f *.zip
rm -rf ${BIN}
${GO} ${GOFLAGS}
