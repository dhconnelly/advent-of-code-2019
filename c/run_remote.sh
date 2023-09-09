#!/usr/bin/env bash

if [ -z "$REMOTE_DEPLOY_HOST" ]; then
    echo "missing variable: REMOTE_DEPLOY_HOST"
    exit 1
fi

if [ -z "$REMOTE_DEPLOY_PATH" ]; then
    echo "missing variable: REMOTE_DEPLOYPATH"
    exit 1
fi

HOST=$REMOTE_DEPLOY_HOST
ROOT=deploy/advent-of-code-2019

echo "deploying to $HOST:$ROOT..." &&
    make clean &&
    ssh $HOST "mkdir -p $ROOT" &&
    scp -r * $HOST:$ROOT/ &&
    echo "building on $HOST..." &&
    ssh $HOST "cd $ROOT && make clean && make all" &&
    echo "executing on $HOST..." &&
    ssh $HOST "cd $ROOT && time ./test.sh" &&
    echo "done."
