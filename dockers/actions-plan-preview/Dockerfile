FROM golang:1.16.5-alpine3.14 as builder
COPY . /app
RUN cd /app && go build -o /plan-preview .

# TODO: Update pipectl version to v0.13.2 once it's rolled out.
FROM gcr.io/pipecd/pipectl:v0.13.1-4-gec4b78b
COPY --from=builder /plan-preview /
ENV PATH $PATH:/app/cmd/pipectl
RUN chmod +x /plan-preview
ENTRYPOINT ["/plan-preview"]
