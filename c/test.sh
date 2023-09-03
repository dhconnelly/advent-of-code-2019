#!/usr/bin/env bash

make test && make runall && diff snapshot.txt <(make -s runall)
