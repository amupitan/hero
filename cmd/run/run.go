package main

import (
	"fmt"

	"github.com/amupitan/hero/lexer"
)

func main() {
	const input = `var a = 's'`
	lex := lexer.New(input)

	tokens, err := lex.Tokenize()
	if err != nil {
		fmt.Println(err)
	}

	for _, token := range tokens {
		fmt.Println(token)

	}
}
