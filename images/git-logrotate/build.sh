#!/bin/bash

ACCT=${1:-drpebcak}
NAME=${2:-git-logrotate}
VERS=${3:-0.1}

docker build -t $ACCT/$NAME:$VERS .
docker push $ACCT/$NAME:$VERS
