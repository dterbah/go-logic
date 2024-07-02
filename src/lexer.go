package logic

import (
	"unicode"

	"github.com/dterbah/gods/list"
	"github.com/dterbah/gods/list/arraylist"
)

type TokenType int

// Alphabet of tokens
const (
	ILLEGAL TokenType = iota
	EOF
	VAR    // variable
	AND    // &
	OR     // |
	NOT    // !
	LPAREN // (
	RPAREN // )
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
		return 0
	}

	list := arraylist.New(comparator)
	return &Lexer{input: input, tokens: list}
}

/*
Parse the input of the Lexer and create Token for each elements
*/
func (lexer *Lexer) Tokenize() list.List[Token] {
	for lexer.pos < len(lexer.input) {
		char := lexer.input[lexer.pos]

		if unicode.IsSpace(rune(char)) {
			lexer.pos++
			continue
		}

		switch {
		case char == '&':
			lexer.tokens.Add(Token{Type: AND, Value: "AND"})
			lexer.pos++
		case char == '|':
			lexer.tokens.Add(Token{Type: OR, Value: "OR"})
			lexer.pos++
		case char == '!':
			lexer.tokens.Add(Token{Type: NOT, Value: "NOT"})
			lexer.pos++
		case char == '(':
			lexer.tokens.Add(Token{Type: RPAREN, Value: "("})
			lexer.pos++
		case char == ')':
			lexer.tokens.Add(Token{Type: LPAREN, Value: ")"})
		default:
			if unicode.IsLetter(rune(char)) {
				start := lexer.pos
				for lexer.pos < len(lexer.input) && unicode.IsLetter(rune(lexer.input[lexer.pos])) {
					lexer.pos++
				}
				lexer.tokens.Add(Token{Type: VAR, Value: lexer.input[start:lexer.pos]})
			} else {
				lexer.tokens.Add(Token{Type: ILLEGAL, Value: string(char)})
				lexer.pos++
			}
		}
	}

	return lexer.tokens
}

/*
Return the string representation of the token
*/
func (token Token) String() string {
	return token.Value
}
