package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"copilot-go-advanced/internal/apperrors"
	"copilot-go-advanced/internal/models"
)

// OWMClient abstracts OpenWeatherMap HTTP calls for testability.
// All requests use units=metric; temperature conversion is the service layer's responsibility.
type OWMClient interface {
	FetchCurrentWeather(lat, lon float64) (*models.OwmCurrentWeatherResponse, error)
	FetchForecast(lat, lon float64) (*models.OwmForecastResponse, error)
}

type owmClientImpl struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

// NewOWMClient returns a production OWMClient.
func NewOWMClient(apiKey, baseURL string) OWMClient {
	return &owmClientImpl{
		apiKey:  apiKey,
		baseURL: baseURL,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *owmClientImpl) FetchCurrentWeather(lat, lon float64) (*models.OwmCurrentWeatherResponse, error) {
	endpoint := fmt.Sprintf("%s/weather", c.baseURL)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, &apperrors.WeatherAPIConnectionError{Cause: err}
	}

	q := url.Values{}
	q.Set("lat", fmt.Sprintf("%f", lat))
	q.Set("lon", fmt.Sprintf("%f", lon))
	q.Set("units", "metric")
	q.Set("appid", c.apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, &apperrors.WeatherAPIConnectionError{Cause: err}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// ok
	case http.StatusNotFound:
		return nil, &apperrors.WeatherAPINotFoundError{Lat: lat, Lon: lon}
	default:
		return nil, &apperrors.WeatherAPIError{StatusCode: resp.StatusCode}
	}

	var result models.OwmCurrentWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, &apperrors.WeatherAPIError{StatusCode: resp.StatusCode}
	}
	return &result, nil
}

func (c *owmClientImpl) FetchForecast(lat, lon float64) (*models.OwmForecastResponse, error) {
	endpoint := fmt.Sprintf("%s/forecast", c.baseURL)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, &apperrors.WeatherAPIConnectionError{Cause: err}
	}

	q := url.Values{}
	q.Set("lat", fmt.Sprintf("%f", lat))
	q.Set("lon", fmt.Sprintf("%f", lon))
	q.Set("units", "metric")
	q.Set("cnt", "40")
	q.Set("appid", c.apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, &apperrors.WeatherAPIConnectionError{Cause: err}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// ok
	case http.StatusNotFound:
		return nil, &apperrors.WeatherAPINotFoundError{Lat: lat, Lon: lon}
	default:
		return nil, &apperrors.WeatherAPIError{StatusCode: resp.StatusCode}
	}

	var result models.OwmForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, &apperrors.WeatherAPIError{StatusCode: resp.StatusCode}
	}
	return &result, nil
}
