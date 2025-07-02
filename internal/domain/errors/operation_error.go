package errors

type OperationError struct {
	message string
}

func NewOperationError(message string) error {
	return &OperationError{message: message}
}

func (e *OperationError) Error() string {
	return e.message
}
