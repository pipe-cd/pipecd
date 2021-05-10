FROM alpine:3.12.1

ARG PIPED_USER=piped
ARG PIPED_USER_GROUP=piped
ARG PIPED_UID=1000
ARG PIPED_GID=1000

ENV HOME=/home/${PIPED_USER}
ENV PIPED_TOOLS_DIR="${HOME}/.piped/tools"
ENV PATH="${PIPED_TOOLS_DIR}:${PATH}"

COPY install-helm.sh /installer/install-helm.sh
COPY install-kubectl.sh /installer/install-kubectl.sh
COPY install-kustomize.sh /installer/install-kustomize.sh
COPY install-terraform.sh /installer/install-terraform.sh
COPY install-nss-wrapper.sh /installer/install-nss-wrapper.sh

RUN \
    addgroup -S -g $PIPED_GID $PIPED_USER_GROUP && \
    adduser -S -u $PIPED_UID -G $PIPED_USER_GROUP -h $HOME $PIPED_USER && \
    apk add --no-cache \
        ca-certificates \
        git \
        openssh \
        curl \
        bash \
        gcc \
        #musl-dev \
        libc-dev \
        make \
        cmake && \

    wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.33-r0/glibc-2.33-r0.apk && \
    apk add glibc-2.33-r0.apk && \

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
    # Install nss_wrapper to add an random UID to "passwd" without having to directly modify /etc/passwd.
    # This is for those whose Piped's user isn't in /etc/passwd
    /installer/install-nss-wrapper.sh && \
    # Delete installer directory.
    rm -rf /installer && \
    rm -f /var/cache/apk/* && \
    chown -R $PIPED_USER:$PIPED_USER_GROUP $HOME && \
    chmod 770 -R $HOME

USER $PIPED_USER
