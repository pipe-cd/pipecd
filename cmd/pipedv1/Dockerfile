# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.26.2 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN make build/go MOD=pipedv1 BUILD_OS=${TARGETOS} BUILD_ARCH=${TARGETARCH}

# ghcr.io/pipe-cd/piped-base:v0.58.0-193-g4253e54
FROM ghcr.io/pipe-cd/piped-base@sha256:054891ce4d2216def0b5a4067cca264e0733582700be8b0e65cc5707e151c464

COPY --from=builder /app/.artifacts/pipedv1 /usr/local/bin/piped

ENTRYPOINT ["/sbin/tini", "--", "piped"]
