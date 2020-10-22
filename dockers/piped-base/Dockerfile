FROM alpine:3.10.1

ENV PIPED_BIN_DIR=/usr/local/piped
ENV PATH="${PIPED_BIN_DIR}:${PATH}"

ADD install-helm.sh /installer/install-helm.sh
ADD install-kubectl.sh /installer/install-kubectl.sh
ADD install-kustomize.sh /installer/install-kustomize.sh
ADD install-terraform.sh /installer/install-terraform.sh

RUN \
    apk add --no-cache \
        ca-certificates \
        git \
        openssh \
        curl \
        bash && \
    update-ca-certificates && \
    mkdir ${PIPED_BIN_DIR} && \
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
    rm -f /var/cache/apk/* && \
    addgroup -S -g 1000 pipecd && \
    adduser -S -u 1000 -G pipecd -h /home/pipecd pipecd

USER pipecd
