package solver

import (
	"fmt"
	"math/rand"
	"time"

	"jssp-oba/internal/oba"
	"jssp-oba/internal/objective"
	"jssp-oba/internal/problem"
	"jssp-oba/internal/reporting"
)

func RunExperiments(
	in *problem.Instance,
	config oba.Config,
	seeds []int64,
) error {

	for i, seed := range seeds {

		fmt.Printf("\n========== CORRIDA %d ==========\n", i+1)
		fmt.Printf("Semilla: %d\n", seed)

		// Crear generador independiente
		rng := rand.New(rand.NewSource(seed))

		// Generar solución inicial
		initial, err := RandomInitialSchedule(in, rng)

		if err != nil {
			return err
		}

		// Evaluar solución inicial
		initialEval, _ := objective.Evaluate(in, initial)

		fmt.Printf(
			"Solución inicial factible: %t\n",
			initialEval.Metrics.Feasible,
		)

		// Ejecutar OBA2-N
		start := time.Now()

		result, err := oba.Solve(
			in,
			initial,
			config,
			rng,
		)

		if err != nil {
			return err
		}

		elapsed := time.Since(start)

		// Mostrar resultados
		fmt.Printf(
			"Makespan final: %d\n",
			result.Best.Metrics.Makespan,
		)

		fmt.Printf(
			"Factible: %t\n",
			result.Best.Metrics.Feasible,
		)

		fmt.Printf(
			"Tiempo: %.4f segundos\n",
			elapsed.Seconds(),
		)

		// Exportar resultado
		// Adaptar esta llamada a la firma real de ExportResult
		_, err = reporting.ExportResult(
			"data/results",
			in,
			*result,
			elapsed,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
