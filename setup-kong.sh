#!/bin/bash

admin_url='http://localhost:8001'

admin() {
    local method="$1"
    shift
    local path="$1"
    shift
    curl -s -H 'content-type: application/json' -X "$method" "$admin_url$path" "$@" | jq
}

# setup a service with a route
admin POST /services -d '{"name": "mockbin", "url": "http://mockbin.org"}'
admin POST /services/mockbin/routes -d '{"name": "mock", "paths": ["/mock"]}'

# setup a user for authentication
admin POST /consumers -d '{"username": "mike"}'
admin POST /consumers/mike/key-auth -d '{"key": "super-secret-key"}'

# enable authentication globally
admin POST /plugins -d '{"name": "key-auth", "config": {"key_names": ["x-api-key"], "key_in_query": false}}'
