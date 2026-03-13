package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Settings holds all application configuration loaded from environment variables.
type Settings struct {
	OpenWeatherMapAPIKey  string
	OpenWeatherMapBaseURL string
	AppName               string
	AppPort               int
	Debug                 bool

	AlertWindSpeedThreshold float64
	AlertTempHighThreshold  float64
	AlertTempLowThreshold   float64
	AlertHumidityThreshold  float64
}

// Load reads configuration from the .env file (if present) and environment variables.
// Environment variables take precedence over .env file values.
func Load() (*Settings, error) {
	// Load .env file if it exists; ignore error (file is optional)
	_ = godotenv.Load()

	port, err := parseInt("APP_PORT", 8080)
	if err != nil {
		return nil, fmt.Errorf("invalid APP_PORT: %w", err)
	}

	debug, err := parseBool("DEBUG", false)
	if err != nil {
		return nil, fmt.Errorf("invalid DEBUG: %w", err)
	}

	windThreshold, err := parseFloat("ALERT_WIND_SPEED_THRESHOLD", 20.0)
	if err != nil {
		return nil, fmt.Errorf("invalid ALERT_WIND_SPEED_THRESHOLD: %w", err)
	}

	tempHighThreshold, err := parseFloat("ALERT_TEMP_HIGH_THRESHOLD", 40.0)
	if err != nil {
		return nil, fmt.Errorf("invalid ALERT_TEMP_HIGH_THRESHOLD: %w", err)
	}

	tempLowThreshold, err := parseFloat("ALERT_TEMP_LOW_THRESHOLD", -20.0)
	if err != nil {
		return nil, fmt.Errorf("invalid ALERT_TEMP_LOW_THRESHOLD: %w", err)
	}

	humidityThreshold, err := parseFloat("ALERT_HUMIDITY_THRESHOLD", 90.0)
	if err != nil {
		return nil, fmt.Errorf("invalid ALERT_HUMIDITY_THRESHOLD: %w", err)
	}

	return &Settings{
		OpenWeatherMapAPIKey:    getEnv("OPENWEATHERMAP_API_KEY", ""),
		OpenWeatherMapBaseURL:   getEnv("OPENWEATHERMAP_BASE_URL", "https://api.openweathermap.org/data/2.5"),
		AppName:                 getEnv("APP_NAME", "Weather App"),
		AppPort:                 port,
		Debug:                   debug,
		AlertWindSpeedThreshold: windThreshold,
		AlertTempHighThreshold:  tempHighThreshold,
		AlertTempLowThreshold:   tempLowThreshold,
		AlertHumidityThreshold:  humidityThreshold,
	}, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func parseInt(key string, defaultVal int) (int, error) {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal, nil
	}
	return strconv.Atoi(val)
}

func parseBool(key string, defaultVal bool) (bool, error) {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal, nil
	}
	return strconv.ParseBool(val)
}

func parseFloat(key string, defaultVal float64) (float64, error) {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal, nil
	}
	return strconv.ParseFloat(val, 64)
}
