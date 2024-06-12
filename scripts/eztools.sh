#!/usr/bin/env bash
BIN="bin"
TMP="$BIN/tmp"
URL="https://f001.backblazeb2.com/file/EricZimmermanTools/net6"

echo "Download tools"

download () {
    if [ ! -f "$BIN/$1.dll" ] ; then
        echo "--- $1"
        curl -s "$URL/$1.zip" | busybox unzip -qq -o -d $TMP -
    else
        return 1
    fi
}

install () {
    cp "$TMP/$1.dll" $BIN
    cp "$TMP/$1.runtimeconfig.json" $BIN
}

installMap () {
    cp -r "$TMP/$2/Maps" $BIN
    install "$2/$1"
}

mkdir -p $BIN

download "EvtxECmd" && installMap "EvtxECmd" "EvtxeCmd"
download "JLECmd" && install "JLECmd"

rm -rf $TMP

export EZTOOLS=$(realpath $BIN)
