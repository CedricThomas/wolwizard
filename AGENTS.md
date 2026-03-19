# Console Development Guide

## Build & Test Commands

### Docker-Based Development

```bash
# Build all services (web, pc-agent, raspberry-agent)
docker compose build

# Build a specific service
docker compose build web
docker compose build pc-agent
docker compose build raspberry-agent

# Start all services with hot reload
docker compose up --watch

# Start a specific service with hot reload
docker compose up --watch web
docker compose up --watch pc-agent
docker compose up --watch raspberry-agent

# Start services in detached mode
docker compose up -d --build

# Stop all services
docker compose down

# Stop and remove volumes
docker compose down -v

# View logs
docker compose logs -f
docker compose logs -f web
docker compose logs -f pc-agent

# Access a running container
docker compose exec web sh
docker compose exec redis redis-cli -a ${REDIS_PASSWORD}

# Rebuild after code changes
docker compose build --no-cache
```

### Native Go Commands (for local testing)

```bash
# Build all targets
go build ./...

# Run all tests
go test ./...

# Run a single test file
go test ./path/to/file_test.go

# Run a specific test function
go test -run TestFunctionName ./...

# Run tests with race detection
go test -race ./...

# Run tests with verbose output
go test -v ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Validate code
go vet ./...

# Tidy dependencies
go mod tidy
```

## Project Structure

```
internal/
  controller/      # Controller interfaces (public API contracts)
  controller/base/ # Controller implementations
  usecase/         # Use case interfaces (business logic contracts)
  usecase/*/base/  # Use case implementations
  service/         # Service interfaces and implementations
  domain/          # Domain entities and types
  input/           # Input layers (web handlers, async consumers)
  config/          # Configuration management
  service/keystore/redis/  # Keystore implementations
  service/token/jwt/       # Token service implementations
```

Follow **Clean Architecture**: controllers depend on usecases, usecases depend on services, services depend on domain.

## Code Style Guidelines

### Imports
Order imports: standard library, then external, then local. Group with blank lines:
```go
import (
    "context"
    "fmt"

    "github.com/caarlos0/env/v11"

    "github.com/CedricThomas/console/internal/domain"
)
```

### Naming Conventions
- Use `camelCase` for functions, variables, and struct fields
- Use `PascalCase` for types, interfaces, and exported symbols
- Interface names: single verb (e.g., `Metrics`, `Publisher`, `Consumer`, `Keystore`) _without_ `er` suffix
- Interface implementations: lowercase with struct name matching interface (e.g., `type metrics struct{}`)
- Constructor functions: `NewXxx` for public constructors, lowercase `newXxx` for private
- Embed interfaces to compose functionality (e.g., `type web struct{ auth }`)

### Types & Structs
- Define structs with private fields in `base` packages
- Use explicit struct tags for JSON/env parsing
- Domain types live only in `internal/domain/`

```go
type Metrics struct {
    OS          OSName
    CPUUsage    float64
    VRAMUsage   float64
    MemoryUsage float64
}
```

### Error Handling
Wrap errors with context using `%w`:
```go
if err := w.publisher.Publish(ctx, asyncapi.BootChannel, bootCmd); err != nil {
    return fmt.Errorf("publish boot command: %w", err)
}
```

Use `context.Context` as first parameter for all functions that perform I/O.

### Validation
Add `Validate()` methods to request structs:
```go
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
    if r.Username == "" {
        return errors.New("'username' is required")
    }
    return nil
}
```

### Middleware
Implement Fiber middleware as factory functions:
```go
func LoggerMiddleware() fiber.Handler {
    return func(c fiber.Ctx) error {
        start := time.Now()
        err := c.Next()
        duration := time.Since(start)
        log.Printf("[%s] %s %s %d - %v", time.Now().Format(time.RFC3339), c.Method(), c.Path(), c.Response().StatusCode(), duration)
        return err
    }
}
```

### Configuration
Use environment variables with `github.com/caarlos0/env/v11`:
```go
type Config struct {
    RedisURL     string `env:"REDIS_URL,required"`
    JWTSecret    string `env:"JWT_SECRET"`
    Port         string `env:"PORT" envDefault:"3000"`
}
```

### Async & Redis
- Publisher interface: `Publish(ctx, channel, message)`
- Consumer interface: `Subscribe(ctx, channel, callback)` returns unsubscribe function
- Use Redis for both caching and pub/sub messaging

### Authentication
Store password hashes with `bcrypt` from `golang.org/x/crypto`
Generate/verify JWT tokens with `github.com/golang-jwt/jwt/v5`
```

## Running the Application

```bash
# Set required environment variables
export REDIS_URL="redis://localhost:6379"
export JWT_SECRET="your-secret-key"
export PORT="3000"

# Run web server
go run ./cmd/web

# Run PC agent
go run ./cmd/pc-agent

# Run Raspberry agent  
go run ./cmd/raspberry-agent
```

## Dependencies

Key dependencies from `go.mod`:
- `github.com/gofiber/fiber/v3` - Web framework
- `github.com/redis/go-redis/v9` - Redis client
- `github.com/golang-jwt/jwt/v5` - JWT handling
- `golang.org/x/crypto` - bcrypt password hashing
- `github.com/caarlos0/env/v11` - Environment config
- `github.com/robfig/cron/v3` - Cron scheduling

## Rules

- Always rebuild the containers when the task is over

