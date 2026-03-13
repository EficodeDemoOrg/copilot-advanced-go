package testhelpers

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"

	"copilot-go-advanced/internal/app"
	"copilot-go-advanced/internal/config"
	"copilot-go-advanced/internal/repository"
	"copilot-go-advanced/internal/services"
)

// TestSettings returns a Settings struct with safe defaults for tests.
func TestSettings() *config.Settings {
	return &config.Settings{
		OpenWeatherMapAPIKey:    "test-key",
		OpenWeatherMapBaseURL:   "https://api.openweathermap.org/data/2.5",
		AppName:                 "Weather App (test)",
		AppPort:                 8080,
		Debug:                   false,
		AlertWindSpeedThreshold: 20.0,
		AlertTempHighThreshold:  40.0,
		AlertTempLowThreshold:   -20.0,
		AlertHumidityThreshold:  90.0,
	}
}

// NewTestApp builds a Gin engine using the provided repo and weather service.
// Use this in integration tests to inject a mock OWM client.
func NewTestApp(repo repository.LocationRepository, weatherSvc services.WeatherService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return app.NewRouter(repo, weatherSvc)
}

// PerformRequest executes an HTTP request against the provided Gin engine and
// returns the recorded response. The caller is responsible for reading the body.
func PerformRequest(r *gin.Engine, method, path string, body *string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, strings.NewReader(*body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	r.ServeHTTP(w, req)
	return w
}
