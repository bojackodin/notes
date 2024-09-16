package speller

import (
	"fmt"
)

type ErrorSpell struct {
	Misspells []Misspell
}

func (e ErrorSpell) Error() string {
	var str string
	lastIndex := len(e.Misspells) - 1
	for i, m := range e.Misspells {
		str += fmt.Sprintf("at %d: %s", m.Pos+1, m.Word)
		// if len(m.Suggestions) > 0 {
		// 	str += fmt.Sprintf(" %v", m.Suggestions)
		// }
		if i != lastIndex {
			str += "; "
		}
	}

	return str
}
