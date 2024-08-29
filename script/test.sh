#!/bin/sh
curl \
  --noproxy '*' \
  -v \
  --request POST \
  --location 'http://99.48.167.162:3030' \
  --header 'Content-Type: application/json' \
  --data-raw '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "getHeight",
  }'
