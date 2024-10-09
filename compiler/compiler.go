package compiler

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type DefNode struct {
	Name     string
	ArgNames []string
	Body     BodyNode
}

type BodyInt int

func (body BodyInt) Value() string {
	return fmt.Sprint(body)
}

type BodyNode interface {
	Value() string
}

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
	COMMA      TokenType = "COMMA"
)

var tokenTypeRegex = map[TokenType]string{
	DEF:        `\bdef\b`,
	END:        `\bend\b`,
	IDENTIFIER: `\b[a-zA-Z]+\b`,
	INTEGER:    `\b[0-9]+\b`,
	OPAREN:     `\(`,
	CPAREN:     `\)`,
	COMMA:      `,`,
}

var tokenTypes = []TokenType{
	DEF,
	END,
	IDENTIFIER,
	INTEGER,
	OPAREN,
	CPAREN,
	COMMA,
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

type Parser struct {
	tokens []Token
}

func NewParser(tokens []Token) Parser {
	return Parser{tokens: tokens}
}

func (parser *Parser) Parse() {
	parser.parseDef()
}

func (parser *Parser) parseDef() DefNode {
	parser.consume(DEF)
	name := parser.consume(IDENTIFIER)
	argNames := parser.parserArgNames()
	body := parser.parseExpr()
	parser.consume(END)

	defNode := DefNode{
		Name:     name.Value,
		ArgNames: argNames,
		Body:     body,
	}

	fmt.Println(defNode)

	return defNode
}

func valuesFromTokens(tokens []Token) []string {
	vals := make([]string, len(tokens))
	for i, tok := range tokens {
		vals[i] = tok.Value
	}

	return vals
}

func (parser *Parser) consume(tokenType TokenType) Token {
	token := parser.tokens[0]
	parser.tokens = parser.tokens[1:]
	if token.Type == tokenType {
		return token
	}

	panic(fmt.Sprintf("Expected token_type %s but got %s", tokenType, token.Type))
}

func (parser *Parser) parserArgNames() []string {
	argNames := []string{}
	parser.consume(OPAREN)
	if parser.peek(IDENTIFIER) {
		argNames = append(argNames, parser.consume(IDENTIFIER).Value)
		for parser.peek(COMMA) {
			parser.consume(COMMA)
			argNames = append(argNames, parser.consume(IDENTIFIER).Value)
		}
	}
	parser.consume(CPAREN)

	return argNames
}

func (parser *Parser) peek(tokenType TokenType) bool {
	return parser.tokens[0].Type == tokenType
}

func (parser *Parser) parseExpr() BodyNode {
	return parser.parseInt()
}

func (paser *Parser) parseInt() BodyInt {
	tok := paser.consume(INTEGER)
	val, _ := strconv.ParseInt(tok.Value, 10, 0)
	return BodyInt(val)
}
