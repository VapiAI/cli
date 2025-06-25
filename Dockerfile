# Build stage
FROM golang:1.24.4-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN make build

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 -S vapi && \
    adduser -u 1000 -S vapi -G vapi

# Copy binary from builder
COPY --from=builder /build/build/vapi /usr/local/bin/vapi

# Set ownership
RUN chown vapi:vapi /usr/local/bin/vapi

# Switch to non-root user
USER vapi

# Set entrypoint
ENTRYPOINT ["vapi"]

# Default command
CMD ["--help"] 