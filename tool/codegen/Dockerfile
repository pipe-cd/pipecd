# Builder image to build go program.
FROM golang:1.17-alpine3.15 as BUILDER

COPY protoc-gen-auth /protoc-gen-auth
RUN cd /protoc-gen-auth \
  && go build -o /usr/local/bin/protoc-gen-auth . \
  && chmod +x /usr/local/bin/protoc-gen-auth

# Codegen image which is actually being used.
FROM golang:1.17-alpine3.15

ENV PROTOC_VER=3.19.4
ENV PROTOC_GEN_GO_VER=1.27.1
ENV PROTOC_GEN_GRPC_WEB_VER=1.3.1
ENV PROTOC_GEN_GO_GRPC_VER=1.2.0
ENV PROTOC_GEN_VALIDATE_VER=0.6.6
ENV GOMOCK_VER=1.6.0
ENV GLIBC_VERSION=2.33-r0

RUN apk --no-cache add wget bash \
    && wget -q https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub -O /etc/apk/keys/sgerrand.rsa.pub \
    && wget -q https://github.com/sgerrand/alpine-pkg-glibc/releases/download/${GLIBC_VERSION}/glibc-${GLIBC_VERSION}.apk -O glibc.apk \
    && apk add glibc.apk \
    && rm /etc/apk/keys/sgerrand.rsa.pub glibc.apk

# protoc
RUN wget -q https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VER}/protoc-${PROTOC_VER}-linux-x86_64.zip -O protoc.zip \
    && unzip protoc.zip -d /usr/local \
    && rm protoc.zip

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

# protoc-gen-auth
COPY --from=BUILDER /usr/local/bin/protoc-gen-auth /usr/local/bin/

# gomock
RUN go install github.com/golang/mock/mockgen@v${GOMOCK_VER}

VOLUME /repo
WORKDIR /repo

COPY ./codegen.sh /
RUN chmod +x /codegen.sh

ENTRYPOINT ["/codegen.sh"]
