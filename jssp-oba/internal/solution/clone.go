package solution

// Clone devuelve una copia independiente de la solución
func (s *Schedule) Clone() *Schedule {
	if s == nil {
		return nil
	}

	clone := NewSchedule(len(s.MachineOrder))

	for machine, sequence := range s.MachineOrder {
		clone.MachineOrder[machine] = append([]int(nil), sequence...)
	}

	return clone
}
