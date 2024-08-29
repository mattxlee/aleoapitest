#!/bin/sh
curl \
  -v \
  --request POST \
  --location 'http://54.193.29.190:3030' \
  --header 'Content-Type: application/json' \
  --data-raw '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "getHeight",
  }'
