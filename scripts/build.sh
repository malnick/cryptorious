#!/bin/bash
set -e

PLATFORM=$(uname | tr [:upper:] [:lower:])
GIT_REF=$(git tag | tail -1)
SOURCE_DIR=$(git rev-parse --show-toplevel)
VERSION=${GIT_REF}
REVISION=$(git rev-parse --short HEAD)
BUILD_DIR="${SOURCE_DIR}/release"
OS=$1

function license_check {
    retval=0
    for source_file in $(find . -type f -name '*.go' -not -path './examples/**' -not -path './vendor/**' -not -path './schema/**'); do
        if ! grep -E 'Copyright [0-9]{4} BlueOwl, LLC.' $source_file &> /dev/null; then
            echo "Missing copyright statement in ${source_file}"
            retval=$((retval + 1))
        fi
        if ! grep 'Licensed under the Apache License, Version 2.0 (the "License");' $source_file &> /dev/null; then
            echo "Missing license header in ${source_file}"
            retval=$((retval + 1))
        fi
    done
    if [[ $retval -gt 0 ]]; then
        echo
        echo "ERROR: found ${retval} cases of missing copyright statements or license headers."
        return $retval
    fi
}

function build_cmd {
  # we might OSS this eventually
  # license_check

	echo "	Version:          ${VERSION}"
	echo "	Revision:         ${REVISION}"
	echo "	Operating System: ${OS}"
	
	GOOS=$OS go build -a -o ${BUILD_DIR}/cryptorious_${OS}_${GIT_REF}_${REVISION}                               \
        -ldflags "-X main.VERSION=${VERSION} -X main.REVISION=${REVISION}" \
        ${SOURCE_DIR}/cryptorious.go
}

function main {
    build_cmd
    if [ -n "$(which tree)" ]; then
        tree $BUILD_DIR 
    fi
}

main "$@"
