package ast

import (
	"github.com/amupitan/hero/ast/core"
	"github.com/amupitan/hero/lexer"
)

// Relationship structure of definitions
// 			Declaration
// 	    	/          \
//         /            \
//     Definition  =Class
//          /    \
//         /      \
//    =Variable    =Function    ==== expressuion
//

type Definition struct {
	// core.Expression TODO(??) should definitions be expressions?
	core.Declaration
	Name      lexer.Token
	Value     core.Expression
	LexerType lexer.Token // TODO use lexer or custom ast type for type
}

func (d *Definition) String() string {
	s := `var ` + d.Name.Value + ` ` + string(d.LexerType.Value)
	if d.Value != nil {
		s += ` = ` + d.Value.String()
	}
	return s
}

func (d *Definition) Type() core.ExpressionType {
	return core.Nil
}
