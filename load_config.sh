#!/bin/bash

#
# this file is where you run the commands to load the configurations in your project.
# there can be many configurations, depending on how the app is structured.  in this case, we are
# just loading the test configurations.
# we use --silent because after the initial collection of data, we do not want to be bothered
# everytime a terminal is opened

# NOTE:  the executable repo_config must be part of the repo for this to work.  in this project
#        we just build it.  in any other project it should be installed through the docker file

if [[ ! -f "repo-config" ]]; then
    go build -o repo-config
fi

output=$(./repo-config collect --json ./end-to-end-test/test-config.json --silent)
result=$(echo "$output" | jq -r .status)

if [[ "$result" != "ok" ]]; then
    echo -n "ERROR: "
    echo "$output" | jq -r .message
else
    echo -n  "config in "
    echo "$output" | jq -r .json_file
fi
