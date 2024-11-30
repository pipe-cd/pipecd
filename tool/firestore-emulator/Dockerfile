ARG GOOGLE_CLOUD_SDK_VERSION=392.0.0-alpine

FROM google/cloud-sdk:$GOOGLE_CLOUD_SDK_VERSION

RUN apk add --update --no-cache openjdk11-jre-headless

ENV FIRESTORE_PROJECT_ID="pipecd-test"
RUN gcloud config set project $FIRESTORE_PROJECT_ID && gcloud components install cloud-firestore-emulator beta --quiet

ENV PORT=8080
EXPOSE "$PORT"

ENTRYPOINT gcloud beta emulators firestore start --host-port="0.0.0.0:${PORT}"
