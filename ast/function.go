package ast

import "github.com/amupitan/hero/ast/core"

type funcBody struct{}

type Param struct {
	Type  interface{}
	value interface{}
}

type Function struct {
	core.Definition
	name        string
	parameters  []Param
	returnTypes []interface{}
	body        funcBody
	owner       interface{}
	private     bool
}
