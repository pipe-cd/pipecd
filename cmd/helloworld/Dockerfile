# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.27.1 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN make build/go MOD=helloworld BUILD_OS=${TARGETOS} BUILD_ARCH=${TARGETARCH}

FROM alpine:3.24.2

COPY --from=builder /app/.artifacts/helloworld /usr/local/bin/helloworld

ENTRYPOINT ["helloworld"]
