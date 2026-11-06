package oba

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
