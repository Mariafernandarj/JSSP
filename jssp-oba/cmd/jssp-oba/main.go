package main

import (
	"log"

	"jssp-oba/internal/oba"
	"jssp-oba/internal/problem"
	"jssp-oba/internal/solver"
)

func main() {
	// Leer instancia
	instance, err := problem.ParseFile(
		"data/instances/example.txt",
	)

	if err != nil {
		log.Fatal(err)
	}

	// Configuración de OBA2-N
	config := oba.DefaultConfig()

	// Semillas de las corridas.
	seeds := []int64{
		10,
		20,
		30,
		40,
		50,
	}

	// Ejecutar experimento.
	err = solver.RunExperiments(
		instance,
		config,
		seeds,
	)

	if err != nil {
		log.Fatal(err)
	}
}
