package ast

import "github.com/amupitan/hero/ast/core"

type funcBody struct{}

type param core.Expression

type Function struct {
	core.Expression
	Definition
	parameters  []core.Expression // TODO(DEV) use para,
	returnTypes []interface{}
	body        funcBody
	owner       interface{}
	private     bool
}

func (f *Function) String() string {
	return `func ` + f.Name + core.StringifyExpressions(f.parameters)
}
