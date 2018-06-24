package ast

import (
	"github.com/amupitan/hero/ast/core"
)

type Loop interface {
	core.Expression
	// evaluate enforces that only ForLoop
	// and RangeLoop implement this interface
	evaluate()
}

// ForLoop structure holds information
// about a parsed AST Loop node
type ForLoop struct {
	core.Statement

	// The name of the loop if it is named
	Name string

	// PreLoop is the expression ran before a loop starts
	PreLoop core.Expression

	// Condition is the condition for evaluation of the loop
	Condition core.Expression

	// PostIteration is the expression ran after every iteration
	PostIteration core.Expression

	// Body is a block for body of the loop
	Body *Block
}

func (l *ForLoop) String() string {
	return `for ` + l.Condition.String() + ` () {}`
}

func (l *ForLoop) evaluate() {}

// RangeLoop represents a for-range loop
type RangeLoop struct {
	// The name of the loop if it is named
	Name string

	// First represents the first identifier in a range for loop
	First string

	// Second represents the first identifier in a range for loop
	Second string

	// Iterable represnets the iterable in for-range loops
	// TODO(DEV) use different type to accommodate for list
	// and map literals
	Iterable string

	// Body is a block for body of the loop
	Body *Block
}

func (r *RangeLoop) String() string {
	return `for ` + r.First + `, ` + r.Second + ` in ` + r.Iterable + ` {}`
}

func (l *RangeLoop) evaluate() {}
