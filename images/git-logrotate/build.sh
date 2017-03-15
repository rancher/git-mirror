#!/bin/bash

ACCT=${ACCT:-llparse}
NAME=git-logrotate
VERS=0.1

docker build -t $ACCT/$NAME:$VERS .
docker push $ACCT/$NAME:$VERS
