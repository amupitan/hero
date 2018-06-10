package core

type Node interface {
	isDefinition() bool
	isExpression() bool
}
