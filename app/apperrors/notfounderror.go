package apperrors

type NotFoundError struct {
	msg string
}

// Create a function Error() string and associate it to the struct.
func (error *NotFoundError) Error() string {
	return error.msg
}

// Now you can construct an error object using MyError struct.
func NewNotFoundError() error {
	return &NotFoundError{"404"}
}
