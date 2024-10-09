package compiler_test

import (
	"learncompiler/compiler"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompiler(t *testing.T) {
	buf, err := os.ReadFile("../test_files/code_test")
	require.NoError(t, err)

	tokenizer := compiler.NewTokenizer(string(buf))
	tokens := tokenizer.Tokenize()
	t.Fatal(tokens)
	require.NotEmpty(t, tokens)
}
