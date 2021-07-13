FROM golang:1.16.5-alpine3.14 as builder
COPY . /app
RUN cd /app && go build -o /plan-preview .

FROM gcr.io/pipecd/pipectl:v0.11.1
COPY --from=builder /plan-preview /
ENV PATH $PATH:/app/cmd/pipectl
RUN chmod +x /plan-preview
ENTRYPOINT ["/plan-preview"]
