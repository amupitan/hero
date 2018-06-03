package main

import (
	"fmt"

	"github.com/amupitan/hero/lexer"
)

func main() {
	const input = "1 + 1"
	lex := lexer.New(input)
	fmt.Println(lex.NextToken())
}
