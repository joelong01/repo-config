# repo-config
A command-line tool to manage repository configurations by collecting and deleting configuration files.

## Overview
``` repo-config``` is a CLI tool designed to streamline the management of repository configuration files. It allows you to specify the settings you need in your repo without specifying the actual values.  The tool will look at a required settings in a JSON file and collect configuration settings to generate corresponding output files that have the actual values.  The output files are in both an ENV and JSON formats and they are stored off of the $HOME directory so that they will never be checked into the repo.

Additionally, it provides functionality to delete these output files when they are no longer needed.

The output files generated by the collect command are stored in the 

```$HOME/repo-config/<git project>``` directory.
The names of the output files are derived from the input JSON file name.
For example, if your input file is cosmosdb_settings.json and the git project is "purchase_service" the output files will be:

```~/.repo-config/purchase_service/.cosmosdb_settings-values.json```
```~/.repo-config/rpurchase_service/.cosmosdb_settings-values.env```

Ensure that you have the necessary permissions to read and write files in the home directory.

The ```repo-config``` program does all interaction through Stderr, except the final result, which is sent to Stdout as a JSON document.  The format of the document looks like:
```json
{
  "status": "ok",
  "message": "",
  "env_file": "~/.repo-config/purchase_service/.cosmosdb_settings-values.json",
  "json_file": "~/.repo-config/rpurchase_service/.cosmosdb_settings-values.env"
} 
```
or as an example error:
```json
{
  "status": "error",
  "message": "JSON file './test-config.1json' not found",
  "env_file": "",
  "json_file": ""
} 
```

The directory end-to-end-test shows a sample on how to use this JSON document with jq to branch on errors and source the environment.

## Installation
To install repo-config, ensure you have Go installed on your system. Build the executable by running:

``` bash
git clone https://github.com/joelong01/repo-config.git
cd repo-config
go build -o repo-config
```
Alternatively, you can install it directly using go install:

``` bash
go install github.com/joelong01/repo-config@latest
```
Make sure that your $GOPATH/bin directory is in your system's PATH so that you can run the 
repo-config command from anywhere.

The dockerfil (.devctainer/Dockerfile) will install zsh and other usefule tools.  It will also modify the .zshrc as follows

``` Docker
RUN echo 'if [[ -f "load_config.sh" ]]; then' >> ~/.zshrc && \
    echo '    source ./load_config.sh' >> ~/.zshrc && \
    echo 'fi' >> ~/.zshrc && \
    echo '' >> ~/.zshrc && \
    echo '# Function to source .ENV files from the corresponding PROJECT_DIR' >> ~/.zshrc && \
    echo 'source_env_files() {' >> ~/.zshrc && \
    echo '    PROJECT_DIR=$(basename "$PWD")' >> ~/.zshrc && \
    echo '    CONFIG_DIR="$HOME/.repo-config/$PROJECT_DIR"' >> ~/.zshrc && \
    echo '    if [[ -d "$CONFIG_DIR" ]]; then' >> ~/.zshrc && \
    echo '        for env_file in "$CONFIG_DIR"/.*.env; do' >> ~/.zshrc && \
    echo '            if [[ -f "$env_file" ]]; then' >> ~/.zshrc && \
    echo '                echo "sourcing $env_file"' >> ~/.zshrc && \
    echo '                source "$env_file"' >> ~/.zshrc && \
    echo '            fi' >> ~/.zshrc && \
    echo '        done' >> ~/.zshrc && \
    echo '    fi' >> ~/.zshrc && \
    echo '}' >> ~/.zshrc && \
    echo '' >> ~/.zshrc && \
    echo '# Call the function when the shell starts' >> ~/.zshrc && \
    echo 'source_env_files' >> ~/.zshrc

```
The ```load_config.sh``` shell script file does the following:

``` bash
if [[ ! -f "repo-config" ]]; then
    go build -o repo-config
fi

output=$(./repo-config collect --json ./end-to-end-test/test-config.json --silent)

result=$(echo $output | jq -r .status)

if [[ "$result" != "ok" ]]; then
    echo -n "ERROR: "
    echo $result | jq .message
fi

```
It checks to see if the ```repo_config``` executable is in the current directory, and if not, tries to build it (useful when writing this program!) Then it will call it to collect the needed settings.  In any other project, the first if clause can be deleted as it will never work. 

The function ```bash source_env_files ``` looks to see if there is any configration based on the current directory. If there is in loads the .ENV files from that directory into the environment.  This should be perserved if repo_config is used in any project.

# Usage
The repo-config tool provides two main commands:
``` 
collect: Collect repository configurations and generate output files.
delete: Delete generated output files.
Collect Command
Description
The collect command reads a JSON configuration file and generates corresponding output files containing the collected configuration values. It can run in interactive or silent mode.
```
Syntax
``` bash
repo-config collect --json <path_to_config.json> [--silent]
Options
--json, -j: (Required) Path to the JSON configuration file containing the configuration items.
--silent, -s: (Optional) Run the command in silent mode. In silent mode, the command operates without interactive prompts and uses default values or existing configuration where possible.
```
## Example

# Interactive mode
```
repo-config collect --json config.json
```
# Silent mode
``` 
repo-config collect --json config.json --silent
```

# Delete Command

The delete command deletes the output files generated by the collect command, such as the .env and .json files derived from the input configuration file. It can run in interactive or silent mode.

Syntax
``` 
repo-config delete --json <path_to_config.json> [--silent]
Options
--json, -j: (Required) Path to the JSON configuration file used to determine which output files to delete.
--silent, -s: (Optional) Run the command in silent mode. In silent mode, the command deletes the output files without prompting for confirmation.
``` 
# Example

## Interactive mode
```
repo-config delete --json config.json
```
# Silent mode
```
repo-config delete --json config.json --silent
```
Configuration JSON File Format
The input JSON configuration file should be structured as a JSON object where each key represents a configuration item. Each configuration item is an object with the following properties:
``` json
description: (Required) A description of the configuration item.
default: (Optional) The default value for the configuration item.
shellscript: (Optional) A shell script to execute for retrieving the value.
tempEnvironmentVariableName: (Optional) The name of a temporary environment variable to set.
requiredAsEnv: (Optional) A boolean indicating whether the configuration item is required as an environment variable.
```
Example config.json:
``` json
{
    "DATABASE_URL": {
        "description": "The URL of the database",
        "default": "",
        "tempEnvironmentVariableName": "DATABASE_URL",
        "requiredAsEnv": true
    },
    "API_KEY": {
        "description": "API key for external service",
        "default": "",
        "tempEnvironmentVariableName": "API_KEY",
        "requiredAsEnv": true
    }
}
```


License
This project is licensed under the MIT License - see the LICENSE file for details.

Contributing
Contributions are welcome! Please submit a pull request or open an issue to discuss improvements or new features.

Contact
For questions or support, please contact your.email@example.com or open an issue on the GitHub repository.