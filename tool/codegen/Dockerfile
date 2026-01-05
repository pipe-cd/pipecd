# Builder stage
FROM golang:1.25.5-bookworm AS builder

COPY protoc-gen-auth /protoc-gen-auth
RUN cd /protoc-gen-auth \
  && go build -o /usr/local/bin/protoc-gen-auth . \
  && chmod +x /usr/local/bin/protoc-gen-auth

ENV PROTOC_GEN_GO_VER=1.27.1
ENV PROTOC_GEN_GO_GRPC_VER=1.2.0
ENV PROTOC_GEN_VALIDATE_VER=0.6.6
ENV GOMOCK_VER=0.6.0

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GEN_GO_VER}
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${PROTOC_GEN_GO_GRPC_VER}
RUN go install github.com/envoyproxy/protoc-gen-validate@v${PROTOC_GEN_VALIDATE_VER}
RUN go install go.uber.org/mock/mockgen@v${GOMOCK_VER}

RUN wget -q https://github.com/envoyproxy/protoc-gen-validate/archive/refs/tags/v${PROTOC_GEN_VALIDATE_VER}.tar.gz -O protoc-gen-validate.tar.gz \
  && mkdir -p /proto-validate \
  && tar xzf protoc-gen-validate.tar.gz -C /proto-validate --strip-components=1 \
  && rm protoc-gen-validate.tar.gz

# Runtime stage
FROM debian:bookworm-slim

ENV PROTOC_GEN_JS_VER=3.21.2
ENV PROTOC_GEN_GRPC_WEB_VER=1.3.1

RUN apt-get update && apt-get install -y --no-install-recommends \
    protobuf-compiler \
    wget \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/bin/protoc-gen-go /usr/local/bin/
COPY --from=builder /go/bin/protoc-gen-go-grpc /usr/local/bin/
COPY --from=builder /go/bin/protoc-gen-validate /usr/local/bin/
COPY --from=builder /go/bin/mockgen /usr/local/bin/
COPY --from=builder /usr/local/bin/protoc-gen-auth /usr/local/bin/
COPY --from=builder /proto-validate /go/src/github.com/envoyproxy/protoc-gen-validate

RUN wget -q https://github.com/grpc/grpc-web/releases/download/${PROTOC_GEN_GRPC_WEB_VER}/protoc-gen-grpc-web-${PROTOC_GEN_GRPC_WEB_VER}-linux-x86_64 \
  && mv protoc-gen-grpc-web-${PROTOC_GEN_GRPC_WEB_VER}-linux-x86_64 /usr/local/bin/protoc-gen-grpc-web \
  && chmod +x /usr/local/bin/protoc-gen-grpc-web

# See: https://github.com/protocolbuffers/protobuf-javascript/issues/127#issuecomment-1361047597
RUN for target in x86_64 aarch_64; do \
  mkdir /protoc-gen-js-${target} && cd /protoc-gen-js-${target} \
  && wget -q https://github.com/protocolbuffers/protobuf-javascript/releases/download/v${PROTOC_GEN_JS_VER}/protobuf-javascript-${PROTOC_GEN_JS_VER}-linux-${target}.tar.gz \
  && tar xzf protobuf-javascript-${PROTOC_GEN_JS_VER}-linux-${target}.tar.gz \
  && chmod +x bin/protoc-gen-js \
  && rm -f protobuf-javascript-${PROTOC_GEN_JS_VER}-linux-${target}.tar.gz; \
  done && \
  mv /protoc-gen-js-aarch_64/ /protoc-gen-js-aarch64/

RUN apt-get purge -y wget ca-certificates && apt-get autoremove -y && rm -rf /var/lib/apt/lists/*

VOLUME /repo
WORKDIR /repo
