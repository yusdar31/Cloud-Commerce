# Build stage - context should be monorepo root
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copy shared-go package
COPY packages/shared-go ./packages/shared-go

# Copy service files
COPY services/user-service ./services/user-service

WORKDIR /app/services/user-service

# Download dependencies
RUN go mod download

# Build from cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/server

# Run stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /main .
EXPOSE 8081
CMD ["./main"]
