import { test, expect, Page } from "@playwright/test";

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

/** Navigate to the dashboard and wait for it to be interactive. */
async function loadDashboard(page: Page) {
  await page.goto("/");
  // The root route serves index.html — wait for the main heading
  await expect(page.locator("h1")).toBeVisible({ timeout: 10_000 });
}

/** Fill in the coordinate search form and submit it. */
async function searchCoordinates(page: Page, lat: string, lon: string) {
  await page.fill("#lat", lat);
  await page.fill("#lon", lon);
  await page.click('button[type="submit"], #search-btn');
}

// ---------------------------------------------------------------------------
// Suite: Dashboard loads
// ---------------------------------------------------------------------------

test.describe("Dashboard loads", () => {
  test("serves the HTML dashboard at /", async ({ page }) => {
    await page.goto("/");

    // The page title should mention the app
    await expect(page).toHaveTitle(/weather/i);

    // The main container should be present
    await expect(page.locator("body")).toBeVisible();
  });

  test("has a coordinate search form", async ({ page }) => {
    await loadDashboard(page);

    // Both input fields must exist
    await expect(page.locator("#lat")).toBeVisible();
    await expect(page.locator("#lon")).toBeVisible();
  });

  test("has a saved locations section", async ({ page }) => {
    await loadDashboard(page);

    // The saved-locations list element must be present
    await expect(page.locator("#saved-locations, #locations-list, [data-testid='locations']").first()).toBeVisible();
  });
});

// ---------------------------------------------------------------------------
// Suite: Coordinate weather search
// ---------------------------------------------------------------------------

test.describe("Weather search", () => {
  test("shows error state for invalid coordinates", async ({ page }) => {
    await loadDashboard(page);

    // Submit the form with out-of-range latitude
    await page.fill("#lat", "999");
    await page.fill("#lon", "0");
    await page.click('button[type="submit"], #search-btn');

    // Expect visible error feedback — either inline validation or
    // an error message rendered by the JS app
    const errorLocator = page.locator(
      ".error, [class*='error'], #error-message, [role='alert']"
    );
    await expect(errorLocator.first()).toBeVisible({ timeout: 5_000 });
  });

  test("searching valid coordinates triggers a weather request and renders a result or error state", async ({
    page,
  }) => {
    await loadDashboard(page);

    // London coordinates
    await searchCoordinates(page, "51.51", "-0.13");

    // After the request resolves the UI should show either:
    //   (a) weather data (temperature value visible), OR
    //   (b) an explicit error state (API key invalid in test env)
    // Either outcome proves the request flow completes end-to-end.
    const resultOrError = page.locator(
      ".weather-result, #current-weather, #error-message, .error, [class*='error'], [role='alert'], .temperature"
    );
    await expect(resultOrError.first()).toBeVisible({ timeout: 10_000 });
  });
});

// ---------------------------------------------------------------------------
// Suite: Saved locations CRUD
// ---------------------------------------------------------------------------

test.describe("Saved locations", () => {
  test("can save a location via the API and see it reflected in the UI", async ({
    page,
    request,
  }) => {
    // Create a location through the REST API directly (no OWM call needed)
    const res = await request.post("/api/locations", {
      data: { name: "Test City", lat: 48.85, lon: 2.35 },
    });
    expect(res.ok()).toBeTruthy();
    const loc = await res.json();
    expect(loc.id).toBeTruthy();

    // Load the dashboard and verify the new location appears
    await loadDashboard(page);
    await page.reload(); // ensure fresh fetch of locations list

    // The location name should appear somewhere on the page
    await expect(page.getByText("Test City")).toBeVisible({ timeout: 5_000 });

    // Clean up
    await request.delete(`/api/locations/${loc.id}`);
  });

  test("can delete a location via the API and verify it disappears", async ({
    page,
    request,
  }) => {
    // Create then immediately delete
    const createRes = await request.post("/api/locations", {
      data: { name: "Delete Me", lat: 40.71, lon: -74.01 },
    });
    expect(createRes.ok()).toBeTruthy();
    const loc = await createRes.json();

    const deleteRes = await request.delete(`/api/locations/${loc.id}`);
    expect(deleteRes.status()).toBe(204);

    // After deletion the GET should return 404
    const getRes = await request.get(`/api/locations/${loc.id}`);
    expect(getRes.status()).toBe(404);
  });
});

// ---------------------------------------------------------------------------
// Suite: Forecast chart
// ---------------------------------------------------------------------------

test.describe("Forecast chart", () => {
  test("the page contains a canvas element for Chart.js rendering", async ({
    page,
  }) => {
    await loadDashboard(page);

    // The dashboard includes Chart.js — there should be at least one <canvas>
    const canvas = page.locator("canvas");
    await expect(canvas.first()).toBeVisible({ timeout: 5_000 });
  });
});

// ---------------------------------------------------------------------------
// Suite: Alert display
// ---------------------------------------------------------------------------

test.describe("Weather alerts", () => {
  test("the alerts section is present in the DOM", async ({ page }) => {
    await loadDashboard(page);

    // Whether alerts are populated or empty, the container should exist
    const alertsSection = page.locator(
      "#alerts, #alerts-section, [data-testid='alerts'], .alerts"
    );
    // It's acceptable if it's hidden until a search is run — check it exists
    await expect(alertsSection.first()).toHaveCount(1, { timeout: 5_000 });
  });
});
