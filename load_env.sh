#!/bin/bash
# Add this to your .zshrc

# Function to source .ENV files from the corresponding PROJECT_DIR
source_env_files() {
    # Get the current directory name (PROJECT_DIR)
    PROJECT_DIR=$(basename "$PWD")
    # Define the path in $HOME/.repo-config where the .ENV files should be located
    CONFIG_DIR="$HOME/.repo-config/$PROJECT_DIR"
    # Check if the directory exists
    if [[ -d "$CONFIG_DIR" ]]; then
        # Loop through all files ending in .env and source them
        for env_file in "$CONFIG_DIR"/.*.env; do
            # Check if any .ENV files are found
            if [[ -f "$env_file" ]]; then
                echo "sourcing $env_file"
                # Source the .ENV file
                # shellcheck source=/dev/null
                source "$env_file"
            fi
        done
    fi
}

# Call the function when the shell starts
source_env_files
