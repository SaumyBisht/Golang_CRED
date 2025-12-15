# Build Stage
FROM golang:1.25.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Builds the binary named 'core-app'
RUN go build -o core-app .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/core-app .
EXPOSE 8081
CMD ["./core-app"]