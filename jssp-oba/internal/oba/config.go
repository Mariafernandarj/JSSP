package oba

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config contiene los parámetros para OBA2-N
type Config struct {
	Iterations int     `json:"iterations"`
	DeltaUnits float64 `json:"delta"`
	EasyAge    int     `json:"easy_age"`
}

// DefaultConfig devuelve parámetros iniciales
func DefaultConfig() Config {
	return Config{
		Iterations: 3000000,
		DeltaUnits: 5,
		EasyAge:    56,
	}
}

func LoadConfig(path string) (Config, error) {
	var config Config

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf(
			"leer configuración: %w",
			err,
		)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf(
			"interpretar configuración: %w",
			err,
		)
	}

	if config.Iterations <= 0 {
		return Config{}, fmt.Errorf(
			"Iterations debe ser mayor que cero",
		)
	}

	if config.DeltaUnits <= 0 {
		return Config{}, fmt.Errorf(
			"Delta debe ser mayor que cero",
		)
	}

	if config.EasyAge < 1 {
		return Config{}, fmt.Errorf(
			"EasyAge debe ser al menos uno",
		)
	}

	return config, nil
}
