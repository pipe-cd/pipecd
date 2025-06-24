FROM golang:1.24.1-alpine3.21

RUN apk update && apk add git

COPY . /app

RUN cd /app && \
  go build -o /gh-release . && \
  chmod +x /gh-release

ENTRYPOINT ["/gh-release"]
