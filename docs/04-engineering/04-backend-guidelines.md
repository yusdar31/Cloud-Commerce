# Backend Guidelines

**Project:** CloudCommerce

**Version:** 1.0.0

**Status:** Approved

**Owner:** Backend Team

**Last Updated:** July 2026

---

# 1. Purpose

Dokumen ini mendefinisikan standar pengembangan backend untuk seluruh microservices CloudCommerce.

Tujuan:

- Konsistensi antar service
- Mempermudah onboarding
- Memastikan kualitas kode
- Mendukung Clean Architecture

---

# 2. Technology Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.25 |
| Framework | Gin |
| Database Driver | pgx |
| SQL Generator | SQLC |
| Migration | Goose |
| Validation | validator/v10 |
| Config | Viper |
| Logging | slog |
| JWT | golang-jwt |
| Testing | Go testing + Testcontainers |

---

# 3. Clean Architecture

Setiap service mengikuti struktur:

```
service-name/
│
├── cmd/server/main.go        # Entry point
├── internal/
│   ├── domain/               # Entity, Value Object, Repository interface
│   ├── application/          # Use cases / service layer
│   ├── infrastructure/       # PostgreSQL, Redis, NATS, MinIO implementations
│   └── transport/            # HTTP handler, middleware, request/response
├── configs/                  # Default config files
├── migrations/               # Goose migration files
├── tests/                    # Integration tests
├── Dockerfile
├── Makefile
└── README.md
```

## Layer Rules

- **Domain**: No external dependencies. Pure Go.
- **Application**: Only depends on Domain.
- **Infrastructure**: Implements Domain interfaces (repository, etc.).
- **Transport**: Depends on Application. Handles HTTP concerns only.

---

# 4. Standard File Structure Per Service

```
internal/
│
├── domain/
│   ├── entity.go             # Entity definitions
│   ├── repository.go         # Repository interfaces
│   └── errors.go             # Domain-specific errors
│
├── application/
│   ├── service.go            # Business logic / use cases
│   └── service_test.go
│
├── infrastructure/
│   ├── postgres/
│   │   ├── repository.go     # SQLC generated code usage
│   │   └── db.go             # Connection pool setup
│   ├── redis/
│   │   └── cache.go
│   └── nats/
│       └── publisher.go
│
└── transport/
    ├── handler.go            # Gin handlers
    ├── middleware.go          # Auth, logging, CORS
    ├── request.go            # Request DTO + validation
    └── response.go           # Standard response format
```

---

# 5. Coding Conventions

## Naming

- File: `snake_case.go`
- Folder: `snake_case`
- Interface: `Finder`, `Creator`, `Repository`
- Constructor: `NewService`, `NewHandler`

## Error Handling

```go
// Always return errors
if err != nil {
    return nil, fmt.Errorf("find product by id: %w", err)
}

// Domain errors
var (
    ErrProductNotFound = errors.New("product not found")
    ErrInsufficientStock = errors.New("insufficient stock")
)
```

## Logging

```go
slog.Info("order created",
    "order_id", order.ID,
    "tenant_id", order.TenantID,
    "amount", order.TotalAmount,
)
```

Never log: passwords, tokens, secrets.

---

# 6. Database

## Migration (Goose)

```
migrations/
│
00001_create_users.sql
00002_create_products.sql
00003_add_inventory.sql
```

Rules:

- One file per change
- Idempotent (use `IF NOT EXISTS`)
- Never edit a committed migration
- Always add a rollback

## SQLC

```sql
-- Query
SELECT * FROM products WHERE tenant_id = $1 AND status = $2;
```

```go
// Generated
products, err := queries.FindByTenantAndStatus(ctx, arg)
```

## Connection Pool

```go
config, _ := pgxpool.ParseConfig(databaseURL)
config.MaxConns = 20
config.MinConns = 5

pool, _ := pgxpool.NewWithConfig(ctx, config)
```

---

# 7. API Handler Standards

## Handler Structure

```go
type ProductHandler struct {
    service application.ProductService
}

func NewProductHandler(service application.ProductService) *ProductHandler {
    return &ProductHandler{service: service}
}

func (h *ProductHandler) Create(c *gin.Context) {
    var req CreateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse{Message: err.Error()})
        return
    }

    product, err := h.service.Create(ctx, req.ToDomain())
    if err != nil {
        c.JSON(500, ErrorResponse{Message: err.Error()})
        return
    }

    c.JSON(201, SuccessResponse{Data: product})
}
```

## Response Format

```go
type SuccessResponse struct {
    Data   interface{} `json:"data"`
    Meta   *Meta       `json:"meta,omitempty"`
    Links  *Links      `json:"links,omitempty"`
}

type ErrorResponse struct {
    Type    string `json:"type"`
    Title   string `json:"title"`
    Status  int    `json:"status"`
    Detail  string `json:"detail"`
    TraceID string `json:"traceId,omitempty"`
}
```

---

# 8. Routing

```go
func SetupRouter(handler *ProductHandler, mw *Middleware) *gin.Engine {
    r := gin.New()
    r.Use(mw.CORS(), mw.RequestID(), mw.Logger())

    api := r.Group("/api/v1")
    api.Use(mw.JWTAuth())

    api.GET("/products", handler.List)
    api.POST("/products", handler.Create)
    api.GET("/products/:id", handler.GetByID)
    api.PUT("/products/:id", handler.Update)
    api.DELETE("/products/:id", handler.Delete)

    return r
}
```

---

# 9. Middleware

All services must have:

```go
func (m *Middleware) JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        // validate JWT, extract tenant_id, user_id, role
        c.Set("tenant_id", claims.TenantID)
        c.Set("user_id", claims.UserID)
        c.Set("role", claims.Role)
        c.Next()
    }
}

func (m *Middleware) TenantIsolation() gin.HandlerFunc {
    return func(c *gin.Context) {
        tenantID := c.GetString("tenant_id")
        // ensure all queries filter by tenant_id
        c.Set("tenant_id", tenantID)
        c.Next()
    }
}
```

---

# 10. Configuration

All config via environment variables:

```go
type Config struct {
    Port        string `mapstructure:"PORT"`
    DatabaseURL string `mapstructure:"DATABASE_URL"`
    RedisURL    string `mapstructure:"REDIS_URL"`
    NATSURL     string `mapstructure:"NATS_URL"`
    JWTSecret   string `mapstructure:"JWT_SECRET"`
}
```

Use Viper:

```go
viper.AutomaticEnv()
viper.SetConfigFile(".env")
```

---

# 11. Health Checks

```go
r.GET("/healthz", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "ok"})
})

r.GET("/readyz", func(c *gin.Context) {
    // check DB connection
    c.JSON(200, gin.H{"status": "ready"})
})
```

---

# 12. Testing

```go
func TestCreateProduct(t *testing.T) {
    // use testcontainers for postgres
    ctx := context.Background()
    pgContainer, _ := postgres.Run(ctx, "postgres:17")

    db, _ := sql.Open("postgres", pgContainer.ConnectionString(ctx))
    // run migration
    // seed data
    // test
}
```

Cover target:

- Unit test domain & application: >80%
- Integration test for repository: happy path + error path
- Handler test: status code, response body

---

# 13. Related Documents

- Technology Stack
- Coding Standards
- API Guidelines
- Database Design
- Monorepo Structure
