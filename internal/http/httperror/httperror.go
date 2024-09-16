package httperror

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type statusError struct {
	err    error
	status int
}

func (e statusError) Error() string {
	if e.err == nil {
		return "status code: " + strconv.Itoa(e.status)
	}
	return e.err.Error()
}

func (e statusError) Unwrap() error { return e.err }

func WithStatus(status int) error {
	return statusError{
		status: status,
	}
}

func WithStatusError(err error, status int) error {
	return statusError{
		err:    err,
		status: status,
	}
}

func HTTPStatus(err error) int {
	if err == nil {
		return 0
	}

	var statusErr statusError
	if errors.As(err, &statusErr) {
		return statusErr.status
	}
	return http.StatusInternalServerError
}

type errorResponse struct {
	Error string `json:"error,omitempty"`
}

func RespondWithError(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err != "" {
		_ = json.NewEncoder(w).Encode(&errorResponse{Error: err})
	}
}
