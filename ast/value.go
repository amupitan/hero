package ast

import "github.com/amupitan/hero/ast/core"

type Value struct {
	core.Expression
	Value string
}

func (v *Value) String() string {
	return v.Value
}
