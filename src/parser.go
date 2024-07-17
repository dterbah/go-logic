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

// Parse parses the entire expression
func (parser *Parser) Parse() (Expression, error) {
	return parser.parseEquivalence()
}

// parseEquivalence parses equivalence expressions (lowest priority)
func (parser *Parser) parseEquivalence() (Expression, error) {
	left, err := parser.parseImplies()
	if err != nil {
		return nil, err
	}

	for parser.peekToken().Is(EQUIVALENCE) {
		parser.pos++
		right, err := parser.parseImplies()
		if err != nil {
			return nil, err
		}
		left = NewEquivalenceExpression(left, right)
	}

	return left, nil
}

// parseImplies parses implication expressions
func (parser *Parser) parseImplies() (Expression, error) {
	left, err := parser.parseXOR()
	if err != nil {
		return nil, err
	}

	for parser.peekToken().Is(IMPLIES) {
		parser.pos++
		right, err := parser.parseXOR()
		if err != nil {
			return nil, err
		}
		left = NewImpliesExpression(left, right)
	}

	return left, nil
}

// parseXOR parses XOR expressions
func (parser *Parser) parseXOR() (Expression, error) {
	left, err := parser.parseOr()
	if err != nil {
		return nil, err
	}

	for parser.peekToken().Is(XOR) {
		parser.pos++
		right, err := parser.parseOr()
		if err != nil {
			return nil, err
		}
		left = NewXORExpression(left, right)
	}

	return left, nil
}

// parseOr parses OR expressions
func (parser *Parser) parseOr() (Expression, error) {
	left, err := parser.parseAnd()
	if err != nil {
		return nil, err
	}

	for parser.peekToken().Is(OR) {
		parser.pos++
		right, err := parser.parseAnd()
		if err != nil {
			return nil, err
		}
		left = NewOrExpression(left, right)
	}

	return left, nil
}

// parseAnd parses AND expressions
func (parser *Parser) parseAnd() (Expression, error) {
	left, err := parser.parseNot()
	if err != nil {
		return nil, err
	}

	for parser.peekToken().Is(AND) {
		parser.pos++
		right, err := parser.parseNot()
		if err != nil {
			return nil, err
		}
		left = NewAndExpression(left, right)
	}

	return left, nil
}

// parseNot parses NOT expressions
func (parser *Parser) parseNot() (Expression, error) {
	if parser.peekToken().Is(NOT) {
		parser.pos++
		next := parser.peekToken()
		if next.IsOperator() {
			return nil, fmt.Errorf("expected var or number after a not operator")
		}

		expr, err := parser.parseNot()
		if err != nil {
			return nil, err
		}
		return NewNotExpression(expr), nil
	}

	return parser.parsePrimary()
}

// parsePrimary parses primary expressions (variables, numbers, parenthesized expressions)
func (parser *Parser) parsePrimary() (Expression, error) {
	token := parser.peekToken()

	switch {
	case token.Is(VAR):
		parser.pos++
		next := parser.peekToken()
		if !next.IsOperator() && !next.Is(EOF) && !next.Is(RPAREN) {
			return nil, fmt.Errorf("expected operator after a variable")
		}
		return NewVarExpression(token.Value), nil
	case token.Is(NUMBER):
		parser.pos++
		next := parser.peekToken()
		if !next.IsOperator() && !next.Is(EOF) && !next.Is(RPAREN) {
			return nil, fmt.Errorf("expected operator after a number")
		}
		value, _ := strconv.Atoi(token.Value)
		return NewNumberExpression(value), nil
	case token.Is(LPAREN):
		parser.pos++
		expr, err := parser.parseEquivalence()
		if err != nil {
			return nil, err
		}
		if !parser.peekToken().Is(RPAREN) {
			return nil, fmt.Errorf("expected closing parenthesis")
		}
		parser.pos++
		return expr, nil
	default:
		return nil, fmt.Errorf("unexpected token: %s", token.Value)
	}
}

func (parser *Parser) peekToken() Token {
	if parser.pos >= parser.tokens.Size() {
		return Token{Type: EOF}
	}

	token, _ := parser.tokens.At(parser.pos)
	return token
}
