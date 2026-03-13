package utils

import "math"

// CelsiusToFahrenheit converts a Celsius temperature to Fahrenheit, rounded to 2 decimal places.
func CelsiusToFahrenheit(c float64) float64 {
	return round2(c*9.0/5.0 + 32)
}

// CelsiusToKelvin converts a Celsius temperature to Kelvin, rounded to 2 decimal places.
func CelsiusToKelvin(c float64) float64 {
	return round2(c + 273.15)
}

// FahrenheitToCelsius converts a Fahrenheit temperature to Celsius, rounded to 2 decimal places.
func FahrenheitToCelsius(f float64) float64 {
	return round2((f - 32) * 5.0 / 9.0)
}

// MpsToKmh converts metres per second to kilometres per hour, rounded to 2 decimal places.
func MpsToKmh(mps float64) float64 {
	return round2(mps * 3.6)
}

// MpsToMph converts metres per second to miles per hour, rounded to 2 decimal places.
func MpsToMph(mps float64) float64 {
	return round2(mps * 2.23694)
}

// DegreesToCompass converts a wind direction in degrees to a 16-point compass rose label.
// Negative degrees are normalised. N is centred at 0°/360°, sectors are 22.5° each.
func DegreesToCompass(deg float64) string {
	compassPoints := [16]string{
		"N", "NNE", "NE", "ENE",
		"E", "ESE", "SE", "SSE",
		"S", "SSW", "SW", "WSW",
		"W", "WNW", "NW", "NNW",
	}

	// Normalise to [0, 360)
	normalized := math.Mod(deg, 360)
	if normalized < 0 {
		normalized += 360
	}

	index := int(math.Round(normalized/22.5)) % 16
	return compassPoints[index]
}

// round2 rounds a float64 to 2 decimal places.
func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
