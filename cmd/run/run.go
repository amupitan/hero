package main

import (
	"fmt"

	"github.com/amupitan/hero/lexer"
)

func main() {
	const input = "a == 3"
	lex := lexer.New(input)
	fmt.Println(lex.NextToken())
	fmt.Println(lex.NextToken())
	fmt.Println(lex.NextToken())
	fmt.Println(lex.NextToken())
	fmt.Println(lex.NextToken())
	fmt.Println(lex.NextToken())
}
