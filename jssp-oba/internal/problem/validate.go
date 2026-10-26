package problem

import "fmt"

// Valdate comprueba la estructura y las cotas de la instancia
func Validate(in *Instance) error {
	if in == nil {
		return fmt.Errorf("la instancia es nil")
	}

	n, m := in.NumJobs, in.NumMachines
	if n <= 0 || m <= 0 {
		return fmt.Errorf("n y m deben ser positivos")
	}

	if len(in.Jobs) != n {
		return fmt.Errorf(
			"se esperaban %d trabajos y hay %d",
			n, len(in.Jobs),
		)
	}

	if len(in.Operations) != n*m {
		return fmt.Errorf(
			"se esperaban %d operaciones y hay %d",
			n*m, len(in.Operations),
		)
	}

	if in.LowerBound <= 0 ||
		in.UpperBound <= in.LowerBound {
		return fmt.Errorf("cotas inválida: inferior=%d, superior=%d",
			in.LowerBound, in.UpperBound,
		)
	}

	seenIDs := make([]bool, n*m)
	machineLoad := make([]int, m)
	maxJobLoad := 0
	totalDuration := 0

	for j, job := range in.Jobs {
		if job.ID != j {
			return fmt.Errorf(
				"ID de trabajo incorrecto en posición %d", j,
			)
		}

		if len(job.Operations) != m {
			return fmt.Errorf(
				"trabajo %d: se esperaban %d operaciones", j, m,
			)
		}

		seenMachines := make([]bool, m)
		jobLoad := 0

		for k, op := range job.Operations {
			if op.ID != j*m+k ||
				op.JobID != j ||
				op.Index != k {
				return fmt.Errorf(
					"identidad incorrecta en trabajo %d, operación %d", j, k,
				)
			}

			if op.ID < 0 || op.ID >= n*m || seenIDs[op.ID] {
				return fmt.Errorf(
					"ID de operación inválida o duplicada: %d", op.ID,
				)
			}
			seenIDs[op.ID] = true

			if op.Machine < 0 || op.Machine >= m {
				return fmt.Errorf("operación %d: máquina %d fuera de rango", op.ID, op.Machine)
			}

			if seenMachines[op.Machine] {
				return fmt.Errorf(
					"trabajo %d usa más de una véz la  máquina %d",
					j, op.Machine,
				)
			}
			seenMachines[op.Machine] = true

			if op.Duration <= 0 {
				return fmt.Errorf(
					"operación %d: duración debe ser positiva",
					op.ID,
				)
			}

			jobLoad += op.Duration
			machineLoad[op.Machine] += op.Duration
			totalDuration += op.Duration
		}

		if jobLoad > maxJobLoad {
			maxJobLoad = jobLoad
		}
	}

	maxMachineLoad := 0
	for _, load := range machineLoad {
		if load > maxMachineLoad {
			maxMachineLoad = load
		}
	}

	computedLower := maxJobLoad
	if maxMachineLoad > computedLower {
		computedLower = maxMachineLoad
	}

	/*
		if in.LowerBound != computedLower {
			return fmt.Errorf(
				"cota inferior incorrecta: archivo=%d, calculada=%d",
				in.LowerBound, computedLower,
			)
		}

		if in.UpperBound != totalDuration {
			return fmt.Errorf(
				"cota superior incorrecta: archivo=%d, calculada=%d",
				in.UpperBound, totalDuration,
			)
		}
	*/
	return nil
}
