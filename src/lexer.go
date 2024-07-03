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
	ILLEGAL TokenType = iota
	EOF               // When no token available
	VAR               // variable
	AND               // &, ., ^
	OR                // |, v, +
	NOT               // !
	LPAREN            // (
	RPAREN            // )
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
	for lexer.pos < len(lexer.input) {
		char := lexer.input[lexer.pos]
		if unicode.IsSpace(rune(char)) {
			lexer.pos++
			continue
		}

		if char == '&' || char == '.' || char == '^' {
			lexer.tokens.Add(Token{Type: AND, Value: "AND"})
			lexer.pos++
		} else if char == '|' || char == '+' || char == 'v' { // todo, v is not working
			lexer.tokens.Add(Token{Type: OR, Value: "OR"})
			lexer.pos++
		} else if char == '!' {
			lexer.tokens.Add(Token{Type: NOT, Value: "NOT"})
			lexer.pos++
		} else if char == '(' {
			lexer.tokens.Add(Token{Type: LPAREN, Value: "("})
			lexer.pos++
		} else if char == ')' {
			lexer.tokens.Add(Token{Type: RPAREN, Value: ")"})
			lexer.pos++
		} else {
			if unicode.IsLetter(rune(char)) {
				lexer.tokens.Add(Token{Type: VAR, Value: string(char)})
				lexer.pos++
			} else {
				// error, char not found
				return nil, fmt.Errorf("error when analysing the char %s", string(char))
			}
		}
	}

	return lexer.tokens, nil
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
