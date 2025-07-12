# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.24.1 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN make build/go MOD=piped BUILD_OS=${TARGETOS} BUILD_ARCH=${TARGETARCH}

# https://github.com/pipe-cd/pipecd/pkgs/container/piped-base/459025001?tag=v0.52.0-97-g3e3df6d
FROM ghcr.io/pipe-cd/piped-base@sha256:fc2bde38211a80ee0df6d90ed168c817f30a7d0dbe8a0aff872a93ffb7d40a66

COPY --from=builder /app/.artifacts/piped /usr/local/bin/piped

ENTRYPOINT ["/sbin/tini", "--", "piped"]
