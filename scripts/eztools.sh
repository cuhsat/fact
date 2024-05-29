#!/usr/bin/env bash
BIN="bin"

mkdir -p ${BIN}

echo "Download tools"

if [ ! -f ${BIN}/EvtxECmd.dll ] ; then
    echo "    EvtxECmd"

    wget -q "https://f001.backblazeb2.com/file/EricZimmermanTools/net6/EvtxECmd.zip" -O ${BIN}/evtx.zip
    unzip -q ${BIN}/evtx.zip -d ${BIN}

    cp ${BIN}/EvtxeCmd/EvtxECmd.dll ${BIN}
    cp ${BIN}/EvtxeCmd/EvtxECmd.runtimeconfig.json ${BIN}
    cp -r ${BIN}/EvtxeCmd/Maps ${BIN}

    rm -rf ${BIN}/EvtxeCmd*
    rm -f ${BIN}/evtx.zip
fi

export EZTOOLS=$(realpath ${BIN})
