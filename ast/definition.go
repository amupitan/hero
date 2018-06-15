package ast

import (
	"github.com/amupitan/hero/ast/core"
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
	Name  string
	Value core.Expression
	Type  string // TODO use lexer or custom ast type for type
}

func (d *Definition) String() string {
	return `var ` + d.Name + ` ` + d.Type + ` = ` + d.Value.String()
}
