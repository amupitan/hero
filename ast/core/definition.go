package core

type Definition struct {
	Statement
}

func (d *Definition) isDefinition() bool { return true }

func (d *Definition) isExpression() bool { return false }
