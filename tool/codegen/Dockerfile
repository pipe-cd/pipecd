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

