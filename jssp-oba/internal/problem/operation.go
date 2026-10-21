package problem

// Operation representa una operación de un trabajo
type Operation struct {
	ID       int
	JobID    int
	Index    int
	Machine  int
	Duration int
}

// Job contine las operaciones en su orden de precendecia
type Job struct {
	ID         int
	Operations []Operation
}
