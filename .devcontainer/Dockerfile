# Use an official Go runtime as a parent image
FROM mcr.microsoft.com/devcontainers/go:1.25

# Install Node.js, npm, and Vite for frontend development
USER root
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && curl -fsSL https://deb.nodesource.com/setup_22.x | bash - \
    && apt-get install -y nodejs \
    && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/*

# Install global node packages
RUN npm install -g npm@latest vite concurrently @anthropic-ai/claude-code

# Install Go tools
RUN go install github.com/air-verse/air@v1.63.0 \
    && go install github.com/swaggo/swag/cmd/swag@v1.16.6 # For live reloading Go apps
RUN chown -R vscode /go
RUN chmod -R u+w /go

# Set back to non-root user
USER vscode

WORKDIR /workspaces
