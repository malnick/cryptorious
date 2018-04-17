#!/bin/bash
SOURCE_DIR=$(git rev-parse --show-toplevel)

function cleanup() {
	docker rmi -f test 
}
trap cleanup INT HUP TERM KILL 

echo "building container for testing"
docker build -t test $SOURCE_DIR
docker run -i --name test --rm test make test
CODE=$?
cleanup
exit $CODE
