FROM alpine:3.12.1

ENV PIPED_BIN_DIR=/home/pipecd/tools
ENV PATH="${PIPED_BIN_DIR}:${PATH}"

COPY install-helm.sh /installer/install-helm.sh
COPY install-kubectl.sh /installer/install-kubectl.sh
COPY install-kustomize.sh /installer/install-kustomize.sh
COPY install-terraform.sh /installer/install-terraform.sh

RUN \
    addgroup -S -g 1000 pipecd && \
    adduser -S -u 1000 -G pipecd -h /home/pipecd pipecd && \
    apk add --no-cache \
        ca-certificates \
        git \
        openssh \
        curl \
        bash && \
    update-ca-certificates && \
    mkdir ${PIPED_BIN_DIR} && \
    chown pipecd ${PIPED_BIN_DIR} && \
    # Pre-install the default version of helm.
    /installer/install-helm.sh && \
    # Pre-install the default version of kubectl.
    /installer/install-kubectl.sh && \
    # Pre-install the default version of kustomize.
    /installer/install-kustomize.sh && \
    # Pre-install the default version of terraform.
    /installer/install-terraform.sh && \
    # Delete installer directory.
    rm -rf /installer && \
    rm -f /var/cache/apk/*

USER pipecd
