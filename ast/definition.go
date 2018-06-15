package ast

import (
	"github.com/amupitan/hero/ast/core"
)

// Relationship structure of definitions
// 			Declaration
// 	    	/          \
//         /            \
//     Definition(Exp)  =Class
//          /    \
//         /      \
//    =Variable    =Function
//

type Definition struct {
	core.Expression
	core.Declaration
	Binary
	Name  string
	Value core.Expression
}

func (d *Definition) String() string {
	return `var ` + d.Name + ` = ` + d.Value.String()
}
