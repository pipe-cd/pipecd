# pipecd/pipecd-base:0.1.0
FROM gcr.io/pipecd/pipecd-base@sha256:f3e98a27b85b8ead610c4f93cec8d936c760a43866cf817d32563daf9b198358

RUN mkdir -p web/static

ADD .artifacts/pipecd pipecd
ADD .artifacts/web-static web/static

ENTRYPOINT ["./pipecd"]
