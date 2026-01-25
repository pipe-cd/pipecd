# Stage 1: Builder - Build Go-based plugins
FROM golang:1.25.2 AS builder

# Version configuration
ARG PROTOC_GEN_GO_VER=1.27.1
ARG PROTOC_GEN_GO_GRPC_VER=1.2.0
ARG PROTOC_GEN_VALIDATE_VER=0.6.6
ARG GOMOCK_VER=0.6.0

# Build protoc-gen-auth from source
COPY protoc-gen-auth /protoc-gen-auth
RUN cd /protoc-gen-auth \
  && go build -o /out/protoc-gen-auth . \
  && chmod +x /out/protoc-gen-auth

# Install Go-based protoc plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GEN_GO_VER} \
  && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${PROTOC_GEN_GO_GRPC_VER} \
  && go install github.com/envoyproxy/protoc-gen-validate@v${PROTOC_GEN_VALIDATE_VER} \
  && go install go.uber.org/mock/mockgen@v${GOMOCK_VER}

# Copy built binaries to output directory
RUN mkdir -p /out \
  && cp /go/bin/protoc-gen-go /out/ \
  && cp /go/bin/protoc-gen-go-grpc /out/ \
  && cp /go/bin/protoc-gen-validate /out/ \
  && cp /go/bin/mockgen /out/

# Download protoc-gen-validate proto files (needed for validate imports)
RUN wget -q https://github.com/envoyproxy/protoc-gen-validate/archive/refs/tags/v${PROTOC_GEN_VALIDATE_VER}.tar.gz -O /tmp/validate.tar.gz \
  && mkdir -p /out/validate-protos \
  && tar xzf /tmp/validate.tar.gz -C /out/validate-protos --strip-components=1 \
  && rm /tmp/validate.tar.gz

# Stage 2: Downloader - Download pre-built binaries
FROM debian:bookworm-slim AS downloader

RUN apt-get update && apt-get install -y --no-install-recommends wget ca-certificates \
  && rm -rf /var/lib/apt/lists/*

# Download protoc-gen-grpc-web
ARG PROTOC_GEN_GRPC_WEB_VER=1.3.1
RUN wget -q https://github.com/grpc/grpc-web/releases/download/${PROTOC_GEN_GRPC_WEB_VER}/protoc-gen-grpc-web-${PROTOC_GEN_GRPC_WEB_VER}-linux-x86_64 \
  -O /protoc-gen-grpc-web \
  && chmod +x /protoc-gen-grpc-web

# Download protoc-gen-js for both architectures
# This is a workaround: https://github.com/protocolbuffers/protobuf-javascript/issues/127#issuecomment-1361047597
ARG PROTOC_GEN_JS_VER=3.21.2
RUN for target in x86_64 aarch_64; do \
    mkdir -p /protoc-gen-js-${target} \
    && wget -q https://github.com/protocolbuffers/protobuf-javascript/releases/download/v${PROTOC_GEN_JS_VER}/protobuf-javascript-${PROTOC_GEN_JS_VER}-linux-${target}.tar.gz \
    && tar xzf protobuf-javascript-${PROTOC_GEN_JS_VER}-linux-${target}.tar.gz -C /protoc-gen-js-${target} \
    && chmod +x /protoc-gen-js-${target}/bin/protoc-gen-js \
    && rm protobuf-javascript-${PROTOC_GEN_JS_VER}-linux-${target}.tar.gz; \
  done \
  && mv /protoc-gen-js-aarch_64 /protoc-gen-js-aarch64

