package main

import (
	"fmt"
	"learncompiler/compiler"
	"os"
	"strings"
)

func main() {
	buf, err := os.ReadFile("./test_files/code_test")
	noError(err)

	tokenizer := compiler.NewTokenizer(string(buf))
	tokens := tokenizer.Tokenize()
	fmt.Println(toString(tokens))
}

func toString(tokens []compiler.Token) string {
	ss := make([]string, len(tokens))
	for i, tok := range tokens {
		ss[i] = tok.Value + "__" + string(tok.Type)
	}

	return strings.Join(ss, ",")
}

func noError(err error) {
	if err != nil {
		panic(err)
	}
}
