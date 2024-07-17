package logic

/* Test autogenerated with the tool go-test-gen. Created 2024-07-03 17:26:57 Wednesday */

import (
	"testing"

	"github.com/dterbah/gods/list/arraylist"
	"github.com/stretchr/testify/assert"
)

func mockComparator(a, b Token) int { return 0 }

func TestNewParser(t *testing.T) {
	assert := assert.New(t)
	parser := NewParser(arraylist.New[Token](mockComparator))

	assert.NotNil(parser)
}

type testCase struct {
	name       string
	expression string
	isError    bool
	variables  map[string]bool
	result     bool
}

func runTestCases(t *testing.T, tests []testCase) {
	assert := assert.New(t)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lexer := NewLexer(test.expression)
			tokens, _ := lexer.Tokenize()
			parser := NewParser(tokens)

			expr, err := parser.Parse()

			if test.isError {
				assert.NotNil(err, test.name)
			} else {
				assert.Nil(err)
				result := expr.Eval(test.variables)
				assert.Equal(result, test.result, test.name)
			}
		})
	}
}

func TestParseNot(t *testing.T) {
	tests := []testCase{
		{"test simple not operation", "!a", false, map[string]bool{"a": true}, false},
		{"test not operator with & after", "!&", true, map[string]bool{"a": true}, false},
		{"test not operator with -> after", "!->", true, map[string]bool{"a": true}, false},
		{"test not operator with | after", "!|", true, map[string]bool{"a": true}, false},
		{"test not operator with ) after", "!)", true, map[string]bool{"a": true}, false},
		{"test simple not operation", "!a", false, map[string]bool{"a": false}, true},
		{"test not operation with parenthesis", "!(avb)", false, map[string]bool{"a": false, "b": false}, true},
		{"test not operation with 1 value", "!1", false, map[string]bool{}, false},
		{"test not operation with 0 value", "!0", false, map[string]bool{}, true},
		{"test not operation with an error after the not", "!(av)", true, map[string]bool{}, true},
	}

	runTestCases(t, tests)
}

func TestParseAnd(t *testing.T) {
	tests := []testCase{
		{"test simple a and a", "a^a", false, map[string]bool{"a": true}, true},
		{"test bad expression a and a ^", "a^a^", true, map[string]bool{"a": true}, true},
		{"test and operator with & after", "&&", true, map[string]bool{"a": true}, false},
		{"test and operator with -> after", "&->", true, map[string]bool{"a": true}, false},
		{"test and operator with | after", "&|", true, map[string]bool{"a": true}, false},
		{"test and operator with ) after", "&)", true, map[string]bool{"a": true}, false},
		{"test simple a and b", "a^b", false, map[string]bool{"a": false, "b": true}, false},
		{"test simple a and b 2", "a^b", false, map[string]bool{"a": true, "b": true}, true},
		{"test simple !a and b", "!a^b", false, map[string]bool{"a": false, "b": true}, true},
		{"test simple !a and !b", "!a^!b", false, map[string]bool{"a": false, "b": true}, false},
		{"test (a and !b)", "(a^!b)", false, map[string]bool{"a": false, "b": true}, false},
	}

	runTestCases(t, tests)
}

func TestParseOr(t *testing.T) {
	tests := []testCase{
		{"test simple a or a", "ava", false, map[string]bool{"a": true}, true},
		{"test or operator with & after", "v&", true, map[string]bool{"a": true}, false},
		{"test or operator with -> after", "v->", true, map[string]bool{"a": true}, false},
		{"test or operator with | after", "v|", true, map[string]bool{"a": true}, false},
		{"test or operator with ) after", "v)", true, map[string]bool{"a": true}, false},
		{"test simple a or b", "avb", false, map[string]bool{"a": false, "b": true}, true},
		{"test simple a or b 2", "avb", false, map[string]bool{"a": false, "b": false}, false},
		{"test simple !a or b ", "!avb", false, map[string]bool{"a": false, "b": true}, true},
		{"test simple !a or !b", "!av!b", false, map[string]bool{"a": true, "b": true}, false},
		{"test simple !a or !b or (a or !b)", "!av!bv(av!b)", false, map[string]bool{"a": true, "b": true}, true},
	}

	runTestCases(t, tests)
}

func TestParseImplies(t *testing.T) {
	tests := []testCase{
		{"test simple a -> a", "a->a", false, map[string]bool{"a": true}, true},
		{"test bad expression a -> a ->", "a->a->", true, map[string]bool{"a": true}, true},
		{"test -> operator with & after", "->&", true, map[string]bool{"a": true}, false},
		{"test -> operator with -> after", "->->", true, map[string]bool{"a": true}, false},
		{"test -> operator with | after", "->|", true, map[string]bool{"a": true}, false},
		{"test -> operator with ) after", "->)", true, map[string]bool{"a": true}, false},
		{"test simple a -> b", "a->b", false, map[string]bool{"a": false, "b": true}, true},
		{"test simple a -> b 2", "a->b", false, map[string]bool{"a": true, "b": false}, false},
		{"test simple !a -> b ", "!a->b", false, map[string]bool{"a": false, "b": true}, true},
		{"test simple !a -> !b", "!a->!b", false, map[string]bool{"a": true, "b": true}, true},
		{"test simple !a -> !b -> (a or !b)", "!av!bv(av!b)", false, map[string]bool{"a": true, "b": true}, true},
	}

	runTestCases(t, tests)
}

func TestParserXOR(t *testing.T) {
	tests := []testCase{
		{"test simple a + a", "a+a", false, map[string]bool{"a": true}, false},
		{"test bad expression a + a +", "a+a+", true, map[string]bool{"a": true}, false},
		{"test + operator with & after", "+&", true, map[string]bool{"a": true}, false},
		{"test + operator with + after", "+", true, map[string]bool{"a": true}, false},
		{"test + operator with | after", "+|", true, map[string]bool{"a": true}, false},
		{"test + operator with ) after", "+)", true, map[string]bool{"a": true}, false},
		{"test simple a + b", "a+b", false, map[string]bool{"a": false, "b": true}, true},
		{"test simple a + b 2", "a+b", false, map[string]bool{"a": true, "b": true}, false},
		{"test simple !a + b ", "!a+b", false, map[string]bool{"a": false, "b": true}, false},
		{"test simple !a + !b", "!a+!b", false, map[string]bool{"a": true, "b": false}, true},
		{"test simple !a + !b + (a or !b)", "!a+!b+(av!b)", false, map[string]bool{"a": true, "b": true}, true},
	}

	runTestCases(t, tests)
}

func TestParserExpression(t *testing.T) {
	tests := []testCase{
		{"test expression (1)", "(1)", false, map[string]bool{}, true},
		{"test expression !(a)", "!(a)", false, map[string]bool{"a": false}, true},
		{"test expression !((a))", "!((a))", false, map[string]bool{"a": false}, true},
		{"test expression (!a)", "(!a)", false, map[string]bool{"a": false}, true},
		{"test expression without closing right parenthesis !((a)", "!((a)", true, map[string]bool{}, true},
		{"test expression with invalid char after a left parenthesis", "(+)", true, map[string]bool{}, true},
	}

	runTestCases(t, tests)
}

func TestParserEquivalence(t *testing.T) {
	tests := []testCase{
		{"test 1<->1", "1<->1", false, map[string]bool{}, true},
	}

	runTestCases(t, tests)
}

func TestParserNumber(t *testing.T) {
	tests := []testCase{
		{"test simple 0 + 1", "0+1", false, map[string]bool{}, true},
		{"test 1 operator with 1 after", "11", true, map[string]bool{}, false},
		{"test 1 operator with 0 after", "10", true, map[string]bool{"a": true}, false},
	}

	runTestCases(t, tests)
}
