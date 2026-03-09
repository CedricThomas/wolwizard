FROM golang:1.25-alpine AS builder

WORKDIR /build

# Install git for go modules
RUN apk add --no-cache git

# Copy go mod files first for better caching
COPY go.mod go.sum* ./
RUN go mod tidy || true

# Copy the rest of the source
COPY . .

# Build the raspberry-agent binary
RUN CGO_ENABLED=0 GOOS=linux go build -o raspberry-agent ./cmd/raspberry-agent

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /

COPY --from=builder /build/raspberry-agent .

CMD ["/raspberry-agent"]