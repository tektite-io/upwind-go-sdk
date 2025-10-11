# Multi-stage build for minimal image size

# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make ca-certificates tzdata

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE} -w -s" \
    -a -installsuffix cgo \
    -o upwind \
    ./cmd/upwind

# Final stage - minimal runtime image
FROM scratch

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /build/upwind /upwind

# Set entrypoint
ENTRYPOINT ["/upwind"]

# Default command shows help
CMD ["help"]

# Labels
LABEL org.opencontainers.image.title="Upwind Go SDK CLI"
LABEL org.opencontainers.image.description="Command-line interface for the Upwind Security API"
LABEL org.opencontainers.image.vendor="Tektite IO"
LABEL org.opencontainers.image.source="https://github.com/tektite-io/upwind-go-sdk"

