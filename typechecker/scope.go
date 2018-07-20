package typechecker

import (
	"fmt"

	lx "github.com/amupitan/hero/lexer"
)

type Scope struct {
	global *Global
	parent *Scope
	child  *Scope
	vars   map[string]*lx.Token
	funcs  map[string]*lx.Token
}

func (s *Scope) New() *Scope {
	if s.child == nil {
		ns := &Scope{global: s.global, parent: s}
		s.child = ns
		return ns
	}

	return s.child.New()
}

// checkForVar checks if an identifier has been used
func (s *Scope) checkForVar(t *lx.Token) error {
	if c, exists := s.global.typedefs[t.Value]; exists {
		return fmt.Errorf("%s is already declared as a type on line %d:%d (%d, %d)", c.Value, c.Line, c.Column, t.Line, t.Column)
	}

	if v, exists := s.vars[t.Value]; exists {
		return fmt.Errorf("%s is already declared on line %d:%d (%d, %d)", v.Value, v.Line, v.Column, t.Line, t.Column)
	}

	if v, exists := s.funcs[t.Value]; exists {
		return fmt.Errorf("%s is already declared on line %d:%d", v.Value, t.Line, t.Column)
	}

	return nil
}

// AddVar adds a variable to the scope
func (s *Scope) AddVar(t *lx.Token, ttype string) error {
	if err := s.checkForVar(t); err != nil {
		return err
	}
	// if err := s.global.checkForType(t *lx.Token)
	s.vars[t.Value] = t
	return nil
}

// AddFunc adds a function name to the scope
func (s *Scope) AddFunc(t *lx.Token) error {
	if err := s.checkForVar(t); err != nil {
		return err
	}

	s.funcs[t.Value] = t
	return nil
}

func (s *Scope) clear() {
	s.parent = nil
	s.child = nil
	s.vars = nil
	s.funcs = nil
}

// Lookup returns true if the variable is present
func (s *Scope) Lookup(t *lx.Token) bool {
	if err := s.checkForVar(t); err == nil {
		return true
	}

	if s.parent != nil {
		return s.parent.Lookup(t)
	}

	return false
}

// Delete removes a scope and returns it's parent
func (s *Scope) Delete() *Scope {
	p := s.parent
	s.clear()
	p.child = nil
	s = nil //TODO(DEV) is this safe
	return p
}
