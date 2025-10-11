# syntax=docker/dockerfile:1
# web builder
# because this stage builds only web assets, we can use any platform
FROM --platform=$BUILDPLATFORM node:20.19.0-alpine3.21 AS web

ARG BUILD_VERSION

WORKDIR /app

RUN apk add --no-cache make git

COPY web/package.json web/yarn.lock ./web/
RUN yarn --cwd web install --prefer-offline

COPY .git ./.git
COPY web ./web
COPY Makefile .
RUN make build/web PIPECD_VERSION=${BUILD_VERSION}

# pipecd builder
FROM --platform=$BUILDPLATFORM golang:1.25.0 AS builder

ARG TARGETOS
ARG TARGETARCH
ARG BUILD_VERSION

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY .git ./.git
COPY pkg/ ./pkg/
COPY cmd/ ./cmd/
COPY Makefile .
RUN make build/go MOD=pipecd BUILD_OS=${TARGETOS} BUILD_ARCH=${TARGETARCH} BUILD_VERSION=${BUILD_VERSION}

FROM alpine:3.21.3

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
