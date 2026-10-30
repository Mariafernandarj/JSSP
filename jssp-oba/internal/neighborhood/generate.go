package neighborhood

import (
	"fmt"
	"jssp-oba/internal/solution"
	"math/rand"
)

// Generate genra un vecino medinat en intercambio adyacnete aleatorio
func Generate(
	schedule *solution.Schedule,
) (*solution.Schedule, error) {
	if schedule == nil {
		return nil, fmt.Errorf("la solución no puede ser nula")
	}

	type move struct {
		machine  int
		position int
	}

	moves := make([]move, 0)

	for machine, order := range schedule.MachineOrder {
		for position := 0; position+1 < len(order); position++ {
			moves = append(moves, move{
				machine:  machine,
				position: position,
			})
		}
	}

	if len(moves) == 0 {
		return nil, fmt.Errorf("no existen movimientos posibles")
	}

	selected := moves[rand.Intn(len(moves))]

	return SwapAdjacent(
		schedule,
		selected.machine,
		selected.position,
	)
}
