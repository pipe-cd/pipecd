FROM golang:1.22.5-alpine3.20

RUN apk update && apk add git

COPY . /app

RUN cd /app && \
  go build -o /gh-release . && \
  chmod +x /gh-release

ENTRYPOINT ["/gh-release"]
