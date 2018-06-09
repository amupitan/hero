package ast

type Expression struct {
	Node
}

func (e *Expression) isDefinition() bool { return false }

func (e *Expression) isExpression() bool { return true }
