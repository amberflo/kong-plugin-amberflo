#!/bin/bash

# setup a service with a route
./admin.sh POST /services -d '{"name": "mockbin", "url": "http://mockbin.org"}'
./admin.sh POST /services/mockbin/routes -d '{"name": "mock", "paths": ["/mock"]}'

# setup a user for authentication
./admin.sh POST /consumers -d '{"username": "mike"}'
./admin.sh POST /consumers/mike/key-auth -d '{"key": "super-secret-key"}'

# enable authentication globally
./admin.sh POST /plugins -d '{"name": "key-auth", "config": {"key_names": ["x-api-key"], "key_in_query": false}}'
