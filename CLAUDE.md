# fit-profi — Claude Code Instructions

## Project Overview

Fitness platform API. Stack: **Go 1.25 · Echo v4 · GORM · MySQL**
Module: `github.com/msskobelina/fit-profi`
Architecture: **DDD + CQRS** inside `internal/`

---

## Commands

```bash
# Run server
go run ./cmd/api/main.go

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/delivery/controller/authorize/...

# Lint
golangci-lint run

# Tidy deps
go mod tidy
```

---

## Architecture

```
cmd/api/main.go                  → bootstrap.Run()
internal/
├── bootstrap/                   # DI wiring + Echo setup
├── domain/
│   ├── model/                   # Pure domain entities (no framework deps)
│   └── repository/              # Repository interfaces
├── application/
│   ├── command/                 # Mutating use-cases (by domain)
│   └── query/                   # Read use-cases (by domain)
├── delivery/
│   ├── boundary/io.go           # JSON decode + go-playground/validator
│   ├── controller/              # net/http handlers (framework-agnostic)
│   ├── middleware.go            # AuthMiddleware, RequireAdmin
│   └── router.go               # Echo route registration + wrap() helper
└── infrastructure/
    ├── email/                   # SMTP sender
    └── repository/              # GORM implementations
pkg/                             # Shared utilities (access, analytics, errors, …)
```

### Key Patterns

- **Controllers** are `http.Handler` — no Echo imports in `delivery/controller/`
- **Echo path params** injected via `controller.PathParam(r, "id")`
- **Auth context**: `"userID"` (`int`) and `"userRole"` (`string`) set by `AuthMiddleware`
- **IO interface** (`delivery/controller/controller.go`): `Read`, `Error`, `Fatal`, `Result`
- **Validation**: struct tags `validate:"required,email"` etc., processed by `boundary/io.go`

---

## Conventions

### Adding a new endpoint

1. **Domain model** in `internal/domain/model/` if needed
2. **Repository interface** in `internal/domain/repository/`
3. **Command/Query** struct + handler in `internal/application/command/` or `query/`
4. **Controller** in `internal/delivery/controller/<domain>/`
5. **GORM implementation** in `internal/infrastructure/repository/<domain>/`
6. **Wire it** in `bootstrap/application.go` + register route in `delivery/router.go`

### Validation rules (required for all write endpoints)

All request structs **must** carry `validate` tags:
- `validate:"required"` — mandatory fields
- `validate:"email"` — email fields
- `validate:"min=6"` — password min length
- `validate:"oneof=val1 val2"` — enum fields (Goal, MealType, Category)
- `validate:"gt=0"` — positive numeric fields (Age, WeightKg)

### Error handling

- Use `io.Error(err, r, w)` for 400 — includes validation errors
- Use `io.Fatal(err, r, w)` for unexpected 500 errors
- Domain errors: return `*errors.Error{Message: "...", Status: http.StatusXxx}` from handlers

### Testing

- Table-driven tests in `*_test.go` alongside the controller file
- Use `boundary.New()` as the real IO implementation in tests
- Mock handlers with inline structs implementing the handler interface
- Use `httptest.NewRecorder()` and `http.NewRequest()`

---

## Security Notes

- **UserID always from auth context** — never from request body
- **Validation required** on every write endpoint to prevent bad data
- **bcrypt cost 14** for passwords
- **JWT HS256** signed with `HMAC_SECRET`, 14-day expiry

---

## Environment Variables

```
MYSQL_USER, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_DATABASE
HTTP_PORT
HMAC_SECRET
ADMIN_USER_FULLNAME, ADMIN_USER_EMAIL
MAIL_HOST, MAIL_USERNAME, MAIL_APP_PASSWORD
MIXPANEL_TOKEN, MIXPANEL_API_HOST
GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, GOOGLE_REDIRECT_URL
LOG_LEVEL
```
