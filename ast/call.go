package ast

import (
	"strings"

	"github.com/amupitan/hero/ast/core"
)

type Call struct {
	Name    string // TODO: take in complete token?
	Args    []core.Expression
	Object  string
	Func    *Function
	Negated bool
}

func (c *Call) String() string {
	s := strings.Builder{}
	// named call
	if c.Name != `` {
		if c.Object != `` {
			s.WriteString(c.Object)
			s.WriteRune('.')
		}
		s.WriteString(c.Name)
	} else {
		// lambda call
		s.WriteString(c.Func.String())
	}

	s.WriteRune('(')
	s.WriteString(core.StringifyExpressions(c.Args))
	s.WriteRune(')')

	return s.String()
}

func (c *Call) Type() core.ExpressionType {
	return core.CallNode
}
