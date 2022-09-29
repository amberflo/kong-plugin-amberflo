#!/bin/bash

base_url='http://localhost:8001'

call() {
    local method="$1"
    shift
    local path="$1"
    shift
    >&2 echo '#' "$method" "$base_url$path" "$@"
    curl -s -H 'content-type: application/json' -X "$method" "$base_url$path" "$@" | jq
}

call "$@"
