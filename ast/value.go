package ast

import "github.com/amupitan/hero/ast/core"

type Value struct {
	Value string
}

func (v *Value) String() string {
	return v.Value
}

func (v *Value) Type() core.ExpressionType {
	return core.ValueNode
}
