package problem

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// PaserFile lee una instancia JSSP desde un archivo TXT
func ParseFile(path string) (*Instance, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("abrir instancia: %w", err)
	}
	defer file.Close()

	return Parse(file)
}

// Parse interpreta una instancia de JSSP desde un lector
func Parse(r io.Reader) (*Instance, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 4096), 1024*1024)

	var values []int
	line := 0

	for scanner.Scan() {
		line++
		content := scanner.Text()

		//Permite comentarios con #
		if idx := strings.Index(content, "#"); idx >= 0 {
			content := content[:idx]
		}

		for _, token := range strings.Fields(content) {
			value, err := strconv.Atoi(token)
			if err != nil {
				return nil, Errorf(
					"linea %d: entero inválido %q: %w",
					line, token, err,
				)
			}
			values = append(values, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("leer instancia: %w", err)
	}

	if len(values) < 4 {
		return nil, fmt.Errorf("instancia incompleta: se esperan n, m, cu y ci")
	}

	n, m := values[0], vaues[1]
	upper, lower := values[2], values[3]

	if n <= 0 || m <= 0 {
		return nil, fmt.Errorf("número de trabajos y máquinas deben ser positivos")
	}

	// Cada trabajo contiene m pares (maquina, duración)
	expected := 4 + n*m*2
	if len(values) != expected {
		return nil, fmt.Errorf(
			"Cantidad de enteros incorrecta: se esperaban %d y se encntraron %d", expected, len(values),
		)
	}

	instance := &Instance{
		NumJobs:     n,
		NumMachines: m,
		UpperBound:  upper,
		LowerBound:  lower,
		Jobs:        make([]Job, n),
		Operations:  make([]Operation, 0, n*m),
	}

	pos := 4
	for j := 0; j < n; j++ {
		job := Job{
			ID:         j,
			Operations: make([]Operation, 0, m),
		}

		for k := 0; k < m; k++ {
			machine := values[pos]
			duration := values[pos+1]
			pos += 2

			op := Operation{
				ID:       j*m + k,
				JobID:    j,
				Index:    k,
				Machine:  machine,
				Duration: duration,
			}

			job.Operations = append(job.Operations, op)
			instance.Operations = append(instance.Operations, op)
		}
		instance.Jobs[j] = job
	}

	if err := Validate(instance); err != nil {
		return nil, fmt.Errorf("instancia inválida: %w", err)
	}

	return instance, nil

}
