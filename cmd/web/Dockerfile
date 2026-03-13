# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files first for better caching
COPY go.mod go.sum* ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the web binary
RUN CGO_ENABLED=0 GOOS=linux go build -o web ./cmd/web

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /

COPY --from=builder /build/web .
COPY --from=builder /build/static ./static

EXPOSE 8080

CMD ["/web"]