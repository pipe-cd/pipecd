# Builder image to build go program.
FROM golang:1.23 AS builder

COPY protoc-gen-auth /protoc-gen-auth
RUN cd /protoc-gen-auth \
  && go build -o /usr/local/bin/protoc-gen-auth . \
  && chmod +x /usr/local/bin/protoc-gen-auth

# Codegen image which is actually being used.
FROM golang:1.23

# This is version of protobuf installed in the image.
# See https://pkgs.alpinelinux.org/packages?name=protobuf&branch=v3.16
# NOTE: Start from protobuf v3.20.1, the protoc-gen-js is not included in protobuf package.
ENV PROTOC_GEN_JS_VER=3.21.2
ENV PROTOC_GEN_GO_VER=1.27.1
ENV PROTOC_GEN_GRPC_WEB_VER=1.3.1
ENV PROTOC_GEN_GO_GRPC_VER=1.2.0
ENV PROTOC_GEN_VALIDATE_VER=0.6.6
ENV GOMOCK_VER=1.6.0

# dependecies and protoc
RUN apt update && apt install -y protobuf-compiler

# protoc-gen-go
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GEN_GO_VER}

# protoc-gen-go-grpc
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${PROTOC_GEN_GO_GRPC_VER}

# protoc-gen-grpc-web
RUN wget https://github.com/grpc/grpc-web/releases/download/${PROTOC_GEN_GRPC_WEB_VER}/protoc-gen-grpc-web-${PROTOC_GEN_GRPC_WEB_VER}-linux-x86_64 \
  && mv protoc-gen-grpc-web-${PROTOC_GEN_GRPC_WEB_VER}-linux-x86_64 /usr/local/bin/protoc-gen-grpc-web \
  && chmod +x /usr/local/bin/protoc-gen-grpc-web

# protoc-gen-validate
RUN go install github.com/envoyproxy/protoc-gen-validate@v${PROTOC_GEN_VALIDATE_VER} \
  && wget -q https://github.com/envoyproxy/protoc-gen-validate/archive/refs/tags/v${PROTOC_GEN_VALIDATE_VER}.tar.gz -O protoc-gen-validate.tar.gz \
  && mkdir -p /go/src/github.com/envoyproxy \
  && tar xvfz protoc-gen-validate.tar.gz -C /go/src/github.com/envoyproxy \
  && rm protoc-gen-validate.tar.gz \
  && mv /go/src/github.com/envoyproxy/protoc-gen-validate-${PROTOC_GEN_VALIDATE_VER} /go/src/github.com/envoyproxy/protoc-gen-validate

# protoc-gen-js
# This is a workaround to use it  https://github.com/protocolbuffers/protobuf-javascript/issues/127#issuecomment-1361047597
RUN for target in x86_64 aarch_64; do \
  mkdir /protoc-gen-js-${target} && cd /protoc-gen-js-${target} \
  && wget https://github.com/protocolbuffers/protobuf-javascript/releases/download/v${PROTOC_GEN_JS_VER}/protobuf-javascript-${PROTOC_GEN_JS_VER}-linux-${target}.tar.gz \
  && tar xvfz protobuf-javascript-${PROTOC_GEN_JS_VER}-linux-${target}.tar.gz \
  && chmod +x bin/protoc-gen-js \
  && rm -rf protobuf-javascript-${PROTOC_GEN_JS_VER}-linux-${target}.tar.gz; \
  done && \
  mv /protoc-gen-js-aarch_64/ /protoc-gen-js-aarch64/ 

# protoc-gen-auth
COPY --from=BUILDER /usr/local/bin/protoc-gen-auth /usr/local/bin/

# gomock
RUN go install github.com/golang/mock/mockgen@v${GOMOCK_VER}

VOLUME /repo
WORKDIR /repo
