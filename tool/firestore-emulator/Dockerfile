FROM openjdk:8-jre-alpine3.9

ARG GOOGLE_CLOUD_SDK_VERSION=392.0.0

RUN apk add --no-cache \
        python3 \
        curl && \
    curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-$GOOGLE_CLOUD_SDK_VERSION-linux-x86.tar.gz && \
    tar -C /usr/local -xvzf google-cloud-sdk-$GOOGLE_CLOUD_SDK_VERSION-linux-x86.tar.gz && \
    rm google-cloud-sdk-$GOOGLE_CLOUD_SDK_VERSION-linux-x86.tar.gz && \
    /usr/local/google-cloud-sdk/install.sh --quiet

ENV PATH $PATH:/usr/local/google-cloud-sdk/bin

RUN gcloud components install cloud-firestore-emulator

EXPOSE 8080
ENTRYPOINT ["gcloud", "beta", "emulators", "firesotre", "start"]
CMD ["--host-port=0.0.0.0:8080"]
