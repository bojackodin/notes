package speller

import (
	"fmt"
)

type SpellError struct {
	Misspells []Misspell
}

func (e *SpellError) Error() string {
	var str string
	lastIndex := len(e.Misspells) - 1
	for i, m := range e.Misspells {
		str += fmt.Sprintf("at %d: %s", m.Pos+1, m.Word)
		if i != lastIndex {
			str += "; "
		}
	}

	return str
}
