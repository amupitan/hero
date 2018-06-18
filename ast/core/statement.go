package core

import "strings"

type Statement interface {
	// string value for debugging purposes
	String() string
}

func StringifyStatements(exps []Statement) string {
	s := strings.Builder{}
	for i, exp := range exps {
		s.WriteString(exp.String())

		// write comma if not last exp
		if i+1 < len(exps) {
			s.WriteString(`, `)
		}
	}

	return s.String()
}
