FROM golang:1.23.3 AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN make build/go MOD=helloworld

FROM alpine:3.21.0

COPY --from=builder /app/.artifacts/helloworld /usr/local/bin/helloworld

ENTRYPOINT ["helloworld"]
