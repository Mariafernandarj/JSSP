package solver

import (
	"fmt"
	"math/rand"

	"jssp-oba/internal/problem"
	"jssp-oba/internal/solution"
	"sort"
)

func RandomInitialSchedule(
	in *problem.Instance,
	rng *rand.Rand,
) (*solution.Schedule, error) {
	if in == nil {
		return nil, fmt.Errorf("la instancia es nil")
	}

	// Agrupar las operaciones por trabajo y ordenarlas por su posición
	numJobs := 0
	for _, op := range in.Operations {
		if op.JobID+1 > numJobs {
			numJobs = op.JobID + 1
		}
	}

	jobOps := make([][]problem.Operation, numJobs)
	for _, op := range in.Operations {
		if op.Machine < 0 || op.Machine >= in.NumMachines {
			return nil, fmt.Errorf("máquina inválida: %d", op.Machine)
		}
		jobOps[op.JobID] = append(jobOps[op.JobID], op)
	}
	for j := range jobOps {
		ops := jobOps[j]
		sort.Slice(ops, func(a, b int) bool {
			return ops[a].Index < ops[b].Index
		})
	}

	// Trabajos que aún tienen operaciones pendientes
	next := make([]int, numJobs) // siguiente operación de cada trabajo
	active := make([]int, 0, numJobs)
	for j := 0; j < numJobs; j++ {
		if len(jobOps[j]) > 0 {
			active = append(active, j)
		}
	}

	// Construcción incremental: elegir un trabajo al azar y agregar
	// su siguiente operación al final de la lista de su máquina
	machineOrder := make([][]int, in.NumMachines)

	for len(active) > 0 {
		k := rng.Intn(len(active))
		j := active[k]

		op := jobOps[j][next[j]]
		machineOrder[op.Machine] = append(machineOrder[op.Machine], op.ID)
		next[j]++

		if next[j] == len(jobOps[j]) {
			// Trabajo terminado: sacarlo de la lista de activos.
			active[k] = active[len(active)-1]
			active = active[:len(active)-1]
		}
	}

	return &solution.Schedule{MachineOrder: machineOrder}, nil

}
