package repositoryerror

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicate      = errors.New("duplicate")
)
