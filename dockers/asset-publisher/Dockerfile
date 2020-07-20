FROM golang:1.14.6-alpine3.12 AS builder
WORKDIR /tool
COPY main.go .
COPY go.mod .
COPY go.sum .
RUN go build -o asset-publisher main.go

FROM alpine:3.12
RUN apk --no-cache add ca-certificates git
COPY --from=builder /tool/asset-publisher ./
RUN chmod +x ./asset-publisher
