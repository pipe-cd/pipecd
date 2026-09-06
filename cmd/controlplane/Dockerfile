# syntax=docker/dockerfile:1
# web builder
# because this stage builds only web assets, we can use any platform
FROM --platform=$BUILDPLATFORM node:24.10.0-alpine AS web

WORKDIR /app

COPY . .

RUN apk add --no-cache make git

RUN make update/web-deps
RUN make build/web

# controlplane builder
FROM --platform=$BUILDPLATFORM golang:1.26.2 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN make build/go MOD=controlplane BUILD_OS=${TARGETOS} BUILD_ARCH=${TARGETARCH}

FROM alpine:3.22.2

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

COPY --from=builder /app/.artifacts/controlplane /usr/local/bin/controlplane
COPY --from=web /app/.artifacts/web-static web/static

ENTRYPOINT ["controlplane"]
