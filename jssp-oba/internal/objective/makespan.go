package objective

import (
	"fmt"
	"jssp-oba/internal/problem"
	"jssp-oba/internal/solution"
)

// Matrices contiene las métricas obtenidas al evaluar la gráfica
type Metrics struct {
	Makespan   int
	CycleNodes int
	Processed  int
	Feasible   bool
}

// CalculateMakespan ejecuta el algoritmo de Kahn y calcula
// los tiempos de finalización más tenpranos
func CalculateMakespan(g *Graph) (Metrics, error) {
	if g == nil {
		return Metrics{}, fmt.Errorf("la gáfica no puede ser nil")
	}

	n := len(g.Adjacency)
	if len(g.InDegree) != n || len(g.Duration) != n {
		return Metrics{}, fmt.Errorf("estructura de la gŕafica es inconsitente")
	}

	degree := append([]int(nil), g.InDegree...)
	finish := make([]int, n)
	queue := make([]int, 0, n)

	for u := 0; u < n; u++ {
		if degree[u] == 0 {
			queue = append(queue, u)
			finish[u] = g.Duration[u]
		}
	}

	processed := 0
	makespan := 0

	for head := 0; head < len(queue); head++ {
		u := queue[head]
		processed++

		if finish[u] > makespan {
			makespan = finish[u]
		}

		for _, v := range g.Adjacency[u] {
			candidate := finish[u] + g.Duration[v]
			if candidate > finish[v] {
				finish[v] = candidate
			}

			degree[v]--
			if degree[v] == 0 {
				queue = append(queue, v)
			}
		}
	}

	// Kahn deja con grado positivo los nodos no procesados
	cycleNodes := 0
	for _, d := range degree {
		if d > 0 {
			cycleNodes++
		}
	}

	return Metrics{
		Makespan:   makespan,
		CycleNodes: cycleNodes,
		Processed:  processed,
		Feasible:   processed == n,
	}, nil
}

// EvaluateSchedule construe la gráfica y obtiene sus métricas
func EvaluateSchedule(
	in *problem.Instance,
	schedule *solution.Shedule,
) (Metrics, error) {
	g, err := BuildGraph(in, schedule)
	if err != nil {
		return Metrics{}, err
	}
	return CalculateMakespan(g)
}
