# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.25.0 AS builder

ARG TARGETOS
ARG TARGETARCH
ARG BUILD_VERSION

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN make build/go MOD=launcher BUILD_OS=${TARGETOS} BUILD_ARCH=${TARGETARCH} BUILD_VERSION=${BUILD_VERSION}

# https://github.com/pipe-cd/pipecd/pkgs/container/piped-base/459025001?tag=v0.52.0-97-g3e3df6d
FROM ghcr.io/pipe-cd/piped-base@sha256:fc2bde38211a80ee0df6d90ed168c817f30a7d0dbe8a0aff872a93ffb7d40a66

COPY --from=builder /app/.artifacts/launcher /usr/local/bin/launcher

ENTRYPOINT ["/sbin/tini", "--", "launcher"]
