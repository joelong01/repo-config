# Use the official Go image as a base
FROM golang:1.21-bullseye

# Install required packages
RUN apt-get update && apt-get install -y \
    git \
    curl \
    ca-certificates \
    gnupg \
    lsb-release \
    sudo \
    zsh \
    jq

# Set the Go PATH environment variable
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH


# Create a non-root user for VS Code to use
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && echo "$USERNAME ALL=(ALL) NOPASSWD:ALL" | tee /etc/sudoers.d/$USERNAME

# Set the working directory
WORKDIR /workspace

# Switch to the non-root user
USER $USERNAME

# Install Cobra CLI
RUN go install github.com/spf13/cobra-cli@latest

# Install GitHub CLI using the standard installation script
RUN curl -sS https://webi.sh/gh | bash
ENV PATH="/home/$USERNAME/.local/bin:${PATH}"
RUN sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" || true



# Expose ports if necessary
EXPOSE 8080

# update the .zshrc to load the configured environment
# Append the environment sourcing function to the .zshrc file
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

