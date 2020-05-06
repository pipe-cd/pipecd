FROM alpine:3.10.1

ENV RUNNER_BIN_DIR=/usr/local/piped
ENV PATH="${RUNNER_BIN_DIR}:${PATH}"

ADD install-helm.sh /installer/install-helm.sh

RUN \
    apk add --no-cache \
        ca-certificates \
        git \
        openssh \
        curl \
        bash && \
    update-ca-certificates && \
    mkdir ${RUNNER_BIN_DIR} && \
    # Pre-install helm.
    /installer/install-helm.sh && \
    # Delete installer directory.
    rm -rf /installer && \
    apk del curl && \
    rm -f /var/cache/apk/*
