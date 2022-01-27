FROM golang:1.17.3-alpine3.15 as builder
WORKDIR /tool
COPY main.go .
COPY go.mod .
COPY go.sum .
RUN go build -o asset-publisher main.go

FROM alpine:3.15
RUN apk --no-cache add ca-certificates git
COPY --from=builder /tool/asset-publisher ./
RUN chmod +x ./asset-publisher
