package error

type DisError struct {
	err string //error description
}

func (e *DisError) Error() string {
	return e.err
}