# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copy shared-go package
COPY packages/shared-go/go.mod packages/shared-go/go.sum ./packages/shared-go/
COPY packages/shared-go/ ./packages/shared-go/

# Copy api-gateway service
COPY services/api-gateway/go.mod services/api-gateway/go.sum ./services/api-gateway/
COPY services/api-gateway/ ./services/api-gateway/

# Download dependencies
WORKDIR /app/services/api-gateway
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/server

# Run stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
