#!/bin/bash

PLATFORMS="darwin/amd64"
PLATFORMS="$PLATFORMS windows/amd64 windows/386"
PLATFORMS="$PLATFORMS linux/amd64 linux/386"

version="${1?Specify version (first argument)}"

for PLATFORM in $PLATFORMS; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    BIN_FILENAME="bin/pwgenapi-${version}-${GOOS}-${GOARCH}"
    if [[ "${GOOS}" == "windows" ]]; then
        BIN_FILENAME="${BIN_FILENAME}.exe";
    fi
    CMD="GOOS=${GOOS} GOARCH=${GOARCH} go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -o ${BIN_FILENAME}"
    echo "${CMD}"
    eval $CMD || FAILURES="${FAILURES} ${PLATFORM}"
done

# eval errors
if [[ "${FAILURES}" != "" ]]; then
    echo ""
    echo "${SCRIPT_NAME} failed on: ${FAILURES}"
    exit 1
fi