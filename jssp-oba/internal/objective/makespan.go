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

// CalculateMakespanRelaxed: Kahn que, al atascarse, fuerza el nodo
// con menos predecesores pendientes y cuenta los arcos ignorados
func CalculateMakespanRelaxed(g *Graph) (Metrics, int, error) {
	n := len(g.Adjacency)
	degree := append([]int(nil), g.InDegree...)
	finish := make([]int, n)
	done := make([]bool, n)
	queue := make([]int, 0, n)
	violations := 0

	for u := 0; u < n; u++ {
		if degree[u] == 0 {
			queue = append(queue, u)
			finish[u] = g.Duration[u]
			done[u] = true
		}
	}

	head, processed, makespan := 0, 0, 0
	for processed < n {
		if head == len(queue) {
			// atascado: forzar el nodo pendiente con menor grado
			best := -1
			for v := 0; v < n; v++ {
				if !done[v] && (best == -1 || degree[v] < degree[best]) {
					best = v
				}
			}
			violations += degree[best]
			degree[best] = 0
			done[best] = true
			finish[best] += g.Duration[best] // se suma a lo que ya acumuló
			queue = append(queue, best)
		}
		u := queue[head]
		head++
		processed++
		if finish[u] > makespan {
			makespan = finish[u]
		}
		for _, v := range g.Adjacency[u] {
			if done[v] {
				continue // arco ignorado, ya se contó como violación
			}
			if c := finish[u] + g.Duration[v]; c > finish[v] {
				finish[v] = c
			}
			degree[v]--
			if degree[v] == 0 {
				done[v] = true
				queue = append(queue, v)
			}
		}
	}
	return Metrics{Makespan: makespan, CycleNodes: violations,
		Processed: processed, Feasible: violations == 0}, violations, nil
}

// EvaluateSchedule construe la gráfica y obtiene sus métricas
func EvaluateSchedule(
	in *problem.Instance,
	schedule *solution.Schedule,
) (Metrics, error) {
	g, err := BuildGraph(in, schedule)
	if err != nil {
		return Metrics{}, err
	}
	metrics, _, err := CalculateMakespanRelaxed(g)
	return metrics, err
}
