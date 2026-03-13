# copilot-go-advanced

> **GitHub Copilot Exercise Environment** — a fully working Go weather service built for a workshop that teaches participants to build agentic Copilot workflows.

See [EXERCISES.md](EXERCISES.md) for the workshop exercises.

---

## What It Does

- Fetches real-time weather data from the [OpenWeatherMap 2.5 API](https://openweathermap.org/current)
- Manages saved locations (in-memory, no database)
- Serves a static HTML/JS dashboard with current weather, 5-day forecast (Chart.js), and custom threshold-based alerts
- Provides a REST API with full CRUD for locations and weather queries
- Interactive API docs (Swagger UI) at `/docs/index.html`

---

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.25+ |
| HTTP framework | Gin |
| Validation | go-playground/validator v10 |
| API docs | Swagger (swaggo/swag) |
| Config | godotenv |
| UUID generation | google/uuid |
| Test assertions | testify |
| E2E tests | Playwright (Node.js) |
| Hot reload | Air |
| Lint | golangci-lint |

---

## Prerequisites

| Requirement | Version | Notes |
|-------------|---------|-------|
| [Go](https://go.dev/dl/) | 1.25+ | Required to build and run the server |
| [Node.js](https://nodejs.org/) | 18+ | Required for Playwright E2E tests only |
| [OpenWeatherMap API key](https://openweathermap.org/appid) | — | Free tier is sufficient |
| `make` | Any | macOS/Linux only — Windows users can run Go commands directly (see below) |
| [golangci-lint](https://golangci-lint.run/usage/install/) | Latest | Optional — only needed for linting |
| [Air](https://github.com/air-verse/air) | Latest | Optional — only needed for hot reload (`make dev`) |

**Verify your Go installation (Windows PowerShell):**
```powershell
$env:PATH = [System.Environment]::GetEnvironmentVariable("PATH","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("PATH","User")
go version
```

---

## Quick Start

### 1. Get an OpenWeatherMap API key

Register for a free key at [openweathermap.org/appid](https://openweathermap.org/appid).

### 2. Install dependencies

**macOS / Linux:**
```bash
make install
cd tests/e2e && npx playwright install   # download Playwright browsers
cd ../..
```

**Windows (PowerShell):**
```powershell
$env:PATH = [System.Environment]::GetEnvironmentVariable("PATH","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("PATH","User")

go mod download
go install github.com/air-verse/air@latest
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest
cd tests/e2e; npm install; npx playwright install   # e2e deps + browsers
cd ../..
```

### 3. Generate Swagger docs

**macOS / Linux:**
```bash
make swagger
```

**Windows (PowerShell):**
```powershell
swag init -g cmd/server/main.go --output docs
```

> **Note:** The `docs/` folder is gitignored (generated code). You must run this once after cloning, or the server will not compile.

### 4. Configure environment

```bash
cp .env.example .env
# Edit .env and set OPENWEATHERMAP_API_KEY=<your-key>
```

### 5. Start the dev server

**macOS / Linux:**
```bash
make dev        # hot reload with Air
# or
make build && make start
```

**Windows (PowerShell):**
```powershell
# Set PATH so Go tools are available
$env:PATH = [System.Environment]::GetEnvironmentVariable("PATH","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("PATH","User")

# Run directly (no make required)
go run ./cmd/server

# Or build a binary and run it
go build -o bin\server.exe .\cmd\server
.\bin\server.exe
```

The dashboard is at **http://localhost:8080** and Swagger UI at **http://localhost:8080/docs/index.html**.

---

## Run Tests

**macOS / Linux:**
```bash
make test               # all Go tests (unit + integration)
make test-unit          # converters, repository, services
make test-integration   # HTTP handler tests via httptest
make test-e2e           # Playwright browser tests
```

**Windows (PowerShell):**
```powershell
$env:PATH = [System.Environment]::GetEnvironmentVariable("PATH","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("PATH","User")

go test ./internal/...                          # all Go tests
go test ./internal/utils/... ./internal/repository/... ./internal/services/...  # unit tests
go test ./internal/handlers/...                 # integration tests
cd tests/e2e; npx playwright test               # e2e tests (server auto-starts via Playwright webServer)
```

> **Note:** E2E tests hit the real OpenWeatherMap API. Ensure your `.env` file exists and `OPENWEATHERMAP_API_KEY` is set before running them.

---

## Lint & Format

**macOS / Linux:**
```bash
make lint    # golangci-lint
make fmt     # gofmt + goimports
```

**Windows (PowerShell):**
```powershell
golangci-lint run
gofmt -w .
goimports -w .
```

---

## Project Structure

```
.
├── cmd/server/          Entry point (main.go + Swagger meta-annotations)
├── internal/
│   ├── app/             Gin router factory
│   ├── apperrors/       Domain error types + HTTP dispatcher
│   ├── config/          Settings loaded from .env / env vars
│   ├── container/       DI container (production wiring)
│   ├── handlers/        HTTP handlers (weather + locations)
│   ├── models/          Shared type definitions + validation tags
│   ├── repository/      In-memory location CRUD
│   ├── services/        Business logic, OWM HTTP client
│   ├── testhelpers/     Factory functions + test utilities
│   └── utils/           Pure conversion functions
├── public/              Static frontend (index.html, style.css, app.js)
├── tests/e2e/           Playwright browser tests
├── docs/                Generated Swagger/OpenAPI spec
├── .github/
│   ├── copilot-instructions.md      Always-on Copilot context
│   ├── instructions/                Scoped instruction files
│   └── agents/                      Custom Copilot agents
├── .env.example
├── Makefile
└── EXERCISES.md
```

---

## API Endpoints

### Weather

| Method | Path | Query Params | Description |
|--------|------|-------------|-------------|
| GET | `/api/weather/current` | `lat`, `lon`, `units` | Current weather for coordinates |
| GET | `/api/weather/forecast` | `lat`, `lon`, `days` (1-5), `units` | Multi-day forecast |
| GET | `/api/weather/alerts` | `lat`, `lon` | Threshold-based weather alerts |

### Locations

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/locations` | List all saved locations |
| POST | `/api/locations` | Save a new location |
| GET | `/api/locations/:id` | Get a saved location |
| PUT | `/api/locations/:id` | Update a saved location |
| DELETE | `/api/locations/:id` | Delete a saved location |
| GET | `/api/locations/:id/weather` | Current weather for a saved location |

Temperature units: `celsius` (default) | `fahrenheit` | `kelvin`

---

## Copilot Custom Instructions

This repo uses GitHub Copilot's instruction system to give Copilot always-on context and scoped rules:

| File | Scope | Purpose |
|------|-------|---------|
| `.github/copilot-instructions.md` | All files | Project overview, architecture, run commands |
| `.github/instructions/go.instructions.md` | `**/*.go` | Go naming, error handling, DI, linting |
| `.github/instructions/testing.instructions.md` | `**/*_test.go` | Test framework, mocking, factory usage |
| `.github/instructions/frontend.instructions.md` | `public/**` | Vanilla JS/CSS, no-build-step rules |
| `.github/agents/teacher.agent.md` | — | Exercise Tutor agent (`@teacher`) |

---

## Backlog

- Add persistent storage (SQLite / Postgres) to replace the in-memory repository
- Add user authentication (JWT)
- Add more weather providers (fallback strategy)
- Dark mode toggle
