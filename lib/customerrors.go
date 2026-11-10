package customerrors

import (
	"net/http"
	"strconv"
)

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

// ErrWebhookVerificationFailed is returned when an incoming webhook fails
// provider signature or token verification. It is only produced when a
// verification secret is configured for the provider (see the `secrets`
// configuration block); with no secret configured, unsigned webhooks are
// accepted as before.
var ErrWebhookVerificationFailed = NewAppError(http.StatusUnauthorized, "Webhook verification failed")
