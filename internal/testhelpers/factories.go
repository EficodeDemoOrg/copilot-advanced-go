package testhelpers

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"copilot-go-advanced/internal/models"
)

// MakeCoordinates returns default London coordinates (51.51, -0.13).
func MakeCoordinates(overrides ...func(*models.Coordinates)) models.Coordinates {
	c := models.Coordinates{Lat: 51.51, Lon: -0.13}
	for _, o := range overrides {
		o(&c)
	}
	return c
}

// MakeLocation returns a valid Location with a generated UUID and current timestamp.
func MakeLocation(overrides ...func(*models.Location)) models.Location {
	loc := models.Location{
		ID:          uuid.New().String(),
		Name:        "London",
		Coordinates: MakeCoordinates(),
		CreatedAt:   time.Now().UTC(),
	}
	for _, o := range overrides {
		o(&loc)
	}
	return loc
}

// MakeLocationCreate returns a valid LocationCreate request body for London.
func MakeLocationCreate(overrides ...func(*models.LocationCreate)) models.LocationCreate {
	lc := models.LocationCreate{Name: "London", Lat: 51.51, Lon: -0.13}
	for _, o := range overrides {
		o(&lc)
	}
	return lc
}

// MakeLocationUpdate returns a pointer-based LocationUpdate with only Name set.
func MakeLocationUpdate(overrides ...func(*models.LocationUpdate)) models.LocationUpdate {
	name := "Updated London"
	lu := models.LocationUpdate{Name: &name}
	for _, o := range overrides {
		o(&lu)
	}
	return lu
}

// MakeCurrentWeather returns a valid CurrentWeather domain object (15°C, London).
func MakeCurrentWeather(overrides ...func(*models.CurrentWeather)) models.CurrentWeather {
	cw := models.CurrentWeather{
		Temperature:   15.0,
		FeelsLike:     13.5,
		Humidity:      72.0,
		Pressure:      1013.0,
		WindSpeed:     5.5,
		WindDirection: 270,
		Description:   "light rain",
		Icon:          "10d",
		Timestamp:     time.Now().Unix(),
		LocationName:  "London",
		Units:         models.UnitCelsius,
	}
	for _, o := range overrides {
		o(&cw)
	}
	return cw
}

// MakeForecastDay returns a valid ForecastDay for tomorrow.
func MakeForecastDay(overrides ...func(*models.ForecastDay)) models.ForecastDay {
	tomorrow := time.Now().AddDate(0, 0, 1).UTC().Format("2006-01-02")
	fd := models.ForecastDay{
		Date:        tomorrow,
		TempMin:     10.0,
		TempMax:     18.0,
		Humidity:    68.0,
		Description: "scattered clouds",
		Icon:        "03d",
	}
	for _, o := range overrides {
		o(&fd)
	}
	return fd
}

// MakeForecast returns a valid 3-day Forecast for London.
func MakeForecast(overrides ...func(*models.Forecast)) models.Forecast {
	days := make([]models.ForecastDay, 3)
	for i := range days {
		offset := i + 1
		days[i] = models.ForecastDay{
			Date:        time.Now().AddDate(0, 0, offset).UTC().Format("2006-01-02"),
			TempMin:     float64(8 + i),
			TempMax:     float64(16 + i),
			Humidity:    65.0,
			Description: "partly cloudy",
			Icon:        "02d",
		}
	}
	f := models.Forecast{LocationName: "London", Units: models.UnitCelsius, Days: days}
	for _, o := range overrides {
		o(&f)
	}
	return f
}

// MakeWeatherAlert returns a high-wind medium-severity alert.
func MakeWeatherAlert(overrides ...func(*models.WeatherAlert)) models.WeatherAlert {
	a := models.WeatherAlert{
		AlertType: "high_wind",
		Message:   "High wind speed: 22.0 m/s",
		Severity:  models.SeverityMedium,
		Value:     22.0,
		Threshold: 20.0,
	}
	for _, o := range overrides {
		o(&a)
	}
	return a
}

// MakeOwmWeatherEntry returns a valid OWM weather condition entry.
func MakeOwmWeatherEntry(overrides ...func(*models.OwmWeatherEntry)) models.OwmWeatherEntry {
	e := models.OwmWeatherEntry{ID: 500, Main: "Rain", Description: "light rain", Icon: "10d"}
	for _, o := range overrides {
		o(&e)
	}
	return e
}

// MakeOwmCurrentWeatherResponse returns a valid raw OWM /weather response for London at 15°C.
func MakeOwmCurrentWeatherResponse(overrides ...func(*models.OwmCurrentWeatherResponse)) models.OwmCurrentWeatherResponse {
	r := models.OwmCurrentWeatherResponse{
		Weather: []models.OwmWeatherEntry{MakeOwmWeatherEntry()},
		Main: models.OwmMain{
			Temp:      15.0,
			FeelsLike: 13.5,
			TempMin:   12.0,
			TempMax:   17.0,
			Pressure:  1013.0,
			Humidity:  72.0,
		},
		Wind: models.OwmWind{Speed: 5.5, Deg: 270},
		Dt:   time.Now().Unix(),
		Name: "London",
	}
	for _, o := range overrides {
		o(&r)
	}
	return r
}

// MakeOwmForecastItem returns a valid OWM 3-hour forecast list entry.
func MakeOwmForecastItem(dtOffset int, overrides ...func(*models.OwmForecastItem)) models.OwmForecastItem {
	t := time.Now().AddDate(0, 0, dtOffset).UTC().Truncate(24 * time.Hour)
	item := models.OwmForecastItem{
		Dt: t.Unix(),
		Main: models.OwmMain{
			Temp:     float64(12 + dtOffset),
			TempMin:  float64(10 + dtOffset),
			TempMax:  float64(15 + dtOffset),
			Humidity: 65.0,
			Pressure: 1012.0,
		},
		Weather: []models.OwmWeatherEntry{MakeOwmWeatherEntry()},
		DtTxt:   fmt.Sprintf("%s 12:00:00", t.Format("2006-01-02")),
	}
	for _, o := range overrides {
		o(&item)
	}
	return item
}

// MakeOwmForecastResponse returns a valid OWM /forecast response with 3 daily entries.
func MakeOwmForecastResponse(overrides ...func(*models.OwmForecastResponse)) models.OwmForecastResponse {
	r := models.OwmForecastResponse{
		List: []models.OwmForecastItem{
			MakeOwmForecastItem(1),
			MakeOwmForecastItem(2),
			MakeOwmForecastItem(3),
		},
		City: models.OwmCity{Name: "London"},
	}
	for _, o := range overrides {
		o(&r)
	}
	return r
}
