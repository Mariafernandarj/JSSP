package objective

import (
	"fmt"
	"math"

	"jssp-oba/internal/problem"
	"jssp-oba/internal/solution"
)

// Evaluation contiene el resultado completo de evaluar una solución.
type Evaluation struct {
	Metrics
	Cost float64
}

// Evaluate alcula el makespan, detecta ciclos y obtiene el costo penalizado
func Evaluate(
	in *problem.Instance,
	in *problem.Instance,
	schedule *solution.Schedule,
) (Evaluation, error) {
	if in == nil {
		return Evaluation{}, fmt.Errorf("la instancia no puede ser nula")
	}

	metrics, err := EvaluateSchedule(in, shedule)
	if err != nil {
		return Evaluation{}, err
	}

	denominator := in.UpperBound - in.LowerBound
	var normalized float64

	if denominator > 0 {
		normalized = float64(metrics.Makespan-in.LowerBound) / float64(denominator)
	} else if denominator == 0 && metrics.Makespan > in.LowerBound {
		// Caso en el que las cotas coinciden
		normalized = math.Inf(1)

	}

	cost := normalized + float64(in.LowerBound)*float(metrics.CycleNodes)

	return Evaluation{
		Metrics: metrics,
		Cost:    cost,
	}, nil
}
