#!/bin/bash

pushd ..
go build -o repo-config
popd || exit
jsonResult=$(../repo-config collect --json ./test-fconfig.json)

# echo $jsonResult | jq .

status=$(echo "$jsonResult" | jq -r .status)
if [[ "$status" == "ok" ]]; then
    env_file=$(echo "$jsonResult" | jq -r .env_file)
    # shellcheck disable=SC1090
    source "$env_file"
else
     echo -n "Error: " 
     message=$(echo "$jsonResult" | jq -r .message)
     echo "$message"   
fi
