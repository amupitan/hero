package ast

type Definition struct {
	Node
}

func (d *Definition) isDefinition() bool { return true }

func (d *Definition) isExpression() bool { return false }
