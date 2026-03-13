package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"copilot-go-advanced/internal/apperrors"
	"copilot-go-advanced/internal/models"
	"copilot-go-advanced/internal/repository"
	"copilot-go-advanced/internal/services"
)

// LocationHandler holds dependencies for location endpoints.
type LocationHandler struct {
	repo       repository.LocationRepository
	weatherSvc services.WeatherService
}

// NewLocationHandler creates a LocationHandler.
func NewLocationHandler(repo repository.LocationRepository, weatherSvc services.WeatherService) *LocationHandler {
	return &LocationHandler{repo: repo, weatherSvc: weatherSvc}
}

// ListLocations godoc
// @Summary      List saved locations
// @Description  Returns all saved locations sorted by creation time
// @Tags         locations
// @Produce      json
// @Success      200  {array}   models.Location
// @Router       /api/locations [get]
func (h *LocationHandler) ListLocations(c *gin.Context) {
	c.JSON(http.StatusOK, h.repo.ListAll())
}

// CreateLocation godoc
// @Summary      Save a new location
// @Description  Saves a new named location with coordinates
// @Tags         locations
// @Accept       json
// @Produce      json
// @Param        body  body      models.LocationCreate  true  "Location data"
// @Success      201   {object}  models.Location
// @Failure      422   {object}  map[string]interface{}
// @Router       /api/locations [post]
func (h *LocationHandler) CreateLocation(c *gin.Context) {
	var body models.LocationCreate
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	loc, err := h.repo.Add(body)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, loc)
}

// GetLocation godoc
// @Summary      Get a saved location
// @Description  Returns a single saved location by ID
// @Tags         locations
// @Produce      json
// @Param        id   path      string  true  "Location ID (UUID)"
// @Success      200  {object}  models.Location
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/locations/{id} [get]
func (h *LocationHandler) GetLocation(c *gin.Context) {
	id := c.Param("id")
	loc, err := h.repo.Get(id)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, loc)
}

// UpdateLocation godoc
// @Summary      Update a saved location
// @Description  Partially updates a saved location (all fields optional)
// @Tags         locations
// @Accept       json
// @Produce      json
// @Param        id    path      string               true  "Location ID (UUID)"
// @Param        body  body      models.LocationUpdate  true  "Fields to update"
// @Success      200   {object}  models.Location
// @Failure      404   {object}  map[string]interface{}
// @Failure      422   {object}  map[string]interface{}
// @Router       /api/locations/{id} [put]
func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	id := c.Param("id")

	var body models.LocationUpdate
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	loc, err := h.repo.Update(id, body)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, loc)
}

// DeleteLocation godoc
// @Summary      Delete a saved location
// @Description  Removes a saved location by ID
// @Tags         locations
// @Param        id  path  string  true  "Location ID (UUID)"
// @Success      204
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/locations/{id} [delete]
func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.Delete(id); err != nil {
		apperrors.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// GetLocationWeather godoc
// @Summary      Get weather for a saved location
// @Description  Returns current weather for the coordinates of a saved location
// @Tags         locations
// @Produce      json
// @Param        id     path      string  true   "Location ID (UUID)"
// @Param        units  query     string  false  "Temperature unit (celsius|fahrenheit|kelvin)" Enums(celsius,fahrenheit,kelvin) default(celsius)
// @Success      200    {object}  models.CurrentWeather
// @Failure      404    {object}  map[string]interface{}
// @Failure      502    {object}  map[string]interface{}
// @Failure      503    {object}  map[string]interface{}
// @Router       /api/locations/{id}/weather [get]
func (h *LocationHandler) GetLocationWeather(c *gin.Context) {
	id := c.Param("id")

	var params models.LocationWeatherQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	units := params.Units
	if units == "" {
		units = models.UnitCelsius
	}

	loc, err := h.repo.Get(id)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	weather, err := h.weatherSvc.GetCurrentWeather(loc.Coordinates.Lat, loc.Coordinates.Lon, units, loc.Name)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, weather)
}
