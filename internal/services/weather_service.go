package services

import (
	"fmt"
	"math"
	"sort"
	"time"

	"copilot-go-advanced/internal/config"
	"copilot-go-advanced/internal/models"
	"copilot-go-advanced/internal/utils"
)

// WeatherService provides business logic on top of OWMClient.
type WeatherService interface {
	GetCurrentWeather(lat, lon float64, units models.TemperatureUnit, locationName string) (*models.CurrentWeather, error)
	GetForecast(lat, lon float64, days int, units models.TemperatureUnit, locationName string) (*models.Forecast, error)
	GetAlerts(lat, lon float64) ([]models.WeatherAlert, error)
}

type weatherService struct {
	client   OWMClient
	settings *config.Settings
}

// NewWeatherService returns a production WeatherService.
func NewWeatherService(client OWMClient, settings *config.Settings) WeatherService {
	return &weatherService{client: client, settings: settings}
}

// GetCurrentWeather fetches and converts current weather for the given coordinates.
// locationName overrides the city name from the API response when non-empty.
func (s *weatherService) GetCurrentWeather(lat, lon float64, units models.TemperatureUnit, locationName string) (*models.CurrentWeather, error) {
	owmResp, err := s.client.FetchCurrentWeather(lat, lon)
	if err != nil {
		return nil, err
	}

	name := owmResp.Name
	if locationName != "" {
		name = locationName
	}

	description := ""
	icon := ""
	if len(owmResp.Weather) > 0 {
		description = owmResp.Weather[0].Description
		icon = owmResp.Weather[0].Icon
	}

	temp := convertTemperature(owmResp.Main.Temp, units)
	feelsLike := convertTemperature(owmResp.Main.FeelsLike, units)

	return &models.CurrentWeather{
		Temperature:   temp,
		FeelsLike:     feelsLike,
		Humidity:      owmResp.Main.Humidity,
		Pressure:      owmResp.Main.Pressure,
		WindSpeed:     owmResp.Wind.Speed,
		WindDirection: owmResp.Wind.Deg,
		Description:   description,
		Icon:          icon,
		Timestamp:     owmResp.Dt,
		LocationName:  name,
		Units:         units,
	}, nil
}

// GetForecast fetches and aggregates the 3-hour forecast into daily summaries.
// locationName overrides the city name from the API response when non-empty.
func (s *weatherService) GetForecast(lat, lon float64, days int, units models.TemperatureUnit, locationName string) (*models.Forecast, error) {
	owmResp, err := s.client.FetchForecast(lat, lon)
	if err != nil {
		return nil, err
	}

	name := owmResp.City.Name
	if locationName != "" {
		name = locationName
	}

	// Group 3-hour intervals by date string (YYYY-MM-DD, UTC)
	type dayData struct {
		temps       []float64
		humidities  []float64
		description map[string]int
		icon        map[string]int
	}
	order := []string{}
	byDate := map[string]*dayData{}

	for _, item := range owmResp.List {
		t := time.Unix(item.Dt, 0).UTC()
		date := t.Format("2006-01-02")
		if _, exists := byDate[date]; !exists {
			order = append(order, date)
			byDate[date] = &dayData{
				description: map[string]int{},
				icon:        map[string]int{},
			}
		}
		d := byDate[date]
		d.temps = append(d.temps, item.Main.Temp)
		d.humidities = append(d.humidities, item.Main.Humidity)
		if len(item.Weather) > 0 {
			d.description[item.Weather[0].Description]++
			d.icon[item.Weather[0].Icon]++
		}
	}

	// Build ForecastDay slice limited to requested days
	if days <= 0 {
		days = 5
	}
	if days > len(order) {
		days = len(order)
	}

	forecastDays := make([]models.ForecastDay, 0, days)
	for _, date := range order[:days] {
		d := byDate[date]
		tempMin, tempMax := minMax(d.temps)
		avgHumidity := average(d.humidities)
		desc := mostFrequent(d.description)
		icon := mostFrequent(d.icon)

		forecastDays = append(forecastDays, models.ForecastDay{
			Date:        date,
			TempMin:     convertTemperature(tempMin, units),
			TempMax:     convertTemperature(tempMax, units),
			Humidity:    math.Round(avgHumidity*100) / 100,
			Description: desc,
			Icon:        icon,
		})
	}

	return &models.Forecast{
		LocationName: name,
		Units:        units,
		Days:         forecastDays,
	}, nil
}

// GetAlerts evaluates the current weather against configured thresholds.
func (s *weatherService) GetAlerts(lat, lon float64) ([]models.WeatherAlert, error) {
	owmResp, err := s.client.FetchCurrentWeather(lat, lon)
	if err != nil {
		return nil, err
	}

	alerts := []models.WeatherAlert{}
	windSpeed := owmResp.Wind.Speed
	temp := owmResp.Main.Temp
	humidity := owmResp.Main.Humidity

	// high_wind
	if windSpeed >= s.settings.AlertWindSpeedThreshold*1.5 {
		alerts = append(alerts, models.WeatherAlert{
			AlertType: "high_wind",
			Message:   fmt.Sprintf("Dangerous wind speed: %.1f m/s", windSpeed),
			Severity:  models.SeverityHigh,
			Value:     windSpeed,
			Threshold: s.settings.AlertWindSpeedThreshold,
		})
	} else if windSpeed >= s.settings.AlertWindSpeedThreshold {
		alerts = append(alerts, models.WeatherAlert{
			AlertType: "high_wind",
			Message:   fmt.Sprintf("High wind speed: %.1f m/s", windSpeed),
			Severity:  models.SeverityMedium,
			Value:     windSpeed,
			Threshold: s.settings.AlertWindSpeedThreshold,
		})
	}

	// extreme_heat
	if temp >= s.settings.AlertTempHighThreshold+5 {
		alerts = append(alerts, models.WeatherAlert{
			AlertType: "extreme_heat",
			Message:   fmt.Sprintf("Extreme heat: %.1f°C", temp),
			Severity:  models.SeverityExtreme,
			Value:     temp,
			Threshold: s.settings.AlertTempHighThreshold,
		})
	} else if temp >= s.settings.AlertTempHighThreshold {
		alerts = append(alerts, models.WeatherAlert{
			AlertType: "extreme_heat",
			Message:   fmt.Sprintf("High temperature: %.1f°C", temp),
			Severity:  models.SeverityHigh,
			Value:     temp,
			Threshold: s.settings.AlertTempHighThreshold,
		})
	}

	// extreme_cold
	if temp <= s.settings.AlertTempLowThreshold-10 {
		alerts = append(alerts, models.WeatherAlert{
			AlertType: "extreme_cold",
			Message:   fmt.Sprintf("Extreme cold: %.1f°C", temp),
			Severity:  models.SeverityExtreme,
			Value:     temp,
			Threshold: s.settings.AlertTempLowThreshold,
		})
	} else if temp <= s.settings.AlertTempLowThreshold {
		alerts = append(alerts, models.WeatherAlert{
			AlertType: "extreme_cold",
			Message:   fmt.Sprintf("Extreme cold: %.1f°C", temp),
			Severity:  models.SeverityHigh,
			Value:     temp,
			Threshold: s.settings.AlertTempLowThreshold,
		})
	}

	// high_humidity
	if humidity >= s.settings.AlertHumidityThreshold {
		alerts = append(alerts, models.WeatherAlert{
			AlertType: "high_humidity",
			Message:   fmt.Sprintf("High humidity: %.1f%%", humidity),
			Severity:  models.SeverityLow,
			Value:     humidity,
			Threshold: s.settings.AlertHumidityThreshold,
		})
	}

	return alerts, nil
}

// convertTemperature converts a Celsius value to the requested unit.
func convertTemperature(celsius float64, units models.TemperatureUnit) float64 {
	switch units {
	case models.UnitFahrenheit:
		return utils.CelsiusToFahrenheit(celsius)
	case models.UnitKelvin:
		return utils.CelsiusToKelvin(celsius)
	default:
		return math.Round(celsius*100) / 100
	}
}

func minMax(values []float64) (min, max float64) {
	if len(values) == 0 {
		return 0, 0
	}
	min, max = values[0], values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}

func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func mostFrequent(counts map[string]int) string {
	if len(counts) == 0 {
		return ""
	}
	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys) // deterministic tiebreak
	best := keys[0]
	for _, k := range keys[1:] {
		if counts[k] > counts[best] {
			best = k
		}
	}
	return best
}
