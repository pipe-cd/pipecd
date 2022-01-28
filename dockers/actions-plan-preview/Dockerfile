FROM golang:1.17.3-alpine3.15 as builder
COPY . /app
RUN cd /app && go build -o /plan-preview .

FROM gcr.io/pipecd/pipectl:v0.24.0
COPY --from=builder /plan-preview /
ENV PATH $PATH:/app/cmd/pipectl
RUN chmod +x /plan-preview
ENTRYPOINT ["/plan-preview"]
