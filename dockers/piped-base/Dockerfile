FROM alpine:3.12.1

ENV PIPED_TOOLS_DIR=/tools
ENV PATH="${PIPED_TOOLS_DIR}:${PATH}"

COPY install-helm.sh /installer/install-helm.sh
COPY install-kubectl.sh /installer/install-kubectl.sh
COPY install-kustomize.sh /installer/install-kustomize.sh
COPY install-terraform.sh /installer/install-terraform.sh

RUN \
    addgroup -S -g 1000 pipecd && \
    adduser -S -u 1000 -G pipecd pipecd && \
    apk add --no-cache \
        ca-certificates \
        git \
        openssh \
        curl \
        bash && \
    update-ca-certificates && \
    mkdir ${PIPED_TOOLS_DIR} && \
    chmod 777 /etc/ssh/ssh_config && \
    # Pre-install the default version of helm.
    /installer/install-helm.sh && \
    # Pre-install the default version of kubectl.
    /installer/install-kubectl.sh && \
    # Pre-install the default version of kustomize.
    /installer/install-kustomize.sh && \
    # Pre-install the default version of terraform.
    /installer/install-terraform.sh && \
    chmod 775 ${PIPED_TOOLS_DIR}/* && \
    # Delete installer directory.
    rm -rf /installer && \
    rm -f /var/cache/apk/*

USER pipecd
