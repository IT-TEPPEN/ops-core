# Use an official Go runtime as a parent image
FROM mcr.microsoft.com/devcontainers/go:1.23

# Install Node.js, npm, and Vite for frontend development
USER root
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && curl -fsSL https://deb.nodesource.com/setup_22.x | bash - \
    && apt-get install -y nodejs \
    && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/*

# Install global node packages
RUN npm install -g npm@latest vite concurrently

# Install Go tools
RUN go install github.com/air-verse/air@latest # For live reloading Go apps

# Set back to non-root user
USER vscode

WORKDIR /workspaces