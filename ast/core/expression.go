package core

import "strings"

type Expression interface {
	Statement
	Type() ExpressionType
}

type ExpressionType int16

const (
	Nil ExpressionType = iota
	AssignmentNode
	AtomNode
	BinaryNode
	CallNode
	FunctionNode
	OperationNode
	StringNode
	ValueNode
)

func StringifyExpressions(exps []Expression) string {
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
