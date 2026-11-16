package main

import (
	"flag"
	"fmt"
	"jssp-oba/internal/oba"
	"jssp-oba/internal/problem"
	"jssp-oba/internal/solver"
	"log"
	"path/filepath"
	"strings"
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

	fmt.Printf("Cota inferior (LB): %d\n", instance.LowerBound)
	fmt.Printf("Cota superior (UB): %d\n", instance.UpperBound)

	// Nombre de la instancia para etiquetar los archivos (p. ej. "JSSP_50_20")
	instanceName := strings.TrimSuffix(
		filepath.Base(*instancePath),
		filepath.Ext(*instancePath),
	)

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
		instanceName,
		config,
		seeds,
	)

	if err != nil {
		log.Fatal(err)
	}
}
