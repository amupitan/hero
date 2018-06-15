package ast

import "github.com/amupitan/hero/ast/core"

type funcBody struct{}

type param core.Expression

type Function struct {
	core.Definition
	name        string
	parameters  []core.Expression // TODO(DEV) use para,
	returnTypes []interface{}
	body        funcBody
	owner       interface{}
	private     bool
}

func (f *Function) String() string {
	return `func ` + f.name + core.StringifyExpressions(f.parameters)
}
