package typechecker

import (
	"fmt"

	"github.com/amupitan/hero/ast"
	lx "github.com/amupitan/hero/lexer"
)

type Env struct {
	*ast.Program
	*Global
	// curr is the current scope
	curr *Scope
}

func New(p *ast.Program) *Env {
	e := &Env{
		Program: p,
		Global:  NewGlobal(),
	}
	e.curr = e.Scope
	return e
}

type Signature struct {
	params      []int
	returnTypes []int
}

type Global struct {
	*Scope
	Signatures []Signature
	typedefs   map[string]*lx.Token
	types      TypeEnv
}

type TypeEnv map[string]string

func NewGlobal() *Global {
	g := &Global{
		Scope:    &Scope{},
		types:    make(TypeEnv),
		typedefs: make(map[string]*lx.Token),
	}

	g.Scope.global = g
	return g
}

// NewScope creates a new scope
func (g *Global) NewScope() *Scope {
	return g.Scope.New()
}

// AddType adds a type to the global type definitions
func (g *Global) AddType(t *lx.Token) error {
	if isBuiltinToken(t) {
		return fmt.Errorf("Cannot create type with name %s on line %d:%d, %s is a builtin type",
			t.Value, t.Line, t.Column, t.Value)
	}
	if c, exists := g.typedefs[t.Value]; exists {
		return fmt.Errorf("%s is already declared as a type on line %d:%d", c.Value, t.Line, t.Column)
	}

	g.typedefs[t.Value] = t
	return nil
}

// hasType returns true if a type with the name exists
func (g *Global) hasType(t string) bool {
	if isBuiltin(t) {
		return true
	}
	_, exists := g.typedefs[t]
	return exists
}

// hasType returns true if a type with the name exists
func (g *Global) getVarType(t string) string {
	return g.types[t]
}

// checkForType checks if a type exists
// and returns an error if the type does not exist
func (g *Global) checkForType(t *lx.Token) error {
	if g.hasType(t.Value) {
		return nil
	}

	return fmt.Errorf("type: %s is not defined (%d:, %d)", t.Value, t.Line, t.Column)
}

// isBuiltin returns true if the name represents a builtin type
func isBuiltin(t string) bool {
	return t == `int` || t == `string` || t == `bool`
}

func isBuiltinToken(t *lx.Token) bool {
	return t.Type == lx.Identifier && isBuiltin(t.Value)
}
