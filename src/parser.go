package logic

import (
	"github.com/dterbah/gods/list"
)

// Parser struct used to parse a boolean expression
type Parser struct {
	tokens list.List[Token]
	pos    int
}

func NewParser(tokens list.List[Token]) *Parser {
	return &Parser{tokens: tokens}
}

func (parser *Parser) Parse() Expression {
	return parser.ParseExpression()
}

func (parser *Parser) ParseExpression() Expression {
	var expr Expression

	token := parser.peekToken()

	if token.Is(NOT) {
		parser.pos++
		return NewNotExpression(parser.ParseExpression())
	} else if token.Is(VAR) {
		expr = NewVarExpression(token.Value)
		parser.pos++
		nextToken := parser.peekToken()

		if nextToken.Is(EOF) {
			return expr
		} else if nextToken.Is(OR) {
			parser.pos++
			return NewOrExpression(expr, parser.ParseExpression())
		}
	}

	return expr
}

func (parser *Parser) peekToken() Token {
	if parser.pos >= parser.tokens.Size() {
		return Token{Type: EOF}
	}

	token, _ := parser.tokens.At(parser.pos)

	return token
}
