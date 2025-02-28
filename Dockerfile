FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o api-server

# Use a smaller image for runtime
FROM alpine:3.18

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/api-server .
# Copy necessary files
COPY api-keys.yaml .
COPY .env .

# Expose port
EXPOSE 8082

# Run the application
CMD ["./api-server"] 