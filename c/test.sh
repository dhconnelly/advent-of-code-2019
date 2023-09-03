#!/usr/bin/env bash

if make -s test && make -s runall && diff snapshot.txt <(make -s runall); then
    echo "SUCCESS"
    exit 0
else
    echo "FAILURE"
    exit 1
fi
