package serviceerror

import (
	"errors"
)

var (
	ErrUserDuplicate = errors.New("user duplicate")
	ErrSpeller       = errors.New("spelling mistakes")
)
