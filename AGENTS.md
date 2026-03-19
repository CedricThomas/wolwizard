# Console Development Guide

A Go-based Clean Architecture console for managing PC and Raspberry Pi agents via Redis pub/sub.

## Build & Test Commands

### Docker Development
```bash
# Build all services
docker compose build

# Start with hot reload
docker compose up --watch web

# Stop services
docker compose down

# View logs
docker compose logs -f web

# Rebuild after changes
docker compose build --no-cache
```

### Native Go Commands
```bash
# Build
go build ./...

# Run all tests
go test ./...

# Run specific test
go test -run TestFunctionName ./...

# Test with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Validate & tidy
go vet ./...
go mod tidy
```

## Project Structure

```
internal/
  controller/      # Public API interfaces
  controller/base/ # Controller implementations
  usecase/         # Business logic interfaces
  usecase/*/base/  # Use case implementations
  service/         # Service interfaces & implementations
  domain/          # Domain entities
  input/           # Web handlers, async consumers, cron
  config/          # Configuration management
```

Follow **Clean Architecture**: controllers depend on usecases, usecases depend on services, services depend on domain.

## Code Style Guidelines

### Imports
Order: standard library, external, local. Separate groups with blank lines:
```go
import (
    "context"
    "fmt"
    "net"

    "github.com/caarlos0/env/v11"

    "github.com/CedricThomas/console/internal/domain"
)
```

### Naming Conventions
- `camelCase`: functions, variables, struct fields
- `PascalCase`: types, interfaces, exported symbols
- Interfaces: single verb without `-er` suffix (`Metrics`, `Publisher`, `Keystore`)
- Implementations: lowercase, same name (`type metrics struct{}`)
- Constructors: `NewXxx` (public), `newXxx` (private)
- Embed interfaces to compose: `type web struct{ auth }`

### Types & Structs
- Private fields in `base` packages
- Use struct tags for JSON/env parsing
- Domain types live only in `internal/domain/`

```go
type Config struct {
    RedisURL     string `env:"REDIS_URL,required"`
    JWTSecret    string `env:"JWT_SECRET"`
    Port         string `env:"PORT" envDefault:"3000"`
    MACAddress   net.HardwareAddr `env:"-"`
}
```

### Error Handling
Wrap errors with context using `%w`. First param is `context.Context`:
```go
func (w web) SendAsyncBootCommand(ctx context.Context, osName domain.OSName) error {
    bootCmd := asyncapi.BootCommand{OSName: osName}
    if err := w.publisher.Publish(ctx, asyncapi.BootChannel, bootCmd); err != nil {
        return fmt.Errorf("publish boot command: %w", err)
    }
    return nil
}
```

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

### Fiber Middleware
Factory functions returning `fiber.Handler`:
```go
func LoggerMiddleware() fiber.Handler {
    return func(c fiber.Ctx) error {
        start := time.Now()
        err := c.Next()
        log.Printf("%s %s %d - %v", c.Method(), c.Path(), c.Response().StatusCode(), time.Since(start))
        return err
    }
}
```

## Services & Interfaces

### Redis (Pub/Sub)
- `Publisher.Publish(ctx, channel, message)`
- `Consumer.Subscribe(ctx, channel, callback)` returns `unsubscribe func()`

### Authentication
- Password hashing: `golang.org/x/crypto/bcrypt`
- JWT tokens: `github.com/golang-jwt/jwt/v5`

### WebSocket
- Use `github.com/gofiber/contrib/websocket` for WS connections
- Manager interface handles connections per client

## Key Dependencies
- `github.com/gofiber/fiber/v3` - Web framework
- `github.com/redis/go-redis/v9` - Redis client
- `github.com/golang-jwt/jwt/v5` - JWT handling
- `golang.org/x/crypto` - bcrypt
- `github.com/caarlos0/env/v11` - Environment config
- `github.com/robfig/cron/v3` - Cron scheduling

## Running
```bash
docker compose up -d
```

## Rules
- Always rebuild Docker containers after completing tasks:
  ```bash
  docker compose build --no-cache
  docker compose up --watch
  ```
}