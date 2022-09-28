#!/bin/bash

base_url='http://localhost:8000'

call() {
    local method="$1"
    shift
    local path="$1"
    shift
    curl -s -H 'content-type: application/json' -X "$method" "$base_url$path" "$@"
}

call "$@" | jq
