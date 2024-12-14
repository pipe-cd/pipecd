# web builder
FROM node:20.18.1-alpine3.20 AS web

WORKDIR /app

COPY . .

RUN apk add --no-cache make
RUN make update/web-deps
RUN make build/web

# pipecd builder
FROM golang:1.23.3 AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN make build/go MOD=pipecd

FROM alpine:3.20

ARG GOOGLE_CLOUD_SDK_VERSION=477.0.0

RUN \
    apk add --no-cache \
        python3 \
        curl && \
    curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${GOOGLE_CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    tar -zxvf ./google-cloud-sdk-${GOOGLE_CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    rm ./google-cloud-sdk-${GOOGLE_CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    ./google-cloud-sdk/install.sh --quiet

ENV PATH="/google-cloud-sdk/bin:${PATH}"

RUN mkdir -p web/static

COPY --from=builder /app/.artifacts/pipecd /usr/local/bin/pipecd
COPY --from=web /app/.artifacts/web-static web/static

ENTRYPOINT ["pipecd"]
