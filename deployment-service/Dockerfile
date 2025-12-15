# Build Stage
FROM golang:1.25.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Builds the binary named 'deploy-app'
RUN go build -o deploy-app

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/deploy-app .
EXPOSE 8082
CMD ["./deploy-app"]