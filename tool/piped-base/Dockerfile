FROM alpine:3.12.1

ARG PIPED_USER=piped
ARG PIPED_USER_GROUP=piped
ARG PIPED_UID=1000
ARG PIPED_GID=1000

ENV HOME=/home/${PIPED_USER}
ENV PIPED_TOOLS_DIR="${HOME}/.piped/tools"
ENV PATH="${PIPED_TOOLS_DIR}:${PATH}"
ENV HELM_CONFIG_HOME="${PIPED_TOOLS_DIR}"
ENV HELM_REGISTRY_CONFIG="${PIPED_TOOLS_DIR}/helm-config.json"
ENV HELM_REPOSITORY_CONFIG="${PIPED_TOOLS_DIR}/helm-repositories.yaml"

COPY install-helm.sh /installer/install-helm.sh
COPY install-kubectl.sh /installer/install-kubectl.sh
COPY install-kustomize.sh /installer/install-kustomize.sh
COPY install-terraform.sh /installer/install-terraform.sh

RUN \
    addgroup -S -g $PIPED_GID $PIPED_USER_GROUP && \
    adduser -S -u $PIPED_UID -G $PIPED_USER_GROUP -h $HOME $PIPED_USER && \
    apk add --no-cache \
        ca-certificates \
        git \
        openssh \
        curl \
        bash && \
    update-ca-certificates && \
    mkdir -p ${PIPED_TOOLS_DIR} && \
    touch ${HELM_REGISTRY_CONFIG} && \
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
    chown -R $PIPED_USER:$PIPED_USER_GROUP $HOME && \
    chmod 770 -R $HOME

USER $PIPED_USER
