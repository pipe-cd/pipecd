FROM golang:1.13.4-alpine3.10 AS builder
COPY main.go /pipecd/main.go
ADD install-helm.sh /installer/install-helm.sh
RUN apk --no-cache add curl bash
RUN /installer/install-helm.sh
RUN cd /pipecd && \
    go mod init pipecd.dev/chart-releaser && \
    go mod tidy && \
    go build -o chart-releaser main.go

FROM alpine:3.10
COPY --from=builder /pipecd/chart-releaser ./
COPY --from=builder /helm /usr/local/bin
RUN apk --no-cache add ca-certificates && \
    chmod +x ./chart-releaser
