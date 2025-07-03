# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.24.1 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN make build/go MOD=piped BUILD_OS=${TARGETOS} BUILD_ARCH=${TARGETARCH}

# https://github.com/pipe-cd/pipecd/pkgs/container/piped-base/452707327?tag=v0.52.0-76-g8a7294e
FROM ghcr.io/pipe-cd/piped-base@sha256:a8f3aba027d8c14aab1b093c9ccb26d40c3b551ba9160902445d186d555a92c7

COPY --from=builder /app/.artifacts/piped /usr/local/bin/piped

ENTRYPOINT ["/sbin/tini", "--", "piped"]
