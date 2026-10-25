package solution

import (
	"fmt"
	"jssp-oba/internal/problem"
)

// FromOperationSequence convierte una cadena de IDs de operaciones
// en el oden de procesamiento de cada máquina
func FromOperationSequence(
	in *problem.Instance,
	sequence []int,
) (*Schedule, error) {
	if in == nil {
		return nil, fmt.Errorf("la instancia es nil")
	}

	if len(sequence) != in.TotalOperations() {
		return nil, fmt.Errorf(
			"se esperaban %d peraciones y hay %d",
			in.TotalOperations(), len(sequence),
		)
	}

	schedule := NewShedule(in.NumMachines)
	seen := make([]bool, in.TotalOperations())

	for _, id := range sequence {
		op, ok := in.OperationByID(id)
		if !ok {
			return nil, fmt.Errorf(
				"ID de operación fuera de rango: %d", id,
			)
		}

		if seen[id] {
			return nil, fmt.Errorf(
				"OPeración duplicada: %d", id,
			)
		}
		seen[id] = true

		schedule.AddOperation(op.Machine, id)
	}

	return schedule, nil
}
