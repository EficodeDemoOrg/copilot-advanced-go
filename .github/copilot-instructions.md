# Weather App — Copilot Instructions

## Project Overview

This is **not** a production application. It is a **GitHub Copilot exercise environment** — a fully working Go weather service that serves as the substrate for a workshop teaching participants to build agentic workflows with GitHub Copilot: custom instructions, agents, skills, MCP integration, and hooks.

The application code is complete and tested. Participants should never need to fix application bugs — they build Copilot tooling **around** a working codebase.

---

## Architecture

The app uses a strict **layered architecture**:

| Layer | Package | Responsibility |
|-------|---------|----------------|
| Entry point | `cmd/server/` | Starts HTTP server, loads config |
| HTTP Handlers | `internal/handlers/` | Validate input, call services, return responses. **No business logic.** |
| Services | `internal/services/` | Business logic. `WeatherService` handles unit conversion, forecast aggregation, and threshold-based alerts. `OWMClient` handles all external HTTP. |
| Repository | `internal/repository/` | In-memory CRUD for `Location` objects. Uses `sync.RWMutex`. No database. |
| Models | `internal/models/` | Shared type definitions with Gin binding tags for validation. |
| App factory | `internal/app/` | Wires Gin router, routes, static files, Swagger UI. |
| DI Container | `internal/container/` | Creates real dependencies from `Settings`. Tests swap these with mocks. |
| Config | `internal/config/` | Loads settings from `.env` / environment variables. |
| Errors | `internal/apperrors/` | Domain error types + `HandleError()` dispatcher (domain errors → HTTP codes). |
| Utils | `internal/utils/` | Pure, stateless conversion functions. No side effects. |
| Test helpers | `internal/testhelpers/` | `TestSettings()`, `NewTestApp()`, `PerformRequest()`, factory functions. |
| Frontend | `public/` | Vanilla JS/CSS/HTML served as static files. No build step. |

### Key Rules

- **Handlers** never contain business logic.
- **Services** never return HTTP errors — they return domain errors.
- `HandleError()` in `apperrors` translates domain errors to HTTP status codes.
- The repository is entirely in-memory; no database, no persistence.
- Dependency injection is achieved by accepting interfaces (`OWMClient`, `WeatherService`, `LocationRepository`) — tests swap in mocks.

---

## Key Conventions

- **Language:** Go 1.25+
- **HTTP framework:** Gin (`github.com/gin-gonic/gin`)
- **Validation:** Gin binding tags via `go-playground/validator/v10`
- **Linter:** `golangci-lint`
- **Formatter:** `gofmt` + `goimports`
- **Naming:** `PascalCase` for exported types/functions, `camelCase` for unexported, `snake_case` for file names

---

## Custom Domain Errors

All domain errors live in `internal/apperrors/errors.go`:

| Error Type | HTTP Code | When |
|-----------|-----------|------|
| `WeatherAPINotFoundError{Lat, Lon}` | 404 | OWM returns 404 |
| `WeatherAPIConnectionError{Cause}` | 503 | Network/timeout error |
| `WeatherAPIError{StatusCode, Message}` | 502 | Any other non-200 from OWM |
| `LocationNotFoundError{ID}` | 404 | Location ID missing from repo |
| `WeatherAppError{Message}` | 500 | Generic application error |

`HandleError(c *gin.Context, err error)` uses `errors.As()` to match and dispatch.

---

## Testing Overview

- **Framework:** `testing` (stdlib) + `testify/assert` + `testify/require`
- **HTTP testing:** `net/http/httptest` via `PerformRequest()` helper
- **No real API calls** — all tests mock the OWM client
- **Unit tests:** mock at the `WeatherService` or `OWMClient` interface boundary
- **Integration tests:** inject `mockOWMClient` into real `WeatherService`, use `NewTestApp()`
- **E2E tests:** Playwright (`tests/e2e/`) — start Go server via `webServer` directive

---

## Dependencies

| Dependency | Purpose |
|-----------|---------|
| `github.com/gin-gonic/gin` | HTTP framework + router |
| `github.com/go-playground/validator/v10` | Struct validation via binding tags |
| `github.com/swaggo/swag` | Swagger annotation → OpenAPI spec generator |
| `github.com/swaggo/gin-swagger` | Swagger UI middleware for Gin |
| `github.com/joho/godotenv` | Load `.env` file |
| `github.com/google/uuid` | UUID generation for location IDs |
| `github.com/stretchr/testify` | Test assertions |
| `@playwright/test` (Node.js) | Frontend browser e2e tests |

---

## Run Commands

```bash
# Install all dependencies (Go + Node e2e)
make install

# Start dev server with hot reload
make dev

# Build binary
make build

# Run all Go tests (unit + integration)
make test

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Run Playwright e2e tests
make test-e2e

# Lint
make lint

# Format
make fmt

# Regenerate Swagger docs
make swagger
```
