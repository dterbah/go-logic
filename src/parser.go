package logic

import (
	"fmt"

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
	return parser.parseExpression()
}

func (parser *Parser) parseExpression() Expression {
	var expr Expression

	token := parser.peekToken()

	if token.Is(NOT) {
		parser.pos++
		return NewNotExpression(parser.parseExpression())
	} else if token.Is(VAR) {
		expr = NewVarExpression(token.Value)
		parser.pos++
		nextToken := parser.peekToken()

		if nextToken.Is(EOF) {
			return expr
		} else if nextToken.Is(OR) {
			parser.pos++
			return NewOrExpression(expr, parser.parseExpression())
		} else if nextToken.Is(AND) {
			parser.pos++
			return NewAndExpression(expr, parser.parseExpression())
		}
	} else if token.Is(LPAREN) {
		parser.pos++
		expr = parser.parseExpression()
		nextToken := parser.peekToken()
		if !nextToken.Is(RPAREN) {
			panic(fmt.Sprintf("error when analyzing the input. Received %s, waiting %s", nextToken.Value, ")"))
		}

		parser.pos++
		nextToken = parser.peekToken()
		if nextToken.Is(EOF) {
			return expr
		} else if nextToken.Is(OR) {
			parser.pos++
			return NewOrExpression(expr, parser.parseExpression())
		} else if nextToken.Is(AND) {
			parser.pos++
			return NewAndExpression(expr, parser.parseExpression())
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
