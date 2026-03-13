---
applyTo: "**/*.go"
---

# Go Coding Conventions

## Language Version & Module System

- Go 1.25+
- Module path: `copilot-go-advanced` (defined in `go.mod`)
- All internal packages use the `copilot-go-advanced/internal/` prefix

## Type Safety

- Avoid `interface{}` / `any` unless marshaling JSON — prefer concrete types and named interfaces
- Define explicit interface types for every dependency boundary that needs to be mocked
- Use struct tags for validation (`binding:"required,min=..."`) and JSON (`json:"fieldName"`)

## Naming Conventions

| Element | Convention | Example |
|---------|-----------|---------|
| Exported types/funcs | `PascalCase` | `WeatherService`, `NewLocationRepository` |
| Unexported types/funcs | `camelCase` | `weatherService`, `newOWMClient` |
| File names | `snake_case` | `weather_service.go`, `location_repo.go` |
| Test files | `<name>_test.go` | `weather_service_test.go` |
| Constants | `PascalCase` (exported) | `UnitCelsius`, `SeverityHigh` |
| Interfaces | Noun (no `I` prefix) | `OWMClient`, `WeatherService` |

## Error Handling

- Services return domain errors from `internal/apperrors/` — never `fmt.Errorf("http 404")` style
- Handlers call `apperrors.HandleError(c, err)` instead of writing responses directly
- Use `errors.As()` for structured error unwrapping (not type assertions)
- Wrap third-party errors in domain errors before returning from the service layer

## Layer Responsibilities

```
Handler → calls Service → calls OWMClient/Repository
   ↓ validates request params via ShouldBindQuery / ShouldBindJSON
   ↓ calls apperrors.HandleError on error
Service → contains ALL business logic (unit conversion, aggregation, alerts)
   ↓ returns domain errors, never HTTP codes
OWMClient → maps HTTP errors to domain errors
Repository → maps "not found" to LocationNotFoundError
```

## Concurrency

- `LocationRepository` uses `sync.RWMutex` — `RLock` for reads, `Lock` for writes
- `OWMClient` uses `http.Client{Timeout: 10 * time.Second}`
- Do not share mutable state between goroutines without synchronization

## Dependency Injection

- Production wiring lives in `internal/container/container.go`
- Tests create mocks that implement the same interface (e.g., `OWMClient`)
- Never use global variables for dependencies; always pass via constructor

## Formatting & Linting

- Format: `gofmt -w . && goimports -w .`
- Lint: `golangci-lint run`
- All exported symbols must have docstrings (`// FunctionName ...`)
- Unused imports are not allowed

## Swagger Annotations

- All handlers have `// @Summary`, `// @Tags`, `// @Produce`, `// @Param`, `// @Success`, `// @Failure` annotations
- Regenerate docs:
  - macOS / Linux: `make swagger`
  - Windows: `swag init -g cmd/server/main.go --output docs`
