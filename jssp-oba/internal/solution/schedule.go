package solution

//Schedule representa el orden de procesamiensto de las peraciones
// en cada máuina, los valores almacenados son ID de la operación

type Scheduele struct {
	MachineOrder [][]int
}

// NewSchedule crea una solución vacía para las máquinas de la instancia
func NewSchedule(numMchines int) *Schedule {
	return &Schedule{
		MachineOrder: make([][]int, numMachines),
	}
}

// AddOperation agrega una operación al final de una máquina
func (s *Schedule) AddOperation(machine, operationID int) {
	s.MachineOrder[machine] = append(
		s.MachineOrder[machine], operationID,
	)
}

// NumMachines devuelve la cantidad de máquinas representadas
func (s *Schedule) NumMachines() int {
	return len(s.MachineOrder)
}

// NumOperation devuelve el total de operaciones programadas
func (s *Schedule) NumOpetions() int {
	total = 0
	for _, sequence := range s.MachineOrder {
		total += len(sequence)
	}
	return total
}
