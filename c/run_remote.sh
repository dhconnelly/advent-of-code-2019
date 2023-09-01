#!/bin/bash

if [ -z "$REMOTE_DEPLOY_HOST" ]; then
    echo "missing variable: REMOTE_DEPLOY_HOST"
    exit 1
fi

if [ -z "$REMOTE_DEPLOY_PATH" ]; then
    echo "missing variable: REMOTE_DEPLOYPATH"
    exit 1
fi

if [ -z "$1" ]; then
    echo "usage: run_remote.sh dayN"
    exit 1
fi

HOST=$REMOTE_DEPLOY_HOST
ROOT=advent-of-code-2019
DAY=$1

echo "deploying to $HOST..." &&
    make clean &&
    ssh $HOST "mkdir -p $ROOT" &&
    scp -r * $HOST:$ROOT/ &&
    echo "building on $HOST..." &&
    ssh $HOST "cd $ROOT && make clean && make all" &&
    echo "executing on $HOST..." &&
    ssh $HOST "cd $ROOT && time ./target/$DAY inputs/$DAY.txt" &&
    echo "done."
