package app

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"copilot-go-advanced/internal/handlers"
	"copilot-go-advanced/internal/repository"
	"copilot-go-advanced/internal/services"
)

// NewRouter builds and returns a configured *gin.Engine.
// It wires all routes, static file serving, and Swagger UI.
func NewRouter(
	locationRepo repository.LocationRepository,
	weatherSvc services.WeatherService,
) *gin.Engine {
	r := gin.Default()

	// Static frontend
	r.Static("/static", "./public")
	r.StaticFile("/", "./public/index.html")

	// Swagger UI
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/api")
	{
		weatherHandler := handlers.NewWeatherHandler(weatherSvc)
		weather := api.Group("/weather")
		{
			weather.GET("/current", weatherHandler.GetCurrentWeather)
			weather.GET("/forecast", weatherHandler.GetForecast)
			weather.GET("/alerts", weatherHandler.GetAlerts)
		}

		locHandler := handlers.NewLocationHandler(locationRepo, weatherSvc)
		locations := api.Group("/locations")
		{
			locations.GET("", locHandler.ListLocations)
			locations.POST("", locHandler.CreateLocation)
			locations.GET("/:id", locHandler.GetLocation)
			locations.PUT("/:id", locHandler.UpdateLocation)
			locations.DELETE("/:id", locHandler.DeleteLocation)
			locations.GET("/:id/weather", locHandler.GetLocationWeather)
		}
	}

	return r
}
