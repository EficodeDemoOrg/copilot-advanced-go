package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"copilot-go-advanced/internal/apperrors"
	"copilot-go-advanced/internal/models"
	"copilot-go-advanced/internal/services"
	"copilot-go-advanced/internal/testhelpers"
)

func newWeatherSvc(mock *mockOWMClient) services.WeatherService {
	return services.NewWeatherService(mock, testhelpers.TestSettings())
}

// ── GetCurrentWeather ─────────────────────────────────────────────────────────

func TestGetCurrentWeather_CelsiusPassthrough(t *testing.T) {
	owmResp := testhelpers.MakeOwmCurrentWeatherResponse()
	svc := newWeatherSvc(&mockOWMClient{currentResp: &owmResp})

	result, err := svc.GetCurrentWeather(51.51, -0.13, models.UnitCelsius, "")

	require.NoError(t, err)
	assert.InDelta(t, 15.0, result.Temperature, 0.01)
	assert.Equal(t, models.UnitCelsius, result.Units)
}

func TestGetCurrentWeather_FahrenheitConversion(t *testing.T) {
	owmResp := testhelpers.MakeOwmCurrentWeatherResponse()
	svc := newWeatherSvc(&mockOWMClient{currentResp: &owmResp})

	result, err := svc.GetCurrentWeather(51.51, -0.13, models.UnitFahrenheit, "")

	require.NoError(t, err)
	assert.InDelta(t, 59.0, result.Temperature, 0.01) // 15°C → 59°F
	assert.Equal(t, models.UnitFahrenheit, result.Units)
}

func TestGetCurrentWeather_KelvinConversion(t *testing.T) {
	owmResp := testhelpers.MakeOwmCurrentWeatherResponse()
	svc := newWeatherSvc(&mockOWMClient{currentResp: &owmResp})

	result, err := svc.GetCurrentWeather(51.51, -0.13, models.UnitKelvin, "")

	require.NoError(t, err)
	assert.InDelta(t, 288.15, result.Temperature, 0.01) // 15°C → 288.15K
}

func TestGetCurrentWeather_NonTemperatureFieldsUnchanged(t *testing.T) {
	owmResp := testhelpers.MakeOwmCurrentWeatherResponse()
	svc := newWeatherSvc(&mockOWMClient{currentResp: &owmResp})

	result, err := svc.GetCurrentWeather(51.51, -0.13, models.UnitFahrenheit, "")

	require.NoError(t, err)
	assert.Equal(t, 72.0, result.Humidity)
	assert.Equal(t, 1013.0, result.Pressure)
	assert.Equal(t, 5.5, result.WindSpeed)
	assert.Equal(t, 270, result.WindDirection)
}

func TestGetCurrentWeather_UsesOWMCityName(t *testing.T) {
	owmResp := testhelpers.MakeOwmCurrentWeatherResponse()
	svc := newWeatherSvc(&mockOWMClient{currentResp: &owmResp})

	result, err := svc.GetCurrentWeather(51.51, -0.13, models.UnitCelsius, "")

	require.NoError(t, err)
	assert.Equal(t, "London", result.LocationName)
}

func TestGetCurrentWeather_OverridesNameWhenProvided(t *testing.T) {
	owmResp := testhelpers.MakeOwmCurrentWeatherResponse()
	svc := newWeatherSvc(&mockOWMClient{currentResp: &owmResp})

	result, err := svc.GetCurrentWeather(51.51, -0.13, models.UnitCelsius, "My London")

	require.NoError(t, err)
	assert.Equal(t, "My London", result.LocationName)
}

func TestGetCurrentWeather_PropagatesAPIError(t *testing.T) {
	apiErr := &apperrors.WeatherAPIConnectionError{Cause: nil}
	svc := newWeatherSvc(&mockOWMClient{currentErr: apiErr})

	_, err := svc.GetCurrentWeather(51.51, -0.13, models.UnitCelsius, "")

	assert.ErrorAs(t, err, &apiErr)
}

// ── GetForecast ───────────────────────────────────────────────────────────────

func TestGetForecast_AggregatesDailyFromThreeHour(t *testing.T) {
	owmResp := testhelpers.MakeOwmForecastResponse()
	svc := newWeatherSvc(&mockOWMClient{forecastResp: &owmResp})

	result, err := svc.GetForecast(51.51, -0.13, 5, models.UnitCelsius, "")

	require.NoError(t, err)
	// MakeOwmForecastResponse builds 3 days of intervals
	assert.LessOrEqual(t, len(result.Days), 3)
	assert.Greater(t, len(result.Days), 0)
}

func TestGetForecast_DaysParamLimitsResults(t *testing.T) {
	// Build a forecast with 5 different days
	items := make([]models.OwmForecastItem, 5)
	for i := range items {
		items[i] = testhelpers.MakeOwmForecastItem(i + 1)
	}
	owmResp := models.OwmForecastResponse{
		List: items,
		City: models.OwmCity{Name: "London"},
	}
	svc := newWeatherSvc(&mockOWMClient{forecastResp: &owmResp})

	result, err := svc.GetForecast(51.51, -0.13, 3, models.UnitCelsius, "")

	require.NoError(t, err)
	assert.Len(t, result.Days, 3)
}

func TestGetForecast_KelvinConversion(t *testing.T) {
	owmResp := testhelpers.MakeOwmForecastResponse()
	svc := newWeatherSvc(&mockOWMClient{forecastResp: &owmResp})

	result, err := svc.GetForecast(51.51, -0.13, 5, models.UnitKelvin, "")

	require.NoError(t, err)
	require.NotEmpty(t, result.Days)
	// Day 1 base temp is 13°C → 286.15K
	assert.Greater(t, result.Days[0].TempMin, 270.0) // well above 0K means conversion happened
	assert.Equal(t, models.UnitKelvin, result.Units)
}

func TestGetForecast_LocationNameOverride(t *testing.T) {
	owmResp := testhelpers.MakeOwmForecastResponse()
	svc := newWeatherSvc(&mockOWMClient{forecastResp: &owmResp})

	result, err := svc.GetForecast(51.51, -0.13, 5, models.UnitCelsius, "Custom Name")

	require.NoError(t, err)
	assert.Equal(t, "Custom Name", result.LocationName)
}

// ── GetAlerts ─────────────────────────────────────────────────────────────────

func TestGetAlerts_NoAlertsWhenBelowThresholds(t *testing.T) {
	owmResp := testhelpers.MakeOwmCurrentWeatherResponse(func(r *models.OwmCurrentWeatherResponse) {
		r.Main.Temp = 20.0   // below 40 threshold
		r.Wind.Speed = 5.0   // below 20 threshold
		r.Main.Humidity = 50 // below 90 threshold
	})
	svc := newWeatherSvc(&mockOWMClient{currentResp: &owmResp})

	alerts, err := svc.GetAlerts(51.51, -0.13)

	require.NoError(t, err)
	assert.Empty(t, alerts)
}

func TestGetAlerts_HighWindMediumSeverity(t *testing.T) {
	owmResp := testhelpers.MakeOwmCurrentWeatherResponse(func(r *models.OwmCurrentWeatherResponse) {
		r.Wind.Speed = 22.0 // ≥20 but <30 (1.5×20)
	})
	svc := newWeatherSvc(&mockOWMClient{currentResp: &owmResp})

	alerts, err := svc.GetAlerts(51.51, -0.13)

	require.NoError(t, err)
	require.Len(t, alerts, 1)
	assert.Equal(t, "high_wind", alerts[0].AlertType)
	assert.Equal(t, models.SeverityMedium, alerts[0].Severity)
}

func TestGetAlerts_HighWindHighSeverity(t *testing.T) {
	owmResp := testhelpers.MakeOwmCurrentWeatherResponse(func(r *models.OwmCurrentWeatherResponse) {
		r.Wind.Speed = 31.0 // ≥1.5×20=30
	})
	svc := newWeatherSvc(&mockOWMClient{currentResp: &owmResp})

	alerts, err := svc.GetAlerts(51.51, -0.13)

	require.NoError(t, err)
	require.Len(t, alerts, 1)
	assert.Equal(t, models.SeverityHigh, alerts[0].Severity)
}

func TestGetAlerts_ExtremeHeat(t *testing.T) {
	owmResp := testhelpers.MakeOwmCurrentWeatherResponse(func(r *models.OwmCurrentWeatherResponse) {
		r.Main.Temp = 46.0 // ≥40+5=45 → extreme
	})
	svc := newWeatherSvc(&mockOWMClient{currentResp: &owmResp})

	alerts, err := svc.GetAlerts(51.51, -0.13)

	require.NoError(t, err)
	require.Len(t, alerts, 1)
	assert.Equal(t, "extreme_heat", alerts[0].AlertType)
	assert.Equal(t, models.SeverityExtreme, alerts[0].Severity)
}

func TestGetAlerts_HighTemp(t *testing.T) {
	owmResp := testhelpers.MakeOwmCurrentWeatherResponse(func(r *models.OwmCurrentWeatherResponse) {
		r.Main.Temp = 41.0 // ≥40 but <45
	})
	svc := newWeatherSvc(&mockOWMClient{currentResp: &owmResp})

	alerts, err := svc.GetAlerts(51.51, -0.13)

	require.NoError(t, err)
	require.Len(t, alerts, 1)
	assert.Equal(t, "extreme_heat", alerts[0].AlertType)
	assert.Equal(t, models.SeverityHigh, alerts[0].Severity)
}

func TestGetAlerts_ExtremeCold(t *testing.T) {
	owmResp := testhelpers.MakeOwmCurrentWeatherResponse(func(r *models.OwmCurrentWeatherResponse) {
		r.Main.Temp = -32.0 // ≤-20-10=-30 → extreme
	})
	svc := newWeatherSvc(&mockOWMClient{currentResp: &owmResp})

	alerts, err := svc.GetAlerts(51.51, -0.13)

	require.NoError(t, err)
	require.Len(t, alerts, 1)
	assert.Equal(t, "extreme_cold", alerts[0].AlertType)
	assert.Equal(t, models.SeverityExtreme, alerts[0].Severity)
}

func TestGetAlerts_HighHumidity(t *testing.T) {
	owmResp := testhelpers.MakeOwmCurrentWeatherResponse(func(r *models.OwmCurrentWeatherResponse) {
		r.Main.Humidity = 92.0 // ≥90
	})
	svc := newWeatherSvc(&mockOWMClient{currentResp: &owmResp})

	alerts, err := svc.GetAlerts(51.51, -0.13)

	require.NoError(t, err)
	require.Len(t, alerts, 1)
	assert.Equal(t, "high_humidity", alerts[0].AlertType)
	assert.Equal(t, models.SeverityLow, alerts[0].Severity)
}

func TestGetAlerts_MultipleThresholdsProduceMultipleAlerts(t *testing.T) {
	owmResp := testhelpers.MakeOwmCurrentWeatherResponse(func(r *models.OwmCurrentWeatherResponse) {
		r.Wind.Speed = 25.0  // high_wind medium
		r.Main.Humidity = 95.0 // high_humidity low
	})
	svc := newWeatherSvc(&mockOWMClient{currentResp: &owmResp})

	alerts, err := svc.GetAlerts(51.51, -0.13)

	require.NoError(t, err)
	assert.Len(t, alerts, 2)

	types := make(map[string]bool)
	for _, a := range alerts {
		types[a.AlertType] = true
	}
	assert.True(t, types["high_wind"])
	assert.True(t, types["high_humidity"])
}

// guard against time.Now().Unix() being int64 (OWM Dt type)
var _ int64 = time.Now().Unix()
