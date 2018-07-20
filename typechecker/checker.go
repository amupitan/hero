package typechecker

import (
	"fmt"

	"github.com/amupitan/hero/ast"
	lx "github.com/amupitan/hero/lexer"
)

func (e *Env) checkAtom() error {
	return nil
}

// checkFunc performs semantic analysis on a Function
func (e *Env) checkFunc(f *ast.Function) (error, bool) {
	isSig := !f.Lambda
	if err := e.AddFunc(&f.Name); err != nil {
		return err, isSig
	}

	// ensure param types are valid
	for _, p := range f.Parameters {
		if err := e.checkForType(&p.Type); err != nil {
			return err, isSig
		}
	}

	// ensure return types are valid
	for _, r := range f.ReturnTypes {
		if err := e.Global.checkForType(&r); err != nil {
			return err, isSig
		}
	}

	// create scope
	e.curr = e.Scope.New()

	err, ret, breakT, continueT := e.checkBlock(f.Body)
	if err != nil {
		return err, isSig
	}
	if breakT != nil {
		return fmt.Errorf("Unexpected `break` (%d, %d)", breakT.Line, breakT.Column), isSig
	}

	if continueT != nil {
		return fmt.Errorf("Unexpected `continue` (%d, %d)", continueT.Line, continueT.Column), isSig
	}

	if len(f.ReturnTypes) > 0 {
		if ret == nil {
			return fmt.Errorf("Expected return at end of function declared at (%d, %d)", f.Name.Line, f.Name.Column), isSig
		}

		// check for right number of returned vars
		if len(ret.Values) != len(f.ReturnTypes) {
			return fmt.Errorf("Expected %d return values but got %s (%d, %d)", len(f.ReturnTypes), len(ret.Values), ret.Line, ret.Column), isSig
		}

		// check the return vars are the right types
		// TODO(DEV)
		for t, r := range ret.Values {
			if t.Type == lx.Identifier {
				if e.getVarType(t.Value) != e.Scope.types[fr] {
					return fmt.Errorf("%s's type does not match return type %s (%d, %d)", r, e.Scope.types[fr], r.Line, ret.Column), isSig
				}
			} else {
				// check literal type
			}
		}
	} else {
		// no returns were expected
		if ret != nil {
			return fmt.Errorf("Unexpected `continue` (%d, %d)", ret.Line, ret.Column), isSig
		}
	}

	// leave scope
	e.Scope = e.Scope.Delete()
	return nil, isSig
}

func (e *Env) checkBlock(b *ast.Block) (err error, returns *ast.Return, breakT *lx.Token, continueT *lx.Token) {
	stmts := b.Statements
	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *ast.Function:
			var isSig bool
			err, isSig = e.checkFunc(s)
			if e.Scope != e.Global.Scope && isSig {
				err = fmt.Errorf("Cannot create function using signature in non-global context (%d, %d)", s.Name.Line, s.Name.Column)
			}
		}
	}

	// TODO(DEV) what?
	return err, nil, nil, nil
}

func (e *Env) Check() error {
	b := e.Program.Body
	err, ret, breakT, continueT := e.checkBlock(b)
	if ret != nil {
		return fmt.Errorf("Unexpected return to non-return function (%d, %d)", ret.Line, ret.Column)
	}

	if breakT != nil {
		return fmt.Errorf("Unexpected `break` (%d, %d)", breakT.Line, breakT.Column)
	}

	if continueT != nil {
		return fmt.Errorf("Unexpected `continue` (%d, %d)", continueT.Line, continueT.Column)
	}

	return err
}
