#!/bin/bash

api_key=''

meter_api_name='kong-api-calls'
customer_header='x-consumer-username'

./admin.sh POST /plugins -d "$(cat <<EOF
{
    "name": "metering",
    "config": {
        "apiKey": "$api_key",
        "meterApiName": "$meter_api_name",
        "customerHeader": "$customer_header"
    }
}
EOF
)"
