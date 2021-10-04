package customerrors

import "strconv"

type AppError struct {
	StatusCode int
	ErrorText  string
}

// The AppError.Error is the generic custom error for the Webhook Go app.
func (ae AppError) Error() string {
	return "status code: " + strconv.Itoa(ae.StatusCode) + ", error message: " + ae.ErrorText
}

// The NewAppError creates a new AppError from an http statuscode and text string.
func NewAppError(statusCode int, errorText string) AppError {
	return AppError{statusCode, errorText}
}
