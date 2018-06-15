package core

type Statement interface {
	isDefinition() bool
	isExpression() bool
}
