# Stage 1: Build static executable
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install ca-certificates and git if needed
RUN apk add --no-cache ca-certificates git

# Cache go modules
COPY go.mod go.sum* ./
RUN go mod download || true

# Copy entire source code
COPY . .

# Compile binary statically with stripped symbols
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/api ./cmd/api

# Stage 2: Minimal non-root runtime container
FROM alpine:3.20

# Create non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/api /app/api

USER appuser

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/app/api"]
