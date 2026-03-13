package services_test

import (
	"copilot-go-advanced/internal/models"
)

// mockOWMClient is a test double that returns configurable responses.
type mockOWMClient struct {
	currentResp *models.OwmCurrentWeatherResponse
	currentErr  error
	forecastResp *models.OwmForecastResponse
	forecastErr  error
}

func (m *mockOWMClient) FetchCurrentWeather(lat, lon float64) (*models.OwmCurrentWeatherResponse, error) {
	return m.currentResp, m.currentErr
}

func (m *mockOWMClient) FetchForecast(lat, lon float64) (*models.OwmForecastResponse, error) {
	return m.forecastResp, m.forecastErr
}
