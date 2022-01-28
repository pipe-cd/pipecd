FROM golang:1.17.3-alpine3.15 AS builder
WORKDIR /tool
COPY main.go go.mod go.sum ./
RUN go build -o asset-publisher main.go

FROM alpine:3.15
RUN apk --no-cache add ca-certificates git
COPY --from=builder /tool/asset-publisher ./
RUN chmod +x ./asset-publisher
