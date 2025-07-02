package errors

type ForbiddenError struct {
	message string
}

func NewForbiddenError(message string) error {
	return &ForbiddenError{message: message}
}

func (e *ForbiddenError) Error() string {
	return e.message
}
