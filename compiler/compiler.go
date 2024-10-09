package compiler

import (
	"fmt"
	"regexp"
	"strings"
)

type Tokenizer struct {
	code string
}

func NewTokenizer(code string) Tokenizer {
	return Tokenizer{
		code: code,
	}
}

type Token struct {
	Value string
	Type  TokenType
}

type TokenType string

const (
	DEF        TokenType = "DEF"
	END        TokenType = "END"
	IDENTIFIER TokenType = "IDENTIFIER"
	INTEGER    TokenType = "INTEGER"
	OPAREN     TokenType = "OPAREN"
	CPAREN     TokenType = "CPAREN"
)

var tokenTypeRegex = map[TokenType]string{
	DEF:        `\bdef\b`,
	END:        `\bend\b`,
	IDENTIFIER: `\b[a-zA-Z]+\b`,
	INTEGER:    `\b[0-9]+\b`,
	OPAREN:     `\(`,
	CPAREN:     `\)`,
}

var tokenTypes = []TokenType{
	DEF,
	END,
	IDENTIFIER,
	INTEGER,
	OPAREN,
	CPAREN,
}

func (tkz *Tokenizer) Tokenize() (tokens []Token) {
	for i := 0; i < 100; i++ {
		isCodeEmpty := len(strings.TrimSpace(tkz.code)) == 0
		if isCodeEmpty {
			return
		}

		token, ok := tkz.tokenizeOne()
		if !ok {
			panic(fmt.Sprintf("Couldnt match token on %v", tkz.code))
		}
		tkz.code = strings.TrimSpace(tkz.code)
		tokens = append(tokens, token)
	}

	return
}

func (tkz *Tokenizer) tokenizeOne() (Token, bool) {
	for _, tokType := range tokenTypes {
		rePattern := fmt.Sprintf(`\A%s`, tokenTypeRegex[tokType])
		re := regexp.MustCompile(rePattern)
		if !re.MatchString(tkz.code) {
			continue
		}

		val := re.FindString(tkz.code)
		tkz.code = tkz.code[len(val):]

		return Token{
			Value: val,
			Type:  tokType,
		}, true
	}

	return Token{}, false
}
