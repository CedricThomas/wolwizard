FROM golang:1.25-alpine AS builder

WORKDIR /build

# Install git for go modules
RUN apk add --no-cache git

# Copy go mod files first for better caching
COPY go.mod go.sum* ./
RUN go mod tidy || true

# Copy the rest of the source
COPY . .

# Build the pc-agent binary
RUN CGO_ENABLED=0 GOOS=linux go build -o pc-agent ./cmd/pc-agent

FROM nvidia/cuda:12.4.1-runtime-ubuntu22.04

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /

COPY --from=builder /build/pc-agent .

CMD ["/pc-agent"]