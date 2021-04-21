FROM alpine:3.12.1

ENV PIPED_DIR=/etc/piped
ENV PIPED_TOOLS_DIR="${PIPED_DIR}/tools"
ENV PATH="${PIPED_TOOLS_DIR}:${PATH}"
ENV HELM_CACHE_HOME="${PIPED_DIR}/.cache/helm"
ENV HELM_CONFIG_HOME="${PIPED_DIR}/.config/helm"
ENV HELM_DATA_HOME="${PIPED_DIR}/.local/share/helm"
ENV KUBECONFIG="${PIPED_DIR}/.kube/config"

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
    mkdir -p ${PIPED_TOOLS_DIR} && \
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
    chown -hR pipecd:pipecd ${PIPED_DIR} && \
    chmod 770 -R ${PIPED_DIR}

USER pipecd
