package main

import (
	"fmt"
	"jssp-oba/internal/problem"
	"jssp-oba/internal/solution"
	"log"
	"os"
)

func main() {
	path := "data/instances/example.txt"

	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	instance, err := problem.ParseFile(path)
	if err != nil {
		log.Fatal("Error al cargar la instancia: %v", err)
	}

	fmt.Println("Instancia cargada correctamente")
	fmt.Printf("Trabajos: %d\n", instance.NumJobs)
	fmt.Printf("Máquinas: %d\n", instance.NumMachines)
	fmt.Printf("Operaciones: %d\n", instance.TotalOperations())
	fmt.Printf("Cota inferior: %d\n", instance.LowerBound)
	fmt.Printf("Cota superior: %d\n", instance.UpperBound)

	//Construimos una cadena inicial en el orden original de las operacions de cada trabajo
	var sequence []int
	for _, job := range instance.Jobs {
		for _, op := range job.Operations {
			sequence = append(sequence, op.ID)
		}
	}

	schedule, err := solution.FromOperationSequence(
		instance, sequence,
	)
	if err != nil {
		log.Fatalf("Error al reprentar la solución: %v", err)
	}

	fmt.Println("\nCadena de máquinas:")
	for machine, operations := range schedule.MachineOrder {
		fmt.Printf("M%d: %v\n", machine, operations)
	}
}
