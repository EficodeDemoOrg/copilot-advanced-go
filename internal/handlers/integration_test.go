// Package handlers_test contains HTTP-level integration tests.
// A real WeatherService is used; only the OWM HTTP client is mocked.
package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"copilot-go-advanced/internal/apperrors"
	"copilot-go-advanced/internal/models"
	"copilot-go-advanced/internal/repository"
	"copilot-go-advanced/internal/services"
	"copilot-go-advanced/internal/testhelpers"
)

// mockOWMClient implements services.OWMClient for integration tests.
// All weather-service business logic uses the real implementation; only OWM HTTP calls are mocked.
type mockOWMClient struct {
	currentResp  *models.OwmCurrentWeatherResponse
	currentErr   error
	forecastResp *models.OwmForecastResponse
	forecastErr  error
}

func (m *mockOWMClient) FetchCurrentWeather(_, _ float64) (*models.OwmCurrentWeatherResponse, error) {
	return m.currentResp, m.currentErr
}

func (m *mockOWMClient) FetchForecast(_, _ float64) (*models.OwmForecastResponse, error) {
	return m.forecastResp, m.forecastErr
}

// newRouter creates an isolated Gin engine per test — fresh repo + real service + mock OWM.
func newRouter(
	currentResp *models.OwmCurrentWeatherResponse,
	currentErr error,
	forecastResp *models.OwmForecastResponse,
	forecastErr error,
) *gin.Engine {
	mock := &mockOWMClient{
		currentResp:  currentResp,
		currentErr:   currentErr,
		forecastResp: forecastResp,
		forecastErr:  forecastErr,
	}
	repo := repository.NewLocationRepository()
	settings := testhelpers.TestSettings()
	weatherSvc := services.NewWeatherService(mock, settings)
	return testhelpers.NewTestApp(repo, weatherSvc)
}

func strp(s string) *string { return &s }

// ── GET /api/weather/current ──────────────────────────────────────────────────

func TestWeatherCurrent_200ValidCoords(t *testing.T) {
	owmResp := testhelpers.MakeOwmCurrentWeatherResponse()
	r := newRouter(&owmResp, nil, nil, nil)

	w := testhelpers.PerformRequest(r, http.MethodGet, "/api/weather/current?lat=51.51&lon=-0.13", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var result models.CurrentWeather
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, "London", result.LocationName)
	assert.Equal(t, models.UnitCelsius, result.Units)
}

func TestWeatherCurrent_200FahrenheitConversion(t *testing.T) {
	owmResp := testhelpers.MakeOwmCurrentWeatherResponse() // 15°C
	r := newRouter(&owmResp, nil, nil, nil)

	w := testhelpers.PerformRequest(r, http.MethodGet, "/api/weather/current?lat=51.51&lon=-0.13&units=fahrenheit", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var result models.CurrentWeather
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, models.UnitFahrenheit, result.Units)
	assert.InDelta(t, 59.0, result.Temperature, 0.1) // 15°C → 59°F
}

func TestWeatherCurrent_422MissingParams(t *testing.T) {
	r := newRouter(nil, nil, nil, nil)
	w := testhelpers.PerformRequest(r, http.MethodGet, "/api/weather/current", nil)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestWeatherCurrent_422InvalidLat(t *testing.T) {
	r := newRouter(nil, nil, nil, nil)
	w := testhelpers.PerformRequest(r, http.MethodGet, "/api/weather/current?lat=999&lon=-0.13", nil)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestWeatherCurrent_404WhenOWMReturns404(t *testing.T) {
	notFound := &apperrors.WeatherAPINotFoundError{Lat: 51.51, Lon: -0.13}
	r := newRouter(nil, notFound, nil, nil)

	w := testhelpers.PerformRequest(r, http.MethodGet, "/api/weather/current?lat=51.51&lon=-0.13", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestWeatherCurrent_502OnOWMServerError(t *testing.T) {
	apiErr := &apperrors.WeatherAPIError{StatusCode: 500, Message: "upstream error"}
	r := newRouter(nil, apiErr, nil, nil)

	w := testhelpers.PerformRequest(r, http.MethodGet, "/api/weather/current?lat=51.51&lon=-0.13", nil)
	assert.Equal(t, http.StatusBadGateway, w.Code)
}

// ── GET /api/weather/forecast ─────────────────────────────────────────────────

func TestWeatherForecast_200WithData(t *testing.T) {
	owmForecast := testhelpers.MakeOwmForecastResponse()
	r := newRouter(nil, nil, &owmForecast, nil)

	w := testhelpers.PerformRequest(r, http.MethodGet, "/api/weather/forecast?lat=51.51&lon=-0.13", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var result models.Forecast
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.NotEmpty(t, result.Days)
}

func TestWeatherForecast_DaysParamLimitsResults(t *testing.T) {
	// 5 items, each on a distinct future date — one item per day
	items := make([]models.OwmForecastItem, 5)
	for i := range items {
		items[i] = testhelpers.MakeOwmForecastItem(i + 1)
	}
	owmForecast := models.OwmForecastResponse{List: items, City: models.OwmCity{Name: "London"}}
	r := newRouter(nil, nil, &owmForecast, nil)

	w := testhelpers.PerformRequest(r, http.MethodGet, "/api/weather/forecast?lat=51.51&lon=-0.13&days=2", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var result models.Forecast
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Len(t, result.Days, 2)
}

func TestWeatherForecast_422DaysExceedsMax(t *testing.T) {
	// days=6 is non-zero, so omitempty does not skip the max=5 validator
	r := newRouter(nil, nil, nil, nil)
	w := testhelpers.PerformRequest(r, http.MethodGet, "/api/weather/forecast?lat=51.51&lon=-0.13&days=6", nil)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestWeatherForecast_KelvinConversion(t *testing.T) {
	owmForecast := testhelpers.MakeOwmForecastResponse()
	r := newRouter(nil, nil, &owmForecast, nil)

	w := testhelpers.PerformRequest(r, http.MethodGet, "/api/weather/forecast?lat=51.51&lon=-0.13&units=kelvin", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var result models.Forecast
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, models.UnitKelvin, result.Units)
	require.NotEmpty(t, result.Days)
	assert.Greater(t, result.Days[0].TempMin, 270.0) // Celsius + 273.15
}

// ── GET /api/weather/alerts ───────────────────────────────────────────────────

func TestWeatherAlerts_200EmptyWhenBelowThresholds(t *testing.T) {
	// Default factory: 15°C, 5.5 m/s, 72% — all below thresholds in TestSettings
	owmResp := testhelpers.MakeOwmCurrentWeatherResponse()
	r := newRouter(&owmResp, nil, nil, nil)

	w := testhelpers.PerformRequest(r, http.MethodGet, "/api/weather/alerts?lat=51.51&lon=-0.13", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var result []models.WeatherAlert
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Empty(t, result)
}

func TestWeatherAlerts_200MultipleAlerts(t *testing.T) {
	// wind=25 > 20 threshold → high_wind; humidity=95 > 90 threshold → high_humidity
	owmResp := testhelpers.MakeOwmCurrentWeatherResponse(func(resp *models.OwmCurrentWeatherResponse) {
		resp.Wind.Speed = 25.0
		resp.Main.Humidity = 95.0
	})
	r := newRouter(&owmResp, nil, nil, nil)

	w := testhelpers.PerformRequest(r, http.MethodGet, "/api/weather/alerts?lat=51.51&lon=-0.13", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var result []models.WeatherAlert
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Len(t, result, 2)
}

// ── POST /api/locations ───────────────────────────────────────────────────────

func TestLocations_201ValidCreation(t *testing.T) {
	r := newRouter(nil, nil, nil, nil)

	w := testhelpers.PerformRequest(r, http.MethodPost, "/api/locations", strp(`{"name":"London","lat":51.51,"lon":-0.13}`))

	assert.Equal(t, http.StatusCreated, w.Code)
	var result models.Location
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, "London", result.Name)
}

func TestLocations_422InvalidLat(t *testing.T) {
	r := newRouter(nil, nil, nil, nil)
	w := testhelpers.PerformRequest(r, http.MethodPost, "/api/locations", strp(`{"name":"London","lat":999,"lon":-0.13}`))
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestLocations_422EmptyName(t *testing.T) {
	r := newRouter(nil, nil, nil, nil)
	w := testhelpers.PerformRequest(r, http.MethodPost, "/api/locations", strp(`{"name":"","lat":51.51,"lon":-0.13}`))
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestLocations_MultipleCreatesHaveUniqueIDs(t *testing.T) {
	r := newRouter(nil, nil, nil, nil)
	body := strp(`{"name":"London","lat":51.51,"lon":-0.13}`)

	var ids [3]string
	for i := range ids {
		w := testhelpers.PerformRequest(r, http.MethodPost, "/api/locations", body)
		require.Equal(t, http.StatusCreated, w.Code)
		var loc models.Location
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &loc))
		ids[i] = loc.ID
	}
	assert.NotEqual(t, ids[0], ids[1])
	assert.NotEqual(t, ids[1], ids[2])
}

// ── GET /api/locations ────────────────────────────────────────────────────────

func TestLocations_200EmptyList(t *testing.T) {
	r := newRouter(nil, nil, nil, nil)
	w := testhelpers.PerformRequest(r, http.MethodGet, "/api/locations", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var result []models.Location
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Empty(t, result)
}

func TestLocations_200ListAfterCreate(t *testing.T) {
	r := newRouter(nil, nil, nil, nil)
	testhelpers.PerformRequest(r, http.MethodPost, "/api/locations", strp(`{"name":"London","lat":51.51,"lon":-0.13}`))

	w := testhelpers.PerformRequest(r, http.MethodGet, "/api/locations", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var result []models.Location
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Len(t, result, 1)
}

// ── GET /api/locations/:id ────────────────────────────────────────────────────

func TestLocations_200GetExisting(t *testing.T) {
	r := newRouter(nil, nil, nil, nil)
	created := createLocation(t, r, `{"name":"London","lat":51.51,"lon":-0.13}`)

	w := testhelpers.PerformRequest(r, http.MethodGet, fmt.Sprintf("/api/locations/%s", created.ID), nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLocations_404GetNonExistent(t *testing.T) {
	r := newRouter(nil, nil, nil, nil)
	w := testhelpers.PerformRequest(r, http.MethodGet, "/api/locations/00000000-0000-0000-0000-000000000000", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ── PUT /api/locations/:id ────────────────────────────────────────────────────

func TestLocations_200UpdateName(t *testing.T) {
	r := newRouter(nil, nil, nil, nil)
	created := createLocation(t, r, `{"name":"London","lat":51.51,"lon":-0.13}`)

	w := testhelpers.PerformRequest(r, http.MethodPut, fmt.Sprintf("/api/locations/%s", created.ID), strp(`{"name":"New York"}`))

	assert.Equal(t, http.StatusOK, w.Code)
	var updated models.Location
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &updated))
	assert.Equal(t, "New York", updated.Name)
}

func TestLocations_404UpdateNonExistent(t *testing.T) {
	r := newRouter(nil, nil, nil, nil)
	w := testhelpers.PerformRequest(r, http.MethodPut, "/api/locations/missing", strp(`{"name":"x"}`))
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ── DELETE /api/locations/:id ─────────────────────────────────────────────────

func TestLocations_204DeleteExisting(t *testing.T) {
	r := newRouter(nil, nil, nil, nil)
	created := createLocation(t, r, `{"name":"London","lat":51.51,"lon":-0.13}`)

	w := testhelpers.PerformRequest(r, http.MethodDelete, fmt.Sprintf("/api/locations/%s", created.ID), nil)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Confirm deletion
	w2 := testhelpers.PerformRequest(r, http.MethodGet, fmt.Sprintf("/api/locations/%s", created.ID), nil)
	assert.Equal(t, http.StatusNotFound, w2.Code)
}

func TestLocations_404DeleteNonExistent(t *testing.T) {
	r := newRouter(nil, nil, nil, nil)
	w := testhelpers.PerformRequest(r, http.MethodDelete, "/api/locations/missing", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ── GET /api/locations/:id/weather ────────────────────────────────────────────

func TestLocations_200WeatherUsesLocationName(t *testing.T) {
	owmResp := testhelpers.MakeOwmCurrentWeatherResponse()
	r := newRouter(&owmResp, nil, nil, nil)
	created := createLocation(t, r, `{"name":"My City","lat":51.51,"lon":-0.13}`)

	w := testhelpers.PerformRequest(r, http.MethodGet, fmt.Sprintf("/api/locations/%s/weather", created.ID), nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var weather models.CurrentWeather
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &weather))
	// The handler passes the saved location name, which overrides the OWM city name.
	assert.Equal(t, "My City", weather.LocationName)
}

func TestLocations_404WeatherNonExistentLocation(t *testing.T) {
	r := newRouter(nil, nil, nil, nil)
	w := testhelpers.PerformRequest(r, http.MethodGet, "/api/locations/00000000-0000-0000-0000-000000000000/weather", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ── shared test helpers ───────────────────────────────────────────────────────

// createLocation POSTs a new location and returns the decoded result.
func createLocation(t *testing.T, r *gin.Engine, body string) models.Location {
	t.Helper()
	w := testhelpers.PerformRequest(r, http.MethodPost, "/api/locations", strp(body))
	require.Equal(t, http.StatusCreated, w.Code, "createLocation: unexpected status %d", w.Code)
	var loc models.Location
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &loc))
	return loc
}
