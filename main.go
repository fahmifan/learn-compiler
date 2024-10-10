package main

import (
	"fmt"
	"learncompiler/compiler"
	"os"
)

func main() {
	buf, err := os.ReadFile("./test_files/code_test")
	noError(err)

	tokenizer := compiler.NewTokenizer(string(buf))
	tokens := tokenizer.Tokenize()

	parser := compiler.NewParser(tokens)
	node := parser.Parse()

	fmt.Println(compiler.GenerateCode(node))
}

func noError(err error) {
	if err != nil {
		panic(err)
	}
}
