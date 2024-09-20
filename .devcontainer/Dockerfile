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
    zsh

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
