package neighborhood

import (
	"fmt"
	"jssp-oba/internal/solution"
)

// SwapAdjacent interca,bia dos operaciones consecutivas de una máquina
func SwapAdjacent(
	schedule *solution.Schedule,
	machine int,
	position int,
) (*solution.Schedule, error) {
	if schedule == nil {
		return nil, fmt.Errorf("la solución no puede ser nil")
	}

	if machine < 0 || machine >= len(schedule.MachineOrder) {
		return nil, fmt.Errorf("máquina fuera de rango: %d", machine)
	}

	order := schedule.MachineOrder[machine]
	if position < 0 || position+1 >= len(order) {
		return nil, fmt.Errorf("posición inválida: %d", position)
	}

	neighbor := schedule.Clone()

	neighbor.MachineOrder[machine][position],
		neighbor.MachineOrder[machine][position+1] =
		neighbor.MachineOrder[machine][position+1],
		neighbor.MachineOrder[machine][position]

	return neighbor, nil

}
