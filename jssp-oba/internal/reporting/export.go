package reporting

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"jssp-oba/internal/oba"
	"jssp-oba/internal/problem"
)

func ExportResult(
	outputDir string,
	instance *problem.Instance,
	result oba.Result,
	elapsed time.Duration,
) (string, error) {

	if instance == nil {
		return "", fmt.Errorf("la instancia es nil")
	}

	if result.BestSchedule == nil {
		return "", fmt.Errorf("la mejor solucion es nil")
	}

	// Crear el directorio si no existe.
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("no se pudo crear el directorio: %w", err)
	}

	// Nombre único para cada ejecución.
	filename := fmt.Sprintf(
		"resultado_%s.txt",
		time.Now().Format("20060102_150405.000000000"),
	)

	path := filepath.Join(outputDir, filename)

	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("no se pudo crear el archivo: %w", err)
	}
	defer file.Close()

	var builder strings.Builder

	// Crear un mapa para recuperar los datos de cada operación.
	operations := make(map[int]problem.Operation)

	for _, op := range instance.Operations {
		operations[op.ID] = op
	}

	// ------------------------------------------------
	// Solución
	// ------------------------------------------------

	builder.WriteString("Solucion:\n[ ")

	first := true

	for machine, operationIDs := range result.BestSchedule.MachineOrder {
		for _, id := range operationIDs {

			op, exists := operations[id]
			if !exists {
				return "", fmt.Errorf(
					"no se encontro la operacion con ID %d",
					id,
				)
			}

			if !first {
				builder.WriteString(", ")
			}

			fmt.Fprintf(
				&builder,
				"[%d, %d, %d]",
				op.ID,
				machine,
				op.Duration,
			)

			first = false
		}
	}

	builder.WriteString(" ]\n\n")

	// ------------------------------------------------
	// Información de la mejor solución
	// ------------------------------------------------

	builder.WriteString("Informacion de la mejor solucion:\n")

	fmt.Fprintf(
		&builder,
		"Makespan: %d\n",
		result.Best.Metrics.Makespan,
	)

	fmt.Fprintf(
		&builder,
		"Nodos en ciclos: %d\n",
		result.Best.Metrics.CycleNodes,
	)

	fmt.Fprintf(
		&builder,
		"Operaciones procesadas: %d\n",
		result.Best.Metrics.Processed,
	)

	fmt.Fprintf(
		&builder,
		"Factible: %t\n",
		result.Best.Metrics.Feasible,
	)

	fmt.Fprintf(
		&builder,
		"Costo: %.6f\n\n",
		result.Best.Cost,
	)

	// ------------------------------------------------
	// Información de la ejecución
	// ------------------------------------------------

	builder.WriteString("Informacion de la ejecucion:\n")

	fmt.Fprintf(
		&builder,
		"Iteraciones: %d\n",
		result.Iterations,
	)

	fmt.Fprintf(
		&builder,
		"Movimientos aceptados: %d\n",
		result.Accepted,
	)

	fmt.Fprintf(
		&builder,
		"Umbral final: %.6f\n",
		result.FinalThreshold,
	)

	fmt.Fprintf(
		&builder,
		"Tiempo de ejecucion: %.6f segundos\n",
		elapsed.Seconds(),
	)

	// Escribir el contenido en el archivo.
	if _, err := file.WriteString(builder.String()); err != nil {
		return "", fmt.Errorf("error al escribir el archivo: %w", err)
	}

	return path, nil
}
