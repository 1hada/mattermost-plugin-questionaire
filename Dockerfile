# podman build -t plugin-builder .
# podman run --rm -v "$PWD":/app -v "$PWD/dist":/app/dist plugin-builder


FROM golang:1.22-bullseye

# Install Node.js 20 and npm
RUN apt-get update && \
    apt-get install -y curl gnupg && \
    curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y nodejs make git && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Default command to build
CMD ["make"]