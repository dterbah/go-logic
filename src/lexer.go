package logic

import (
	"fmt"
	"unicode"

	"github.com/dterbah/gods/list"
	"github.com/dterbah/gods/list/arraylist"
)

type TokenType int

// Alphabet of tokens
const (
	ILLEGAL     TokenType = iota
	EOF                   // When no token available
	VAR                   // variable
	AND                   // &, ., ^
	OR                    // |, v
	XOR                   // +
	NOT                   // !
	LPAREN                // (
	RPAREN                // )
	IMPLIES               // ->
	NUMBER                // 1, 0
	EQUIVALENCE           // <->
)

// Defines the Token struct
type Token struct {
	Type  TokenType // Type associated to the token
	Value string    // The associated value
}

// Defines the Lexer struct
type Lexer struct {
	input  string
	pos    int
	tokens *arraylist.ArrayList[Token]
}

/*
Create a new Lexer
*/
func NewLexer(input string) *Lexer {
	comparator := func(a, b Token) int {
		if a.Type == b.Type && a.Value == b.Value {
			return 0
		}

		return -1
	}

	list := arraylist.New(comparator)
	return &Lexer{input: input, tokens: list}
}

/*
Parse the input of the Lexer and create Token for each elements
*/
func (lexer *Lexer) Tokenize() (list.List[Token], error) {
	if len(lexer.input) == 0 {
		return nil, fmt.Errorf("the input should not be empty")
	}

	for lexer.pos < len(lexer.input) {
		char := lexer.input[lexer.pos]

		switch {
		case unicode.IsSpace(rune(char)):
			lexer.pos++
			continue
		case isAndOperator(char):
			lexer.tokens.Add(Token{Type: AND, Value: "AND"})
		case isOrOperator(char):
			lexer.tokens.Add(Token{Type: OR, Value: "OR"})
		case isNotOperator(char):
			lexer.tokens.Add(Token{Type: NOT, Value: "NOT"})
		case char == '(':
			lexer.tokens.Add(Token{Type: LPAREN, Value: "("})
		case char == ')':
			lexer.tokens.Add(Token{Type: RPAREN, Value: ")"})
		case isXOROperator(char):
			lexer.tokens.Add(Token{Type: XOR, Value: "XOR"})
		case isNumber(char):
			lexer.tokens.Add(Token{Type: NUMBER, Value: string(char)})
		case char == '<':
			if lexer.pos+2 >= len(lexer.input) {
				return nil, fmt.Errorf("erorr when analyzing equivalence operator")
			}
			dash := lexer.input[lexer.pos+1]
			rightArrow := lexer.input[lexer.pos+2]

			if dash != '-' {
				return nil, fmt.Errorf("error when analyzing equivalence operator, found %s, expected '-'", string(dash))
			}

			if rightArrow != '>' {
				return nil, fmt.Errorf("error when analyzing equivalence operator, found %s, expected '>'", string(rightArrow))
			}
			lexer.pos += 2
			lexer.tokens.Add(Token{Type: EQUIVALENCE, Value: "<->"})

		case char == '-':
			lexer.pos++
			if lexer.pos >= len(lexer.input) || lexer.input[lexer.pos] != '>' {
				return nil, fmt.Errorf("error when analyzing implies operator, found %s, expected '>'", string(lexer.input[lexer.pos]))
			}
			lexer.tokens.Add(Token{Type: IMPLIES, Value: "->"})
		case unicode.IsLetter(rune(char)):
			lexer.tokens.Add(Token{Type: VAR, Value: string(char)})
		default:
			return nil, fmt.Errorf("error when analyzing the char %s", string(char))
		}

		lexer.pos++
	}

	return lexer.tokens, nil
}

// Private functions
func isAndOperator(char byte) bool {
	return char == '&' || char == '.' || char == '^'
}

func isOrOperator(char byte) bool {
	return char == '|' || char == 'v'
}

func isNotOperator(char byte) bool {
	return char == '!'
}

func isXOROperator(char byte) bool {
	return char == '+'
}

func isNumber(char byte) bool {
	return char == '0' || char == '1'
}

/*
Return the string representation of the token
*/
func (token Token) String() string {
	return token.Value
}

/*
Return true if the current token has the type passed in parameter, else false
*/
func (token Token) Is(tokenType TokenType) bool {
	return token.Type == tokenType
}
