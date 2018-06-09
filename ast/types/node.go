package ast

type Node interface {
	isDefinition() bool
	isExpression() bool
}
