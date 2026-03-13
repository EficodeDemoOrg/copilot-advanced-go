import { defineConfig, devices } from "@playwright/test";

/**
 * Playwright configuration for the Go weather app dashboard e2e tests.
 *
 * The webServer directive starts the compiled Go binary before running tests
 * and shuts it down afterwards — no manual server start required.
 *
 * Prerequisites:
 *   1. Build the Go server:  cd ../..  &&  go build -o tmp/server ./cmd/server
 *   2. Set OWM env vars or create a .env file (tests use real API unless mocked)
 *
 * Run:  cd tests/e2e  &&  npm install  &&  npx playwright test
 */
export default defineConfig({
  testDir: ".",
  testMatch: "**/*.spec.ts",

  // Fail the whole run on the first test failure in CI
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  workers: 1,

  reporter: [["list"], ["html", { open: "never" }]],

  use: {
    baseURL: "http://localhost:8080",
    // Capture screenshots and traces on failure for debugging
    screenshot: "only-on-failure",
    trace: "on-first-retry",
  },

  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
  ],

  webServer: {
    // Start the Go server directly with `go run` (works on all platforms).
    // cwd is set to the project root so relative paths (./public) resolve correctly.
    command: "go run ./cmd/server",
    cwd: "../..",
    url: "http://localhost:8080",
    reuseExistingServer: !process.env.CI,
    timeout: 60_000,
    env: {
      APP_PORT: "8080",
      // Provide a placeholder key so the server starts; real OWM calls will
      // return 401/404 and the UI will show an error state — which the tests
      // assert on where relevant.
      OPENWEATHERMAP_API_KEY: process.env.OPENWEATHERMAP_API_KEY ?? "test-key",
      OPENWEATHERMAP_BASE_URL:
        process.env.OPENWEATHERMAP_BASE_URL ??
        "https://api.openweathermap.org/data/2.5",
    },
  },
});
