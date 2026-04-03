# Console Development Guide

A distributed Go-based system for remote PC management, dual-boot OS switching, and Wake-on-LAN control following Clean Architecture principles.

## 🏗️ Architecture Overview

### Three Main Services

| Service | Port | Platform | Purpose |
|---------|-----|----- |---------|
| **Web** | 8080 | Any | REST API, WebSocket, JWT Auth |
| **PC Agent** | 8081 | Linux/Windows | Commands, Metrics, Boot Management |
| **Raspberry Agent** | N/A | Linux/Windows | Wake-on-LAN (WoL) packets |

### Communication Pattern

```
Web Service ──(HTTP/WS)──► Browser
     │
     │ publish to Redis
     ▼
  Redis (pub/sub)
     │
     ├─► PC Agent (subscribes: boot, shutdown)
     │     └─► publishes metrics/pc-agent
     │
     └─► Raspberry Agent (subscribes: wol)
```

### Async Message Channels

| Channel | Direction | Payload Type | Description |
|---------|------ |----- | |---------|----------|----- |---------|
| `boot` | Web → PC | `BootCommand` | Set next boot OS |
| `shutdown` | Web → PC | `ShutdownCommand` | System shutdown |
| `metrics/pc-agent` | PC → Web | `MetricsCommand` | System metrics push |

---

## 🛠️ Build & Test Commands

### Docker Development
```bash
# Build all services
docker compose build

# Start with hot reload
docker compose up --watch

# Start specific service
docker compose up --watch web

# View all logs
docker compose logs -f

# Rebuild (clean cache)
docker compose build --no-cache

# Stop all
docker compose down
```

### Native Go Commands
```bash
# Build all binaries
go build ./...

# Run all tests
go test ./...

# Run specific test
go test -run TestFunctionName ./...

# Verbose test output
go test -v -run TestFunctionName ./...

# Test with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Linting and cleanup
go vet ./...
go mod tidy
```

### Building Individual Services
```bash
# Web server
go build -o web ./cmd/web

# PC agent (cross-platform)
GOOS=linux go build -o pc-agent ./cmd/pc-agent
GOOS=windows go build -o pc-agent.exe ./cmd/pc-agent

# Raspberry agent
go build -o raspberry-agent ./cmd/raspberry-agent
```

---

## 🔧 Mock Generation

All service/usecase interfaces use `go:generate` with [`mockgen`](https://pkg.go.dev/go.uber.org/mock/mockgen).

### Install mockgen (one-time)
```bash
go install go.uber.org/mock/mockgen@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Generate All Mocks
```bash
go generate ./...
```

### Generate for Specific Packages
```bash
go generate ./internal/service/...     # Service interface mocks
go generate ./internal/usecase/...      # Use case interface mocks
go generate ./internal/controller       # Controller interface mocks
go generate ./internal/input/async      # Message consumer mocks
```

**Note:** Mock files are stored in each package's `mock/` subdirectory. **Do not edit them manually** - they will be overwritten.

### Mock Package Structure
```
service/
  command/
    command.go          # Interface definition
    mock/
      command.go        # Auto-generated MockCommandExecutor
```

---

## 📁 Project Structure

```
console/
├── cmd/                        # Application entry points
│   ├── web/                    # Web server main()
│   ├── pc-agent/               # PC agent main()
│   └── raspberry-agent/        # WoL agent main()
├── internal/
│   ├── config/                 # Environment config, Redis client
│   ├── domain/                 # Domain entities, custom errors
│   │   └── boot.go             # BootEntry, OSName, BootError
│   │   └── metrics.go          # Metrics struct
│   ├── usecase/                # Business logic interfaces
│   │   ├── auth/              # UserAuth interface
│   │   │   └── base/          # Auth implementation
│   │   ├── boot/              # Boot interface
│   │   │   └── base/          # Boot implementation
│   │   └── token/             # Token interface
│   ├── service/                # External service interfaces
│   │   ├── async/             # Redis publisher
│   │   ├── command/           # Shutdown, reboot, GRUB
│   │   │   ├── linux/         # Linux implementation
│   │   │   └── windows/       # Windows implementation
│   │   ├── keystore/          # Redis key-value store
│   │   ├── metrics/           # Platform-specific collectors
│   │   │   ├── linux/         # Linux metrics
│   │   │   └── windows/       # Windows metrics
│   │   ├── token/jwt/         # JWT service
│   │   ├── websocket/         # WebSocket manager
│   │   └── wol/               # WoL magic packet sender
│   ├── controller/             # Controller interfaces
│   │   └── base/              # Controller implementations
│   │       ├── web.go         # WebController
│   │       ├── pc_agent.go    # PC Agent Controller
│   │       ├── raspberry_agent.go
│   │       └── auth.go        # AuthController
│   └── input/                  # Input layer (external interface)
│       ├── web/               # HTTP handlers, middleware, router
│       │   ├── fiber/         # Fiber framework setup
│       │   │   ├── handler/   # HTTP request handlers
│       │   │   ├── middleware/  # CORS, auth, logging
│       │   │   └── router/    # Route definitions
│       │   └── presenters/    # Response DTOs
│       ├── async/             # Redis pub/sub handlers
│       │   ├── api/           # Async message DTOs
│       │   ├── redis/         # Redis pub/sub consumer
│       │   ├── handlers/      # Message handlers
│       │   └── subscriptions/ # Channel subscriptions
│       └── cron/              # Periodic jobs
│           └── robfig/        # Cron implementation
├── static/                    # Web frontend assets
├── docker-compose.yml
├── go.mod
└── .env.example
```

### Dependency Flow (Clean Architecture)

```
Input Layer (handlers, cron)
        ↓
Controller Layer (web, pc_agent, raspberry)
        ↓
UseCase Layer (auth, boot, metrics, token)
        ↓
Service Layer (command, metrics, keystore, wol, websocket, async)
        ↓
Domain Layer (entities: BootEntry, Metrics, OSName)
```

**Rule:** Inner layers have no knowledge of outer layers

---

## 📝 Code Style Guidelines

### Imports
**Order:** Standard library → External → Local (separated by blank lines)

```go
import (
    "context"
    "errors"
    "fmt"
    "log"
    "net"
    "time"

    "github.com/caarlos0/env/v11"
    "github.com/gofiber/fiber/v3"
    "github.com/golang-jwt/jwt/v5"

    "github.com/CedricThomas/console/internal/domain"
)
```

### Naming Conventions

| Pattern | Usage | Examples |
|---------|-------|--|------|
| `camelCase` | Functions, variables, struct fields | `username`, `jwtSecret` |
| `PascalCase` | Types, interfaces, exported | `UserController`, `AuthService` |
| `Interface name` | Single verb, **no -er suffix** | `Publisher`, `Consumer`, `Collector` |
| `Implementation` | Lowercase same name | `type publisher struct{}` |
| `Constructor` | `NewXxx()` public, `newXxx()` private | `NewRedisClient()`, `newToken()` |
| `Interface embed` | Use for composition | `type webController struct{ authUseCase UserAuth }` |

### Types & Structs

```go
// Domain entity (internal/domain/)
type Metrics struct {
    OS          string  `json:"os"`
    CPUUsage    float64 `json:"cpu_usage"`
    MemoryUsage float64 `json:"memory_usage"`
    VRAMUsage   float64 `json:"vram_usage"`
}

// Config with env tags
type WebConfig struct {
    JWTSecret             string        `env:"JWT_SECRET" envDefault:"secret"`
    JWTExpirySeconds      int64         `env:"JWT_EXPIRY_SECONDS" envDefault:"86400"`
    Port                  string        `env:"PORT" envDefault:"8080"`
    LastMetricsTTLSeconds int           `env:"LAST_METRICS_KEY_TTL_SECONDS" envDefault:"5"`
}

// Async message DTO
type BootCommand struct {
    OSName domain.OSName `json:"os_name"`
}

// Request validation
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
    if r.Username == "" {
        return errors.New("'username' is required")
    }
    if len(r.Password) < 8 {
        return errors.New("'password' must be at least 8 characters")
    }
    return nil
}
```

### Error Handling
**Always wrap errors with context using `%w`.**

```go
func (c *webController) SendBootCommand(ctx context.Context, osName domain.OSName) error {
    bootCmd := asyncapi.BootCommand{OSName: osName}
    
    if err := c.publisher.Publish(ctx, asyncapi.BootChannel, bootCmd); err != nil {
        return fmt.Errorf("publish boot command: %w", err)
    }
    return nil
}

// Custom errors in domain packages
type ErrGRUBEntryNotFound struct {
    Entry string
}

func (e *ErrGRUBEntryNotFound) Error() string {
    return fmt.Sprintf("GRUB entry not found: %s", e.Entry)
}
```

### Method Signatures
**First parameter should always be `context.Context`:**

```go
// Good
func (w *webController) Authenticate(ctx context.Context, req LoginRequest) error

// Bad
func (w *webController) Authenticate(req LoginRequest) error
```

### Fiber Middleware
**Factory functions returning `fiber.Handler`:**

```go
func LoggerMiddleware() fiber.Handler {
    return func(c fiber.Ctx) error {
        start := time.Now()
        err := c.Next()
        log.Printf("%s %s %d - %v", 
            c.Method(), 
            c.Path(), 
            c.Response().StatusCode(),
            time.Since(start))
        return err
    }
}

func AuthMiddleware(next fiber.Handler, verifier func(string) (bool, string)) fiber.Handler {
    return func(c fiber.Ctx) error {
        token := c.Get("Authorization")
        token = strings.TrimPrefix(token, "Bearer ")
        
        if valid, username := verifier(token); !valid {
            return fiber.ErrUnauthorized
        }
        
        c.Locals("username", username)
        return next(c)
    }
}
```

---

## 🔍 Services & Interfaces

### Redis Pub/Sub
**Publisher:**
```go
type Publisher interface {
    Publish(ctx context.Context, channel string, message interface{}) error
}

// Usage
if err := publisher.Publish(ctx, "boot", BootCommand{OSName: "linux"}); err != nil {
    return fmt.Errorf("publish: %w", err)
}
```

**Consumer:** 
```go
type Callback func(ctx context.Context, message string) error

type Consumer interface {
    Subscribe(ctx context.Context, channel string, callback Callback) (func() error, error)
}

// Usage
unsub, err := consumer.Subscribe(ctx, "metrics/pc-agent", handleMetrics)
if err != nil {
    return err
}
// Cleanup
defer unsub()
```

### Authentication Flow

```go
// Registration
POST /auth/register
{ "username": "user1", "password": "strongpassword" }
⬇️
bcrypt hash password
⬇️
Store in Redis: redis.Set("user:user1", hashedPassword)
⬇️
Return JWT token

// Login
POST /auth/login
{ "username": "user1", "password": "mypassword" }
⬇️
Fetch hash from Redis
⬇️
bcrypt.CompareHashAndPassword
⬇️
Return JWT token

// Token Claims
{
    "sub": "user1",
    "exp": 1234567890,
    "iat": 1234567800
}
```

### WebSocket Metrics Flow

```
PC Agent (cron job)
    ↓ publish to "metrics/pc-agent"
Redis
    ↓ broadcast to subscribers
Web Service (websocket manager)
    ↓ push to connected clients
Browser
    ↓ update UI with latest metrics
```

### Boot Command Flow (GRUB)

```
Web API POST /boot {"os_name": "linux"}
    ↓
Publish to "boot" channel
    ↓
PC Agent receives
    ↓
Find GRUB entry name
    ↓
grub-reboot "entry-name"
    ↓
Set for next reboot
    ↓
PC Agent reboots
    ↓
Ubuntu boots
```

---

## 🔑 Key Dependencies

| Package | Version | Purpose |
|---------|---- |---------|
| `github.com/gofiber/fiber/v3` | v3.1.0 | High-performance web framework |
| `github.com/redis/go-redis/v9` | v9.18.0 | Redis client (pub/sub, cache) |
| `github.com/golang-jwt/jwt/v5` | v5.3.1 | JWT token creation/verification |
| `golang.org/x/crypto` | v0.49.0 | bcrypt password hashing |
| `github.com/caarlos0/env/v11` | v11.4.0 | Environment variable parsing |
| `github.com/robfig/cron/v3` | v3.0.1 | Cron job scheduling |
| `github.com/fasthttp/websocket` | v1.5.8 | WebSocket connections |
| `go.uber.org/mock` | v0.6.0 | Mock generation |
| `github.com/stretchr/testify` | v1.11.1 | Testing utilities |

---

## 🚀 Running

### Docker (Recommended)
```bash
cd console

# First time setup
docker compose build

# Start all services
docker compose up --watch

# Access web UI
open http://localhost:8080
```

### Environment Setup
Create `.env` in project root:

```bash
# Redis
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=yourpassword

# Web Server
PORT=8080
JWT_SECRET=a_secret_key_minimum_32_chars_long
JWT_EXPIRY_SECONDS=86400

# PC Agent
METRICS_REPORTING_SCHEDULE=@every 5s
BOOT_OS_TTL_SECONDS=300

# Raspberry Agent
SERVER_MAC_ADDRESS=aa:bb:cc:dd:ee:ff
SERVER_NETWORK_ADDRESS=255.255.255.255:9
```

### Manual Testing
```bash
# Terminal 1: Start Redis
docker run -d -p 6379:6379 redis:alpine

# Terminal 2: Web Server
export REDIS_URL='redis://localhost:6379'
export JWT_SECRET='your_secret_key'
export PORT='8080'
go run ./cmd/web

# Terminal 3: PC Agent
export REDIS_URL='redis://localhost:6379'
go run ./cmd/pc-agent

# Terminal 4: Raspberry Agent
export SERVER_MAC_ADDRESS='aa:bb:cc:dd:ee:ff'
export SERVER_NETWORK_ADDRESS='255.255.255.255:9'
go run ./cmd/raspberry-agent
```

---

## 📋 Checklist for New Features

- [ ] Define interface in appropriate layer (service/usecase)
- [ ] Run `go generate ./...` to create mocks
- [ ] Implement business logic in `base` package
- [ ] Create controller method
- [ ] Add handlers/routes (if web)
- [ ] Add subscriptions (if async)
- [ ] Add unit tests
- [ ] Update documentation
- [ ] Rebuild Docker: `docker compose build --no-cache`

---

## 🛡️ Security Rules

1. **Never commit `.env` files** - use `.env.example` as template
2. **Use strong JWT secrets** (32+ random characters)
3. **Redis must be password-protected**
4. **Input validation** on all endpoints
5. **Use HTTPS** in production
6. **Network isolation** for WoL broadcasts

---

## 🐛 Troubleshooting

### Redis Connection Fails
```bash
# Verify Redis is running
curl http://localhost:6379

# Check REDIS_URL in .env file
```

### Mocks Not Working
```bash
# Regenerate all mocks
go generate ./...

# Verify mockgen installed
which mockgen
```

### Metrics Not Updating
1. Check `METRICS_REPORTING_SCHEDULE` cron syntax
2. Verify WebSocket connection in browser DevTools
3. Check Redis pub/sub: `redis-cli PUBLISH metrics/pc-agent "test"`

### WoL Not Working
1. Verify MAC address is correct
2. Check network allows UDP broadcasts
3. Ensure PC BIOS has WoL enabled
4. Test with: `wakeonlan AA:BB:CC:DD:EE:FF`

---

## ✅ Development Rules

**After every code change:**

1. ✅ Run tests: `go test ./...`
2. ✅ Generate mocks if interfaces changed: `go generate ./...`
3. ✅ Tidy dependencies: `go mod tidy`
4. ✅ Rebuild Docker: `docker compose build --no-cache`
5. ✅ Verify with: `docker compose logs -f`

**Code Style:**
- ✅ Wrap all errors: `fmt.Errorf("message: %w", err)`
- ✅ First param `context.Context`
- ✅ Single responsibility per function
- ✅ Private structs in `base` packages
- ✅ Public interfaces define contracts

---

## 📞 Getting Help

- Check existing issues on the repository
- Review architecture diagrams in README.md
- Ensure you followed the checklist above
- Verify environment variables are correctly set

Happy coding! 🐳
