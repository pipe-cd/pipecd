# syntax=docker/dockerfile:1
# web builder
# because this stage builds only web assets, we can use any platform
FROM --platform=$BUILDPLATFORM node:20.19.0-alpine3.21 AS web

WORKDIR /app

COPY . .

RUN apk add --no-cache make git

RUN make update/web-deps
RUN make build/web

# pipecd builder
FROM --platform=$BUILDPLATFORM golang:1.24.1 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN make build/go MOD=control-plane BUILD_OS=${TARGETOS} BUILD_ARCH=${TARGETARCH}

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

COPY --from=builder /app/.artifacts/control-plane /usr/local/bin/control-plane
COPY --from=web /app/.artifacts/web-static web/static

ENTRYPOINT ["control-plane"]
