FROM golang:1.24.7-alpine3.22 AS builder
WORKDIR /app
COPY ../go.mod ../go.sum ./
RUN go mod download
COPY . .
RUN go build -o app ./cmd/api/main.go

FROM alpine:latest
WORKDIR /app
# COPY --chown=65534:65534 --from=builder /app/pkg/mail/templates /app/pkg/mail/templates
COPY --chown=65534:65534 --from=builder /app/app ./app
USER 65534

ENTRYPOINT ["./app"]