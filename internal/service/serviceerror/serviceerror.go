package serviceerror

import (
	"errors"
)

var ErrUserDuplicate = errors.New("user duplicate")
