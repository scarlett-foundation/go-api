FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger docs
ARG API_HOST=localhost:8082
ENV API_HOST=${API_HOST}
RUN /go/bin/swag init -g cmd/main.go -d ./ -o ./docs/swagger

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o api-server ./cmd

# Use a smaller image for runtime
FROM alpine:3.18

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/api-server .
# Copy necessary files
COPY api-keys.yaml .
# Copy Swagger docs
COPY --from=builder /app/docs/swagger ./docs/swagger

# Set permissions for app and docs
RUN chmod +x /app/api-server && \
    chmod -R 755 /app/docs

# Expose port
EXPOSE 8082

# Run the application
CMD ["./api-server"] 