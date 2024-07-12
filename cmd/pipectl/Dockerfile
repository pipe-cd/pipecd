FROM alpine:3.20

RUN apk add --no-cache git

ADD .artifacts/pipectl /usr/local/bin/pipectl

ENTRYPOINT ["pipectl"]
