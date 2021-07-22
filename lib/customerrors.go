package customerrors

import "strconv"

type AppError struct {
	StatusCode int
	ErrorText  string
}

func (ae AppError) Error() string {
	return "status code: " + strconv.Itoa(ae.StatusCode) + ", error message: " + ae.ErrorText
}

func NewAppError(statusCode int, errorText string) AppError {
	return AppError{statusCode, errorText}
}
