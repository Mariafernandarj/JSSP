package main

import (
	"fmt"
	"jssp-oba/internal/oba"
	"jssp-oba/internal/objective"
	"jssp-oba/internal/problem"
	"jssp-oba/internal/solution"
	"log"
	"os"
	"time"
)

const (
	MINIC_EXIT_OK       = 0
	MINIC_EXIT_INVALID  = 1
	MINIC_EXIT_INTERNAL = 2
)

func main() {
	os.Exit(run())
}

func run() int {
	instancePath := "data/instances/example.txt"

	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "Uso: go run ./cmd/jssp-oba [archivo_instancia]")
		return MINIC_EXIT_INVALID
	}

	fmt.Println("========================================")
	fmt.Println(" Job Shop Scheduling Problem - OBA2-N")
	fmt.Println("========================================")
	fmt.Printf("Instancia:      %s\n", instancePath)
	fmt.Printf("Trabajos:       %d\n", instance.NumJobs)
	fmt.Printf("Máquinas:       %d\n", instance.NumMachines)
	fmt.Printf("Operaciones:    %d\n", instance.TotalOperations())
	fmt.Printf("Cota inferior:  %d\n", instance.LowerBound)
	fmt.Printf("Cota superior:  %d\n", instance.UpperBound)

	// se crea una solución inicial
	sequence := make([]int, 0, instance.TotalOperations())

	for _, job := range instance.Jobs {
		for _, op := range job.Operations {
			sequence = append(sequence, op.ID)
		}
	}

	initialEval, err := objective.Avaluate(instance, initial)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al evaluar la solución inicial: %v\n", err)
		return MINIC_EXIT_INTERNAL
	}

	fmt.Println("\n--- Solución inicial ---")
	printSchedule(initial)
	printEvaluation(initialEval)

	//Cunfirar OBA2-N
	config := oba.DefaultConfig()

	fmt.Println("\n--- Configuración OBA2-N ---")
	fmt.Printf("Iteraciones: %d\n", config.Iterations)
	fmt.Printf("Delta:      %.6f\n", config.Delta)
	fmt.Printf("EasyAge:    %d\n", config.EasyAge)

	// Ejecutar algoritmo y medir tiempo
	start := time.Now()

	result, err := oba.Solve(instance, initial, config)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error durante la búsqueda OBA2-N: %v\n", err)
		return MINIC_EXIT_INTERNAL
	}

	// 7. Mostrar la mejor solución encontrada.
	fmt.Println("\n--- Resultado de OBA2-N ---")
	printSchedule(result.BestSchedule)
	printEvaluation(result.Best)

	fmt.Printf("Iteraciones ejecutadas: %d\n", result.Iterations)
	fmt.Printf("Movimientos aceptados:  %d\n", result.Accepted)
	fmt.Printf("Umbral final:           %.6f\n", result.FinalThreshold)
	fmt.Printf("Tiempo de ejecución:    %s\n", elapsed)

	return MINIC_EXIT_OK
}

// printSchedule muestra el orden de procesamiento de cada máquina
func printSchedule(schedule *solution.Schedule) {
	if schedule == nil {
		fmt.Println("No hay solución disponible")
		return
	}

	fmt.Println("Orden de operaciones por máquina:")
	for machine, order := range schedule.MachineOrder {
		fmt.Printf(" Máquina%d", machine)

		for i, operationID := range order {
			if i > 0 {
				fmt.Print(" -> ")
			}
			fmt.Printf("0%d", operationID)
		}
		fmt.Println()
	}
}

// printEvaluation muestar las métricas y el costo de una solución
func printEvaluation(evaluation objective.Evaluation) {
	fmt.Printf("Makespan:           %d\n", evaluation.Makespan)
	fmt.Printf("Nodos en ciclos:     %d\n", evaluation.CycleNodes)
	fmt.Printf("Operaciones procesadas: %d\n", evaluation.Processed)
	fmt.Printf("Factible:            %t\n", evaluation.Feasible)
	fmt.Printf("Costo penalizado:    %.6f\n", evaluation.Cost)
}
