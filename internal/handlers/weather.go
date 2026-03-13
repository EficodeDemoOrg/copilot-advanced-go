package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"copilot-go-advanced/internal/apperrors"
	"copilot-go-advanced/internal/models"
	"copilot-go-advanced/internal/services"
)

// WeatherHandler holds dependencies for weather endpoints.
type WeatherHandler struct {
	weatherSvc services.WeatherService
}

// NewWeatherHandler creates a WeatherHandler.
func NewWeatherHandler(weatherSvc services.WeatherService) *WeatherHandler {
	return &WeatherHandler{weatherSvc: weatherSvc}
}

// GetCurrentWeather godoc
// @Summary      Get current weather
// @Description  Returns current weather conditions for the provided coordinates
// @Tags         weather
// @Accept       json
// @Produce      json
// @Param        lat    query     number  true   "Latitude (-90..90)"
// @Param        lon    query     number  true   "Longitude (-180..180)"
// @Param        units  query     string  false  "Temperature unit (celsius|fahrenheit|kelvin)" Enums(celsius,fahrenheit,kelvin) default(celsius)
// @Success      200    {object}  models.CurrentWeather
// @Failure      422    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Failure      502    {object}  map[string]interface{}
// @Failure      503    {object}  map[string]interface{}
// @Router       /api/weather/current [get]
func (h *WeatherHandler) GetCurrentWeather(c *gin.Context) {
	var params models.WeatherQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	units := params.Units
	if units == "" {
		units = models.UnitCelsius
	}

	weather, err := h.weatherSvc.GetCurrentWeather(params.Lat, params.Lon, units, "")
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, weather)
}

// GetForecast godoc
// @Summary      Get weather forecast
// @Description  Returns a multi-day daily forecast for the provided coordinates
// @Tags         weather
// @Accept       json
// @Produce      json
// @Param        lat    query     number  true   "Latitude (-90..90)"
// @Param        lon    query     number  true   "Longitude (-180..180)"
// @Param        days   query     int     false  "Number of forecast days (1..5)" default(5)
// @Param        units  query     string  false  "Temperature unit (celsius|fahrenheit|kelvin)" Enums(celsius,fahrenheit,kelvin) default(celsius)
// @Success      200    {object}  models.Forecast
// @Failure      422    {object}  map[string]interface{}
// @Failure      503    {object}  map[string]interface{}
// @Router       /api/weather/forecast [get]
func (h *WeatherHandler) GetForecast(c *gin.Context) {
	var params models.ForecastQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	units := params.Units
	if units == "" {
		units = models.UnitCelsius
	}
	days := params.Days
	if days == 0 {
		days = 5
	}

	forecast, err := h.weatherSvc.GetForecast(params.Lat, params.Lon, days, units, "")
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, forecast)
}

// GetAlerts godoc
// @Summary      Get weather alerts
// @Description  Returns threshold-based weather alerts for the provided coordinates
// @Tags         weather
// @Accept       json
// @Produce      json
// @Param        lat  query     number  true  "Latitude (-90..90)"
// @Param        lon  query     number  true  "Longitude (-180..180)"
// @Success      200  {array}   models.WeatherAlert
// @Failure      422  {object}  map[string]interface{}
// @Failure      503  {object}  map[string]interface{}
// @Router       /api/weather/alerts [get]
func (h *WeatherHandler) GetAlerts(c *gin.Context) {
	var params models.AlertQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	alerts, err := h.weatherSvc.GetAlerts(params.Lat, params.Lon)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, alerts)
}
