package apperrors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WeatherAppError is the base domain error type.
type WeatherAppError struct {
	Message string
}

func (e *WeatherAppError) Error() string {
	return e.Message
}

// WeatherAPIError represents a non-200 response from the OWM API.
type WeatherAPIError struct {
	StatusCode int
	Message    string
}

func (e *WeatherAPIError) Error() string {
	return e.Message
}

// WeatherAPINotFoundError represents a 404 from the OWM API.
type WeatherAPINotFoundError struct {
	Lat float64
	Lon float64
}

func (e *WeatherAPINotFoundError) Error() string {
	return fmt.Sprintf("weather data not found for coordinates (%.4f, %.4f)", e.Lat, e.Lon)
}

// WeatherAPIConnectionError represents a network/timeout error calling OWM API.
type WeatherAPIConnectionError struct {
	Cause error
}

func (e *WeatherAPIConnectionError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("weather API connection error: %s", e.Cause.Error())
	}
	return "weather API connection error"
}

// LocationNotFoundError is returned when a saved location ID does not exist.
type LocationNotFoundError struct {
	ID string
}

func (e *LocationNotFoundError) Error() string {
	return fmt.Sprintf("location not found: %s", e.ID)
}

// HandleError maps domain errors to HTTP responses and writes the JSON to the context.
// This must be called instead of c.JSON in handlers.
func HandleError(c *gin.Context, err error) {
	var locNotFound *LocationNotFoundError
	if errors.As(err, &locNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Location not found: %s", locNotFound.ID)})
		return
	}

	var apiNotFound *WeatherAPINotFoundError
	if errors.As(err, &apiNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Weather data not found for the given location"})
		return
	}

	var connErr *WeatherAPIConnectionError
	if errors.As(err, &connErr) {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": connErr.Error()})
		return
	}

	var apiErr *WeatherAPIError
	if errors.As(err, &apiErr) {
		c.JSON(http.StatusBadGateway, gin.H{"error": apiErr.Message})
		return
	}

	var appErr *WeatherAppError
	if errors.As(err, &appErr) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}
