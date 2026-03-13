package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"copilot-go-advanced/internal/utils"
)

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		celsius  float64
		expected float64
	}{
		{0, 32},
		{100, 212},
		{-40, -40},   // crossover point
		{37, 98.6},
		{-20, -4},
		{15, 59},
	}
	for _, tt := range tests {
		got := utils.CelsiusToFahrenheit(tt.celsius)
		assert.InDelta(t, tt.expected, got, 0.01, "CelsiusToFahrenheit(%v)", tt.celsius)
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		celsius  float64
		expected float64
	}{
		{0, 273.15},
		{100, 373.15},
		{-273.15, 0},
		{15, 288.15},
	}
	for _, tt := range tests {
		got := utils.CelsiusToKelvin(tt.celsius)
		assert.InDelta(t, tt.expected, got, 0.01, "CelsiusToKelvin(%v)", tt.celsius)
	}
}

func TestFahrenheitToCelsius(t *testing.T) {
	tests := []struct {
		fahrenheit float64
		expected   float64
	}{
		{32, 0},
		{212, 100},
		{-40, -40}, // crossover point
		{98.6, 37},
	}
	for _, tt := range tests {
		got := utils.FahrenheitToCelsius(tt.fahrenheit)
		assert.InDelta(t, tt.expected, got, 0.01, "FahrenheitToCelsius(%v)", tt.fahrenheit)
	}
}

func TestMpsToKmh(t *testing.T) {
	tests := []struct {
		mps      float64
		expected float64
	}{
		{0, 0},
		{1, 3.6},
		{10, 36},
		{5.5, 19.8},
	}
	for _, tt := range tests {
		got := utils.MpsToKmh(tt.mps)
		assert.InDelta(t, tt.expected, got, 0.01, "MpsToKmh(%v)", tt.mps)
	}
}

func TestMpsToMph(t *testing.T) {
	tests := []struct {
		mps      float64
		expected float64
	}{
		{0, 0},
		{1, 2.24},
		{10, 22.37},
	}
	for _, tt := range tests {
		got := utils.MpsToMph(tt.mps)
		assert.InDelta(t, tt.expected, got, 0.01, "MpsToMph(%v)", tt.mps)
	}
}

func TestDegreesToCompass(t *testing.T) {
	tests := []struct {
		degrees  float64
		expected string
	}{
		{0, "N"},
		{22.5, "NNE"},
		{45, "NE"},
		{67.5, "ENE"},
		{90, "E"},
		{112.5, "ESE"},
		{135, "SE"},
		{157.5, "SSE"},
		{180, "S"},
		{202.5, "SSW"},
		{225, "SW"},
		{247.5, "WSW"},
		{270, "W"},
		{292.5, "WNW"},
		{315, "NW"},
		{337.5, "NNW"},
		{360, "N"},     // wraps back to N
		{-90, "W"},     // negative degrees
		{11.24, "N"},   // boundary: just before NNE sector
		{348.75, "N"},  // boundary: sector boundary rounds up to N (15.5 rounds to 16 → 16%16=0)
	}
	for _, tt := range tests {
		got := utils.DegreesToCompass(tt.degrees)
		assert.Equal(t, tt.expected, got, "DegreesToCompass(%v)", tt.degrees)
	}
}

func TestRoundingPrecision(t *testing.T) {
	// Verify 2-decimal rounding
	assert.Equal(t, 98.6, utils.CelsiusToFahrenheit(37))
	assert.Equal(t, 288.15, utils.CelsiusToKelvin(15))
}
