package errors

import "fmt"

type ErrNotFound struct{}

func (e *ErrNotFound) Error() string {
	return "not found"
}

type ErrApiError struct {
	StatusCode int
	Message    string
}

func (e *ErrApiError) Error() string {
	if e.StatusCode == 0 && e.Message == "" {
		return "unknown api error"
	}

	msg := fmt.Sprintf("call failed with status %d", e.StatusCode)
	if e.Message != "" {
		msg += " " + e.Message
	}

	return msg
}

type ErrInvalidParams struct {
	ParamName string
}

func (e *ErrInvalidParams) Error() string {
	return fmt.Sprintf("missing parameter: %s", e.ParamName)
}

type ErrInvalidToken struct {
	ProviderName string
}

func (e *ErrInvalidToken) Error() string {
	if e.ProviderName != "" {
		return "invalid token for " + e.ProviderName
	}
	return "invalid token"
}
