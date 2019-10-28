#!/bin/bash
set -e
set -o pipefail

DUMPFILE=""

main(){

    cat "$DUMPFILE" | grep "main" | grep -oP "CALL (\b([a-z]|\.|\/|[A-Z]|\_|[0-9])+\b)+" | awk -F' ' '{ print $2 }'

}

DUMPFILE="caps-dump"

main $@