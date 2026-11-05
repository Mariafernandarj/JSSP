package solver

import (
	"fmt"
	"math/rand"

	"jssp-oba/internal/problem"
	"jssp-oba/internal/solution"
)

func RandomInitialSchedule(
	in *problem.Instance,
	seed int64,
) (*solution.Schedule, error) {
	if in == nil {
		return nil, fmt.Errorf("la instancia es nil")
	}

	//Crear una lista de operaciones para cada máquina
	machineOrder := make([][]int, in.NumMachines)

	for _, op := range in.Operations {
		if op.Machine < 0 || op.Machine >= in.NumMachines {
			return nil, fmt.Errorf(
				"maquina invalida %d en operacion %d",
				op.Machine,
				op.ID,
			)
		}

		machineOrder[op.Machine] = append(
			machineOrder[op.Machine],
			op.ID,
		)
	}

	// Generador aleatorio.
	rng := rand.New(rand.NewSource(seed))

	// Mezclar las operaciones de cada máquina.
	for machine := range machineOrder {

		rng.Shuffle(
			len(machineOrder[machine]),
			func(i, j int) {
				machineOrder[machine][i],
					machineOrder[machine][j] =
					machineOrder[machine][j],
					machineOrder[machine][i]
			},
		)
	}

	return &solution.Schedule{
		MachineOrder: machineOrder,
	}, nil

}
