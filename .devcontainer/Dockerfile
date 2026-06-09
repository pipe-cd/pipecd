FROM mcr.microsoft.com/devcontainers/go:1.22-bookworm

# Install Node.js (for the web UI)
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y nodejs

# Install Kubernetes tools: kubectl, Helm, and Kind (for local clusters)
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && \
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash && \
    curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.22.0/kind-linux-amd64 && \
    chmod +x ./kind && mv ./kind /usr/local/bin/kind

# Install PipeCD specific dependencies if any (e.g., golangci-lint)
RUN curl -sSfL https://raw.githubusercontent.com/golangci-lint/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.56.2
