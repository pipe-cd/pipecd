FROM golang:1.22.4-alpine3.20 AS builder
COPY main.go .
RUN go build -o /server main.go

FROM alpine:3
RUN apk --no-cache add ca-certificates

COPY --from=builder /server ./
RUN chmod +x ./server

COPY public /public

EXPOSE 8080
ENTRYPOINT ["./server"]
