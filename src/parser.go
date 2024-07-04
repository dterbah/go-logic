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

/*
Return the expression associated to the tokens
*/
func (parser *Parser) Parse() (Expression, error) {
	token := parser.peekToken()

	if token.Is(VAR) {
		return parser.parseVar(token)
	} else if token.Is(NOT) {
		return parser.parseNot()
	} else if token.Is(LPAREN) {
		return parser.parseExpression()
	} else {
		return nil, fmt.Errorf("a logic expresion should not begin with a %s", token.Value)
	}
}

func (parser *Parser) parseVar(token Token) (Expression, error) {
	// either OR, AND or XOR
	expr := NewVarExpression(token.Value)
	parser.pos++

	nextToken := parser.peekToken()
	if nextToken.Is(EOF) {
		return expr, nil
	} else if nextToken.Is(AND) {
		return parser.parseAnd(expr)
	} else if nextToken.Is(RPAREN) {
		return expr, nil
	} else {
		return expr, fmt.Errorf("you should not have a %s after a variable", nextToken.Value)
	}
}

func (parser *Parser) parseNot() (Expression, error) {
	parser.pos++
	nextToken := parser.peekToken()

	if nextToken.Is(VAR) {
		expr, err := parser.parseVar(nextToken)
		return NewNotExpression(expr), err
	} else if nextToken.Is(LPAREN) {
		expr, err := parser.parseExpression()
		return NewNotExpression(expr), err
	} else if nextToken.Is(EOF) {
		return nil, fmt.Errorf("you sould have either a variable, or a ( after a not operator")
	} else {
		return nil, fmt.Errorf("you should not have a %s after a not operator", nextToken.Value)
	}
}

func (parser *Parser) parseAnd(left Expression) (Expression, error) {
	parser.pos++
	nextToken := parser.peekToken()

	if nextToken.Is(VAR) {
		expr, err := parser.parseVar(nextToken)
		return NewAndExpression(left, expr), err
	} else if nextToken.Is(NOT) {
		expr, err := parser.parseNot()
		return NewAndExpression(left, expr), err
	} else if nextToken.Is(EOF) {
		return nil, fmt.Errorf("you should have either a variable, a !, or a ( after a and operator")
	} else {
		return nil, fmt.Errorf("you should not have a %s after a and operator", nextToken.Value)
	}
}

func (parser *Parser) parseExpression() (Expression, error) {
	parser.pos++
	nextToken := parser.peekToken()
	var expr Expression
	var err error

	if nextToken.Is(VAR) {
		expr, err = parser.parseVar(nextToken)
	} else if nextToken.Is(NOT) {
		expr, err = parser.parseNot()
	} else if nextToken.Is(EOF) {
		return nil, fmt.Errorf("you should have a variable or a ! after a (")
	} else {
		return nil, fmt.Errorf("you should not have a %s after a (", nextToken.Value)
	}

	nextToken = parser.peekToken()

	if !nextToken.Is(RPAREN) {
		return nil, fmt.Errorf("you should close your expression when you use parenthesis")
	}

	return expr, err
}

// func (parser *Parser) parseExpression() Expression {
// 	var expr Expression

// 	token := parser.peekToken()

// 	if token.Is(NOT) {
// 		parser.pos++
// 		return NewNotExpression(parser.parseExpression())
// 	} else if token.Is(VAR) {
// 		expr = NewVarExpression(token.Value)
// 		parser.pos++
// 		nextToken := parser.peekToken()

// 		if nextToken.Is(EOF) {
// 			return expr
// 		} else if nextToken.Is(OR) {
// 			parser.pos++
// 			return NewOrExpression(expr, parser.parseExpression())
// 		} else if nextToken.Is(AND) {
// 			parser.pos++
// 			return NewAndExpression(expr, parser.parseExpression())
// 		}
// 	} else if token.Is(LPAREN) {
// 		parser.pos++
// 		expr = parser.parseExpression()
// 		nextToken := parser.peekToken()
// 		if !nextToken.Is(RPAREN) {
// 			panic(fmt.Sprintf("error when analyzing the input. Received %s, waiting %s", nextToken.Value, ")"))
// 		}

// 		parser.pos++
// 		nextToken = parser.peekToken()
// 		if nextToken.Is(EOF) {
// 			return expr
// 		} else if nextToken.Is(OR) {
// 			parser.pos++
// 			return NewOrExpression(expr, parser.parseExpression())
// 		} else if nextToken.Is(AND) {
// 			parser.pos++
// 			return NewAndExpression(expr, parser.parseExpression())
// 		}
// 	}

// 	return expr
// }

func (parser *Parser) peekToken() Token {
	if parser.pos >= parser.tokens.Size() {
		return Token{Type: EOF}
	}

	token, _ := parser.tokens.At(parser.pos)

	return token
}
