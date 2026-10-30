package oba

// Config contiene los parámetros para OBA2-N
type Config struct {
	Iterations int     `json:"iterations"`
	Delta      float64 `json:"delta"`
	EasyAge    int     `json:"easy_age"`
}

// DefaultConfig devuelve parámetros iniciales
func DefaultConfig() Config {
	return Config{
		Iterations: 10000,
		Delta:      0.1,
		EasyAge:    10,
	}
}
