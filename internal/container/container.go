package container

import (
	"copilot-go-advanced/internal/config"
	"copilot-go-advanced/internal/repository"
	"copilot-go-advanced/internal/services"
)

// Container holds all application-level singletons.
type Container struct {
	Settings    *config.Settings
	LocationRepo repository.LocationRepository
	OWMClient   services.OWMClient
	WeatherSvc  services.WeatherService
}

// New builds and returns a fully-wired Container from config.Settings.
func New(settings *config.Settings) *Container {
	locationRepo := repository.NewLocationRepository()
	owmClient := services.NewOWMClient(settings.OpenWeatherMapAPIKey, settings.OpenWeatherMapBaseURL)
	weatherSvc := services.NewWeatherService(owmClient, settings)

	return &Container{
		Settings:    settings,
		LocationRepo: locationRepo,
		OWMClient:   owmClient,
		WeatherSvc:  weatherSvc,
	}
}
