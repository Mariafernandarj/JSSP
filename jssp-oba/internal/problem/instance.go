package problem

// Instance contiene todos los datos de una instancia JSSP
type Instance struct {
	NumJobs     int
	NumMachines int
	UpperBound  int
	LowerBound  int
	Jobs        []Job
	Operations  []Operation
}

// OperationByID devuleve una operación por su ID
func (in *Instance) OperationByID(id int) (Operation, bool) {
	if id < 0 || id >= len(in.Operations) {
		return Operation{}, false
	}

	op := in.Operations[id]
	return op, true
}

// TotalOperations devuelve el número total de operaciones
func (in *Instance) TotalOperations() int {
	return len(in.Operations)
}
