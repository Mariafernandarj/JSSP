package main

import (
	"flag"
	"log"

	"jssp-oba/internal/oba"
	"jssp-oba/internal/problem"
	"jssp-oba/internal/solver"
)

func main() {
	// Leer instancia
	instancePath := flag.String(
		"instance",
		"",
		"Ruta del archivo de instancia JSSP",
	)

	seedsPath := flag.String(
		"seeds",
		"",
		"Ruta del archivo de semillas",
	)

	configPath := flag.String(
		"config",
		"",
		"Ruta del archivo de configuración OBA2-N",
	)

	flag.Parse()

	if *instancePath == "" || *seedsPath == "" || *configPath == "" {
		log.Fatal("Debes especificar -instance, -seeds y -config")
	}

	// Leer instancia
	instance, err := problem.ParseFile(*instancePath)
	if err != nil {
		log.Fatal(err)
	}

	// Leer configuración
	config, err := oba.LoadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Leer semillas
	seeds, err := solver.LoadSeeds(*seedsPath)
	if err != nil {
		log.Fatal(err)
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
