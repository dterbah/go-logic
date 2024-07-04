package logic

/* Test autogenerated with the tool go-test-gen. Created 2024-07-02 16:56:05 Tuesday */

import (
	"testing"

	"github.com/dterbah/gods/list/arraylist"
	"github.com/stretchr/testify/assert"
)

func TestNewLexer(t *testing.T) {
	assert := assert.New(t)
	lexer := NewLexer("a&b")
	assert.NotNil(lexer)
}

func TestTokenize(t *testing.T) {
	assert := assert.New(t)
	mockTokenCompare := func(a, b Token) int {
		if a.Value == b.Value && a.Type == b.Type {
			return 0
		}

		return -1
	}

	tests := []struct {
		name    string
		tokens  *arraylist.ArrayList[Token]
		input   string
		isError bool
	}{
		{"test with illegal char", arraylist.New(mockTokenCompare, Token{Type: VAR, Value: "a"}, Token{Type: AND, Value: "AND"}, Token{Type: VAR, Value: "b"}), "a=b", true},
		{"test with an AND expression with &", arraylist.New(mockTokenCompare, Token{Type: VAR, Value: "a"}, Token{Type: AND, Value: "AND"}, Token{Type: VAR, Value: "b"}), "a&b", false},
		{"test with an AND expression with .", arraylist.New(mockTokenCompare, Token{Type: VAR, Value: "a"}, Token{Type: AND, Value: "AND"}, Token{Type: VAR, Value: "b"}), "a.b", false},
		{"test with an AND expression with ^", arraylist.New(mockTokenCompare, Token{Type: VAR, Value: "a"}, Token{Type: AND, Value: "AND"}, Token{Type: VAR, Value: "b"}), "a^b", false},
		{"test with an OR expression with |", arraylist.New(mockTokenCompare, Token{Type: VAR, Value: "a"}, Token{Type: OR, Value: "OR"}, Token{Type: VAR, Value: "b"}), "a|b", false},
		{"test with an OR expression with v", arraylist.New(mockTokenCompare, Token{Type: VAR, Value: "a"}, Token{Type: OR, Value: "OR"}, Token{Type: VAR, Value: "b"}), "avb", false},
		{"test with an NOT expression", arraylist.New(mockTokenCompare, Token{Type: NOT, Value: "NOT"}, Token{Type: VAR, Value: "a"}), "!a", false},
		{"test with parenthesis", arraylist.New(mockTokenCompare, Token{Type: LPAREN, Value: "("}, Token{Type: VAR, Value: "a"}, Token{Type: AND, Value: "AND"}, Token{Type: VAR, Value: "b"}, Token{Type: RPAREN, Value: ")"}), "(a&b)", false},
		{"test with an AND expression", arraylist.New(mockTokenCompare, Token{Type: NOT, Value: "NOT"}, Token{Type: VAR, Value: "a"}), "!a", false},
		{"test with spaces", arraylist.New(mockTokenCompare, Token{Type: NOT, Value: "NOT"}, Token{Type: VAR, Value: "a"}), " ! a ", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lexer := NewLexer(test.input)
			tokens, err := lexer.Tokenize()
			if test.isError {
				assert.NotNil(err)
				assert.Nil(tokens)
			} else {
				assert.Nil(err)
				assert.True(test.tokens.ContainsAll(tokens))
			}
		})
	}
}

func TestString(t *testing.T) {
	assert := assert.New(t)
	token := Token{Type: AND, Value: "AND"}

	assert.Equal("AND", token.String())
}
