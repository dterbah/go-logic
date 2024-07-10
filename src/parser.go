package logic

import (
	"fmt"
	"strconv"

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

// --- State machine --- //

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
	} else if token.Is(NUMBER) {
		return parser.parseNumber(token)
	} else {
		return nil, fmt.Errorf("a logic expression should not begin with a %s", token.Value)
	}
}

/*
Parse VAR boolean expression
*/
func (parser *Parser) parseVar(token Token) (Expression, error) {
	expr := NewVarExpression(token.Value)
	parser.pos++
	return parser.parseOperator(expr)
}

func (parser *Parser) parseNumber(token Token) (Expression, error) {
	value, _ := strconv.Atoi(token.Value)
	expr := NewNumberExpression(value)
	parser.pos++
	return parser.parseOperator(expr)
}

/*
Parse NOT boolean expression
*/
func (parser *Parser) parseNot() (Expression, error) {
	parser.pos++
	nextToken := parser.peekToken()

	var expr Expression
	var err error

	if nextToken.Is(VAR) {
		expr = NewNotExpression(NewVarExpression(nextToken.Value))
		parser.pos++
	} else if nextToken.Is(LPAREN) {
		expr, err = parser.parseExpression()
		expr = NewNotExpression(expr)
	} else if nextToken.Is(NUMBER) {
		value, _ := strconv.Atoi(nextToken.Value)
		expr = NewNotExpression(NewNumberExpression(value))
		parser.pos++
	} else {
		return nil, fmt.Errorf("you should have a variable or a ( after a not operator")
	}

	if err != nil {
		return nil, err
	}

	return parser.parseOperator(expr)
}

/*
Parse boolean expression between parenthesis
*/
func (parser *Parser) parseExpression() (Expression, error) {
	parser.pos++
	nextToken := parser.peekToken()

	var expr Expression
	var err error

	if nextToken.Is(VAR) {
		expr, err = parser.parseVar(nextToken)
	} else if nextToken.Is(NOT) {
		expr, err = parser.parseNot()
	} else if nextToken.Is(NUMBER) {
		expr, err = parser.parseNumber(nextToken)
	} else if nextToken.Is(LPAREN) {
		expr, err = parser.parseExpression()
	} else {
		return nil, fmt.Errorf("you should have a variable or a ! after a (")
	}

	if err != nil {
		return nil, err
	}

	nextToken = parser.peekToken()
	if !nextToken.Is(RPAREN) {
		return nil, fmt.Errorf("you should close your expression when you use parenthesis")
	}

	parser.pos++
	return parser.parseOperator(expr)
}

/*
Parse the operator following an expression
*/
func (parser *Parser) parseOperator(left Expression) (Expression, error) {
	for {
		nextToken := parser.peekToken()

		switch {
		case nextToken.Is(EOF) || nextToken.Is(RPAREN):
			return left, nil
		case nextToken.Is(AND):
			parser.pos++
			right, err := parser.Parse()
			if err != nil {
				return nil, err
			}
			left = NewAndExpression(left, right)
		case nextToken.Is(OR):
			parser.pos++
			right, err := parser.Parse()
			if err != nil {
				return nil, err
			}
			left = NewOrExpression(left, right)
		case nextToken.Is(IMPLIES):
			parser.pos++
			right, err := parser.Parse()
			if err != nil {
				return nil, err
			}
			left = NewImpliesExpression(left, right)
		case nextToken.Is(XOR):
			parser.pos++
			right, err := parser.Parse()
			if err != nil {
				return nil, err
			}
			left = NewXORExpression(left, right)
		default:
			return nil, fmt.Errorf("unexpected token %s", nextToken.Value)
		}
	}
}

func (parser *Parser) peekToken() Token {
	if parser.pos >= parser.tokens.Size() {
		return Token{Type: EOF}
	}

	token, _ := parser.tokens.At(parser.pos)
	return token
}
