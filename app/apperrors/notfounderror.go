package apperrors

// NotFoundError struct which uses for generating '404' answer
type NotFoundError struct {
	msg string
}

// Error() returns a string value of 'msg' field.
func (error *NotFoundError) Error() string {
	return error.msg
}

// NewNotFoundError creates a new NotFoundError object.
func NewNotFoundError() error {
	return &NotFoundError{"404"}
}
