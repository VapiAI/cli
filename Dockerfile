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

# Build the binary and generate man pages
RUN make build && make man-pages

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS and mandoc for manual pages
RUN apk add --no-cache ca-certificates mandoc

# Create non-root user
RUN addgroup -g 1000 -S vapi && \
    adduser -u 1000 -S vapi -G vapi

# Create man page directory
RUN mkdir -p /usr/local/share/man/man1

# Copy binary from builder
COPY --from=builder /build/build/vapi /usr/local/bin/vapi

# Copy man pages from builder
COPY --from=builder /build/man/*.1 /usr/local/share/man/man1/

# Set ownership (mandoc doesn't need database updates)
RUN chown vapi:vapi /usr/local/bin/vapi && \
    chown -R vapi:vapi /usr/local/share/man/man1/

# Switch to non-root user
USER vapi

# Set environment for man pages
ENV MANPATH="/usr/local/share/man:/usr/share/man"

# Set entrypoint
ENTRYPOINT ["vapi"]

# Default command
CMD ["--help"] 