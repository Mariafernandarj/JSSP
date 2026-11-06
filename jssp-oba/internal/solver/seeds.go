package solver

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadSeeds(path string) ([]int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("abrir archivo de semillas: %w", err)
	}
	defer file.Close()

	var seeds []int64

	scanner := bufio.NewScanner(file)
	line := 0

	for scanner.Scan() {
		line++

		value := strings.TrimSpace(scanner.Text())

		// Ignorar líneas vacías y comentarios
		if value == "" || strings.HasPrefix(value, "#") {
			continue
		}

		seed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"semilla inválida en línea %d: %s",
				line,
				value,
			)
		}

		seeds = append(seeds, seed)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("leer archivo de semillas: %w", err)
	}

	if len(seeds) == 0 {
		return nil, fmt.Errorf("el archivo no contiene semillas")
	}

	return seeds, nil
}
