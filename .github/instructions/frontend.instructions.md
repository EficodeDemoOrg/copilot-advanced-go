---
applyTo: "public/**"
---

# Frontend Conventions

## Technology

The frontend is intentionally **simple** — vanilla JavaScript, CSS, and HTML with no build step, no framework, and no compilation.

| Asset | File | Purpose |
|-------|------|---------|
| Page | `public/index.html` | Single-page dashboard |
| Styles | `public/style.css` | Layout, theming, component styles |
| Logic | `public/app.js` | All JS: fetching, rendering, chart |

## Serving

The Go backend serves the `public/` directory as static files via Gin:

```go
r.Static("/static", "./public")
r.StaticFile("/", "./public/index.html")
```

The dashboard is available at `http://localhost:8080/`.

## No Build Step

- Do **not** introduce npm, webpack, Vite, or any bundler for the frontend.
- Do **not** add TypeScript or JSX to the frontend files.
- Dependencies (e.g., Chart.js) are loaded via CDN `<script>` tags in `index.html`.

## CSS Custom Properties

Theming uses CSS custom properties defined on `:root`:

```css
:root {
  --primary: #2563eb;
  --bg: #f8fafc;
  --surface: #ffffff;
  --border: #e2e8f0;
  /* ... */
}
```

Add new theme values as custom properties — never hardcode hex values in rules.

## JavaScript Patterns

- All API calls use `fetch()` with `async/await`
- DOM manipulation uses `document.querySelector` / `getElementById` — no jQuery
- Error states are rendered inline in the UI (not via `alert()`)
- Chart.js is used for the 5-day forecast bar/line chart

## API Endpoints Used

The JS makes calls to the Go REST API:

| JS call | Endpoint |
|---------|----------|
| Current weather | `GET /api/weather/current?lat=…&lon=…&units=…` |
| Forecast | `GET /api/weather/forecast?lat=…&lon=…&days=5` |
| Alerts | `GET /api/weather/alerts?lat=…&lon=…` |
| List locations | `GET /api/locations` |
| Create location | `POST /api/locations` |
| Delete location | `DELETE /api/locations/:id` |
| Location weather | `GET /api/locations/:id/weather` |

Location coordinates are accessed as `loc.coordinates.lat` and `loc.coordinates.lon` (matching the Go JSON output).

## E2E Testing

The frontend is tested with **Playwright** (`tests/e2e/dashboard.spec.ts`). The test suite:
- Validates the dashboard loads
- Tests coordinate search form interaction
- Tests saved-location CRUD via the API
- Checks forecast chart canvas is present
- Checks alerts container exists

**macOS / Linux:**
```bash
make test-e2e
```

**Windows (PowerShell):**
```powershell
cd tests/e2e; npx playwright test
```

Requires Node.js. The Go server is auto-started via `go run` using Playwright's `webServer` directive — no manual build needed. Ensure your `.env` file exists with a valid `OPENWEATHERMAP_API_KEY`.
