package compiler

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type Node interface {
	isNode()
}

type DefNode struct {
	Name     string
	ArgNames []string
	Body     Node
}

func (DefNode) isNode() {}

type BodyInt int

func (body BodyInt) isNode() {}

type BodyFnCall struct {
	Name     string
	ArgExprs []Node
}

func (body BodyFnCall) isNode() {}

type BodyVarRef string

func (BodyVarRef) isNode() {}

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
	for {
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

func (parser *Parser) Parse() Node {
	return parser.parseDef()
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

	spew.Dump(defNode)

	return defNode
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

func (parser *Parser) peekOffset(tokenType TokenType, offset int) bool {
	if offset >= len(parser.tokens) {
		return false
	}
	return parser.tokens[offset].Type == tokenType
}

func (parser *Parser) parseExpr() Node {
	if parser.peek(INTEGER) {
		return parser.parseInt()
	}
	if parser.peek(IDENTIFIER) && parser.peekOffset(OPAREN, 1) {
		return parser.parseCall()
	}

	return parser.parseVarRef()
}

func (paser *Parser) parseInt() BodyInt {
	tok := paser.consume(INTEGER)
	val, _ := strconv.ParseInt(tok.Value, 10, 0)
	return BodyInt(val)
}

func (parser *Parser) parseCall() BodyFnCall {
	name := parser.consume(IDENTIFIER)
	argExprs := parser.parseArgExprs()
	return BodyFnCall{
		Name:     name.Value,
		ArgExprs: argExprs,
	}
}

func (parser *Parser) parseArgExprs() []Node {
	argExprs := []Node{}
	parser.consume(OPAREN)

	if !parser.peek(CPAREN) {
		argExprs = append(argExprs, parser.parseExpr())
		for parser.peek(COMMA) {
			parser.consume(COMMA)
			argExprs = append(argExprs, parser.parseExpr())
		}
	}

	parser.consume(CPAREN)

	return argExprs
}

func (parser *Parser) parseVarRef() BodyVarRef {
	return BodyVarRef(parser.consume(IDENTIFIER).Value)
}

func GenerateCode(node Node) string {
	switch val := node.(type) {
	case DefNode:
		return fmt.Sprintf(`function %s(%s) { return %s }`, val.Name, strings.Join(val.ArgNames, ","), GenerateCode(val.Body))
	case BodyFnCall:
		args := make([]string, len(val.ArgExprs))
		for i, expr := range val.ArgExprs {
			args[i] = GenerateCode(expr)
		}
		return fmt.Sprintf("%s(%s)", val.Name, strings.Join(args, ","))
	case BodyInt:
		return fmt.Sprint(val)
	case BodyVarRef:
		return string(val)
	default:
		panic(fmt.Sprintf("unknown type %v", reflect.TypeOf(val)))
	}
}
