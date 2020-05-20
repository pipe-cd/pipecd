FROM alpine:3.10.1

ENV PIPED_BIN_DIR=/usr/local/piped
ENV PATH="${PIPED_BIN_DIR}:${PATH}"

ADD install-helm.sh /installer/install-helm.sh
ADD install-kubectl.sh /installer/install-kubectl.sh
ADD install-kustomize.sh /installer/install-kustomize.sh

RUN \
    apk add --no-cache \
        ca-certificates \
        git \
        openssh \
        curl \
        bash && \
    update-ca-certificates && \
    mkdir ${PIPED_BIN_DIR} && \
    # Pre-install helm.
    /installer/install-helm.sh && \
    # Pre-install kubectl.
    /installer/install-kubectl.sh && \
    # Pre-install kustomize.
    /installer/install-kustomize.sh && \
    # Delete installer directory.
    rm -rf /installer && \
    apk del curl && \
    rm -f /var/cache/apk/*
