#!/usr/bin/env bash
GOBIN="bin"

mkdir -p ${GOBIN}

echo "Download tools"

if [ ! -f ${GOBIN}/EvtxECmd.dll ] ; then
    echo "  EvtxECmd"

    wget -q "https://f001.backblazeb2.com/file/EricZimmermanTools/net6/EvtxECmd.zip" -O ${GOBIN}/evtx.zip
    unzip -q ${GOBIN}/evtx.zip -d ${GOBIN}

    cp ${GOBIN}/EvtxeCmd/EvtxECmd.dll ${GOBIN}
    cp ${GOBIN}/EvtxeCmd/EvtxECmd.runtimeconfig.json ${GOBIN}
    cp -r ${GOBIN}/EvtxeCmd/Maps ${GOBIN}

    rm -rf ${GOBIN}/EvtxeCmd*
    rm -f ${GOBIN}/evtx.zip
fi

export EZTOOLS=$(realpath ${GOBIN})
