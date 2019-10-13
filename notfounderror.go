package main

type notFoundError struct {
	msg string
}

// Create a function Error() string and associate it to the struct.
func (error *notFoundError) Error() string {
	return error.msg
}

// Now you can construct an error object using MyError struct.
func newNotFoundError() error {
	return &notFoundError{"404"}
}
