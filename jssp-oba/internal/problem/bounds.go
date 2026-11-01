package problem

// CalculateBounds calcula automáticamente las cotas inferior y superior
func CalculateBounds(in *Instance) {
	if in == nil {
		return
	}

	var upperBound int
	var maxJobTime int

	// Acumulación de las duraciones correspondientes a cada máquina
	machineTimes := make([]int, in.NumMachines)

	// Calcula la carga de cada trabajo y la suma total
	for _, job := range in.Jobs {
		var jobTime int

		for _, op := range job.Operations {
			jobTime += op.Duration
			upperBound += op.Duration

			if op.Machine >= 0 && op.Machine < in.NumMachines {
				machineTimes[op.Machine] += op.Duration
			}
		}

		// Obtiene el trabajo con mayor tiempo de procesamiento
		if jobTime > maxJobTime {
			maxJobTime = jobTime
		}
	}

	// Obtiene la máquina con mayor carga de trabajo
	var maxMachineTime int

	for _, machineTime := range machineTimes {
		if machineTime > maxMachineTime {
			maxMachineTime = machineTime
		}
	}

	// Asigna las cotas calculadas a la instancia
	in.UpperBound = upperBound

	if maxJobTime > maxMachineTime {
		in.LowerBound = maxJobTime
	} else {
		in.LowerBound = maxMachineTime
	}
}
