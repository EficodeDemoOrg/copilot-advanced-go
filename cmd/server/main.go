// @title           Weather App API
// @version         1.0
// @description     A weather service with saved locations and threshold-based alerts.
// @contact.name    GitHub Copilot Exercise
// @host            localhost:8080
// @BasePath        /

package main

import (
	"fmt"
	"log"

	"copilot-go-advanced/internal/app"
	"copilot-go-advanced/internal/config"
	"copilot-go-advanced/internal/container"

	_ "copilot-go-advanced/docs"
)

func main() {
	settings, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	c := container.New(settings)
	r := app.NewRouter(c.LocationRepo, c.WeatherSvc)

	addr := fmt.Sprintf(":%d", settings.AppPort)
	log.Printf("Starting %s on %s", settings.AppName, addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
