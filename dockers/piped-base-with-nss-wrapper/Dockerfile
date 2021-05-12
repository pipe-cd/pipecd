FROM gcr.io/pipecd/piped-base:0.2.2

ARG PIPED_UID=1000
ARG PIPED_GID=0

USER root

COPY install-nss-wrapper.sh /installer/install-nss-wrapper.sh

RUN \
    apk add --no-cache \
        gcc \
        libc-dev \
        musl-nscd-dev \
        make \
        cmake && \

    # Install glibc to be used for building nss_wrapper.
    wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.33-r0/glibc-2.33-r0.apk && \
    apk add glibc-2.33-r0.apk && \

    # Install nss_wrapper to add an random UID to "passwd" at runtime without having to directly modify /etc/passwd.
    # This is for those whose Piped's user isn't in /etc/passwd
    /installer/install-nss-wrapper.sh && \
    # Delete installer directory.
    rm -rf /installer && \
    rm -f /var/cache/apk/* && \
    grep -v -e ^default /etc/passwd > $HOME/passwd && \
    chown -R $PIPED_UID:$PIPED_GID $HOME && \
    chmod 770 -R $HOME

ENV LD_PRELOAD=/usr/local/lib64/libnss_wrapper.so
ENV NSS_WRAPPER_PASSWD=${HOME}/passwd
ENV NSS_WRAPPER_GROUP=/etc/group

USER $PIPED_USER
