package apperrors

import "testing"

func TestNotFoundError(t *testing.T) {
	err := NewNotFoundError()

	expected := "404"
	result := (err).Error()
	if result != expected {
		t.Errorf("NewNotFoundError returned unexpected result: got \n %v \nwant\n %v",
			result, expected)
	}

}
