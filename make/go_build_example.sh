#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
DESTDIR=${SCRIPT_DIR}/../build

for file in "$@"
do
    EXE=$(basename "$(dirname "$file")")
    go build -o "$DESTDIR"/"$EXE" "$file" ;\
        [ -e "$DESTDIR"/"$EXE" ] && chmod 755 "$DESTDIR"/"$EXE"
done
