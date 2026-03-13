---
applyTo: "**/*_test.go"
---

# Testing Conventions

## Framework & Tools

- **Standard library:** `testing`
- **Assertions:** `github.com/stretchr/testify/assert` (non-fatal) and `testify/require` (fatal)
- **HTTP recorder:** `net/http/httptest`
- **Helper utilities:** `internal/testhelpers` package

## Test Organization

| Type | Location | Pattern |
|------|----------|---------|
| Unit tests | `internal/<package>/` (same package as code) | `<file>_test.go` |
| Integration tests | `internal/handlers/integration_test.go` | `package handlers_test` |
| E2E (Playwright) | `tests/e2e/` | `*.spec.ts` |

## Running Tests

**macOS / Linux:**
```bash
make test               # all Go tests (unit + integration)
make test-unit          # utils, models, repository, services
make test-integration   # handlers (HTTP level, mock OWM client)
make test-e2e           # Playwright browser tests
```

**Windows (PowerShell):**
```powershell
go test ./internal/...                                                           # all Go tests
go test ./internal/utils/... ./internal/repository/... ./internal/services/...  # unit tests
go test ./internal/handlers/...                                                  # integration tests
cd tests/e2e; npx playwright test                                                # e2e tests
```

> **Note:** E2E tests hit the real OpenWeatherMap API. Ensure your `.env` file exists with a valid `OPENWEATHERMAP_API_KEY`.

## Naming Conventions

- Test functions: `TestSubject_Condition` → `TestWeatherCurrent_404WhenOWMReturns404`
- Test files mirror the source file: `weather_service.go` → `weather_service_test.go`
- Mock types are named `mock<InterfaceName>` (camelCase, unexported within the test package)

## AAA Pattern

All test cases follow **Arrange → Act → Assert**:

```go
func TestLocationRepo_GetNonExistent(t *testing.T) {
    // Arrange
    repo := repository.NewLocationRepository()

    // Act
    _, err := repo.Get("non-existent-id")

    // Assert
    var notFound *apperrors.LocationNotFoundError
    require.ErrorAs(t, err, &notFound)
    assert.Equal(t, "non-existent-id", notFound.ID)
}
```

## Factory Usage

Always use factory functions from `internal/testhelpers/factories.go` — never inline raw struct literals in test bodies.

```go
// ✅ correct — use factory with override
owmResp := testhelpers.MakeOwmCurrentWeatherResponse(func(r *models.OwmCurrentWeatherResponse) {
    r.Main.Temp = 42.0
})

// ❌ wrong — inline raw struct
owmResp := models.OwmCurrentWeatherResponse{Main: models.OwmMain{Temp: 42.0}}
```

## Mocking Strategy

| Test type | Mock boundary | How |
|-----------|--------------|-----|
| Unit (service) | `OWMClient` | `mockOWMClient` struct in `*_test.go`, implements `services.OWMClient` |
| Integration (handler) | `OWMClient` | Same mock injected into real `WeatherService` via `services.NewWeatherService(mock, settings)` |
| E2E | None — real server | Playwright auto-starts the server via `go run` using the `webServer` directive — no manual build needed |

## Test Helpers

```go
// Build a Gin engine with a real WeatherService + mock OWM client
r := testhelpers.NewTestApp(repo, weatherSvc)

// Fire a request and get the recorder
w := testhelpers.PerformRequest(r, http.MethodGet, "/api/weather/current?lat=51.51&lon=-0.13", nil)

// Fire a POST with a JSON body
body := `{"name":"London","lat":51.51,"lon":-0.13}`
w := testhelpers.PerformRequest(r, http.MethodPost, "/api/locations", &body)
```

## Integration Test Pattern

```go
// Each test creates its own isolated router (no shared state)
mock := &mockOWMClient{currentResp: &owmResp}
repo := repository.NewLocationRepository()
settings := testhelpers.TestSettings()
weatherSvc := services.NewWeatherService(mock, settings)
r := testhelpers.NewTestApp(repo, weatherSvc)

w := testhelpers.PerformRequest(r, http.MethodGet, "/api/weather/current?lat=51.51&lon=-0.13", nil)
assert.Equal(t, http.StatusOK, w.Code)
```

## No Real API Calls

Tests must **never** make real HTTP calls to OpenWeatherMap. Always inject a mock `OWMClient`. If you see a real API call in a test, that is a bug.
