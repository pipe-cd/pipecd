FROM golang:1.13.4-alpine3.10 AS builder
COPY main.go go.mod go.sum .
RUN go build -o /chart-releaser main.go

FROM alpine:3.10
RUN apk --no-cache add ca-certificates
COPY --from=builder /chart-releaser ./
