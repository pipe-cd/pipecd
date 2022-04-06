FROM gcr.io/pipecd/runner:1.0.0

RUN \
    # Install Google Cloud SDK according to the instruction:
    # https://cloud.google.com/sdk/docs/downloads-apt-get
    echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | \
    tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
    apt-get -y install gnupg && \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add - && \
    apt-get update && \
    apt-get -y install \
      google-cloud-sdk=311.0.0-0 \
      google-cloud-sdk-firestore-emulator
