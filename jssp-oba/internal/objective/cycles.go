package objective

import (
	"fmt"
	"jssp-oba/internal/problem"
	"jssp-oba/internal/solution"
)

// Graph representa las precendencias entre oeraciones
type Graph struct {
	Adjacency [][]int
	InDegree  []int
	Duration  []int
}

// addEdge agrega una arista u -> v
// Se evita duplicar aristas para mantener correcto el grado del al entrada
func (g *Graph) addEdge(u, v int) {
	for _, neighbor := range g.Adjacency[u] {
		if neighbor == v {
			return
		}
	}

	g.Adjacency[u] = append(g.Adjacency[u], v)
	g.InDegree[v]++
}

// BuildGraph construye la grafica de precedencia de una solución
func BuildGraph(
	in *problem.Instance,
	schedule *solution.Schedule,
) (*Graph, error) {
	if in == nil || schedule == nil {
		return nil, fmt.Errorf("la instanca y la solución no puede ser nil")
	}

	if schedule.NumMachines() != in.NumMachines {
		return nil, fmt.Errorf(
			"la solución tiene %d máuquinas, se esperaban %d",
			schedule.NumMachines(),
			in.NumMachines,
		)
	}

	n := in.TotalOperations()

	g := &Graph{
		Adjacency: make([][]int, n),
		InDegree:  make([]int, n),
		Duration:  make([]int, n),
	}

	// Se gurdan la duración por ID de operación
	for _, op := range in.Operations {
		if op.ID < 0 || op.ID >= n {
			return nil, fmt.Errorf("ID de operación fuera de rango: %d", op.ID)
		}
		g.Duration[op.ID] = op.Duration
	}

	// Precedencias dentro de cada trabajo
	for _, job := range in.Jobs {
		for i := 0; i+1 < len(job.Operations); i++ {
			u := job.Operations[i].ID
			v := job.Operations[i+1].ID

			if u < 0 || u >= n || v < 0 || v >= n {
				return nil, fmt.Errorf("ID de operación inválido en trabajo %d", job.ID)
			}

			g.addEdge(u, v)
		}
	}

	// Precedencias impuestas por el orden de la mpquina
	seen := make([]bool, n)

	for machine, order := range schedule.MachineOrder {
		for _, id := range order {
			if id < 0 || id >= n {
				return nil, fmt.Errorf(
					"ID de operación fuera de rango en máquina %d: %d", machine, id,
				)
			}

			if seen[id] {
				return nil, fmt.Errorf(
					"la operación %d aparece más de una vez", id,
				)
			}
			seen[id] = true

			op, ok := in.OperationByID(id)
			if !ok {
				return nil, fmt.Errorf("no existe la operación %d", id)
			}
			if op.Machine != machine {
				return nil, fmt.Errorf("la operación %d pertenece a la máquina %d, no la %d",
					id, op.Machine, machine,
				)
			}
		}

		for i := 0; i+1 < len(order); i++ {
			g.addEdge(order[i], order[i+1])
		}
	}

	for id, exists := range seen {
		if !exists {
			return nil, fmt.Errorf(
				"la operación %d no aparece en ninguna máquina",
				id,
			)
		}
	}
	return g, nil
}
