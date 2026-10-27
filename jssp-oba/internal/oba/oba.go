package oba

import (
	"fmt"
	"math/rand"

	"jssp-oba/internal/neighborhood"
	"jssp-oba/internal/objective"
	"jssp-oba/internal/problem"
	"jssp-oba/internal/solution"
)

// Result contiene la meor solución encontrada y sus métricas
type Result struct {
	BestSchedule   *solution.Schedule
	Best           objective.Evaluation
	Iterations     int
	Accepted       int
	FinalThreshold float64
}

// Solve ejecuta OBA-N desde una solución inicial
func Solve(
	in *problem.Instance,
	initial *solution.Schedule,
	config Config,
) (*Result, error) {
	if in == nil || initial == nil {
		return nil, fmt.Errorf("la instancia y la solución inicial son obligatorias")
	}
	if config.Iterations <= 0 {
		return nil, fmt.Errorf("Iterations debe ser mayor que cero")
	}
	if config.Delta <= 0 {
		return nil, fmt.Errorf("Delta debe ser mayor que cero")
	}
	if config.EasyAge < 1 {
		return nil, fmt.Errorf("EasyAge debe ser al menos uno")
	}

	current := initial.Clone()
	currentEval, err := objective.Evaluate(in, current)
	if err != nil {
		return nil, fmt.Errorf("evaluar solución inicial: %w", err)
	}

	best := current.Clone()
	bestEval := currentEval

	threshold := 0.0
	count := 1
	age := config.Iterations // Simula la edad inicial grande
	accepted := 0
	completed := 0

	for i := 0; i < config.Iterations; i++ {
		completed++
		neighbor, err := neighborhood.Generate(current)
		if err != nil {
			completed--
			break
		}

		neighborEval, err := objective.Evaluate(in, neighbor)
		if err != nil {
			return nil, fmt.Errorf("evaluar vecino: %w", err)
		}

		// OBA2-N no permite umbrales negativos
		threshold = clampThreshold(threshold)

		if neighborEval.Cost < currentEval.Cost+threshold {
			//ACeptar el vecino
			current = neighbor
			currentEval = neighborEval
			accepted++
			age = 0

			if currentEval.Cost < bestEval.Cost {
				best = current.Clone()
				bestEval = currentEval
			}

			// Si el movimiento anterior también fue fácil de aceptar
			// intensificamos el descenso del umbral
			if prevAge < config.EasyAge {
				count++
			} else {
				count = 1
			}

			threshold -= float64(count) * config.Delta
		} else {
			//Rechaza el vecino e incrementa la edad
			age++
			threshold += config.Delta / float64(config.EasyAge)
			count = 1
		}

		threshold = clampThreshold(threshold)
	}

	return &Result{
		BestSchedule:   best,
		Best:           bestEval,
		Iterations:     config.Iterations,
		Accepted:       accepted,
		FinalThreshold: threshold,
	}, nil
}
