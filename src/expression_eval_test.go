package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type evalTestCase struct {
	name           string
	expr           Expression
	variables      map[string]bool
	expectedResult bool
}

func runEvalTestCases(t *testing.T, tests []evalTestCase) {
	assert := assert.New(t)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(test.expr.Eval(test.variables), test.expectedResult, test.name)
		})
	}
}

func TestNotExpressionEval(t *testing.T) {
	tests := []evalTestCase{
		{
			"test eval not expression with var inside",
			NewNotExpression(NewVarExpression("a")),
			map[string]bool{"a": true},
			false,
		},
	}

	runEvalTestCases(t, tests)
}

func TestVarExpressionEval(t *testing.T) {
	tests := []evalTestCase{
		{
			"test eval var expression with true result",
			NewVarExpression("a"),
			map[string]bool{"a": true},
			true,
		},
		{
			"test eval var expression with true result",
			NewVarExpression("a"),
			map[string]bool{"a": false},
			false,
		},
	}

	runEvalTestCases(t, tests)
}

func TestOrExpressionEval(t *testing.T) {
	tests := []evalTestCase{
		{
			"test or expression with two variables",
			NewOrExpression(NewVarExpression("a"), NewVarExpression("b")),
			map[string]bool{"a": true, "b": false},
			true,
		},
		{
			"test or expression with one variable and one number",
			NewOrExpression(NewVarExpression("a"), NewNumberExpression(0)),
			map[string]bool{"a": true},
			true,
		},
	}

	runEvalTestCases(t, tests)
}

func TestAndExpressionEval(t *testing.T) {
	tests := []evalTestCase{
		{
			"test and expression with two variables",
			NewAndExpression(NewVarExpression("a"), NewVarExpression("b")),
			map[string]bool{"a": true, "b": false},
			false,
		},
		{
			"test and expression with one variable and one number",
			NewOrExpression(NewVarExpression("a"), NewNumberExpression(0)),
			map[string]bool{"a": false},
			false,
		},
	}

	runEvalTestCases(t, tests)
}

func TestImpliesExpressionEval(t *testing.T) {
	tests := []evalTestCase{
		{
			"test implies expression like 0 -> 0",
			NewImpliesExpression(NewVarExpression("a"), NewVarExpression("b")),
			map[string]bool{"a": false, "b": false},
			true,
		},
		{
			"test implies expression like 0 -> 1",
			NewImpliesExpression(NewVarExpression("a"), NewVarExpression("b")),
			map[string]bool{"a": false, "b": true},
			true,
		},
		{
			"test implies expression like 1 -> 0",
			NewImpliesExpression(NewVarExpression("a"), NewVarExpression("b")),
			map[string]bool{"a": true, "b": false},
			false,
		},
		{
			"test implies expression like 1 -> 1",
			NewImpliesExpression(NewVarExpression("a"), NewVarExpression("b")),
			map[string]bool{"a": true, "b": true},
			true,
		},
	}

	runEvalTestCases(t, tests)
}

func TestXORExpressionEval(t *testing.T) {
	tests := []evalTestCase{
		{
			"test xor expression like 0 + 0",
			NewXORExpression(NewVarExpression("a"), NewVarExpression("b")),
			map[string]bool{"a": false, "b": false},
			false,
		},
		{
			"test xor expression like 0 + 1",
			NewXORExpression(NewVarExpression("a"), NewVarExpression("b")),
			map[string]bool{"a": false, "b": true},
			true,
		},
		{
			"test xor expression like 1 + 0",
			NewXORExpression(NewVarExpression("a"), NewVarExpression("b")),
			map[string]bool{"a": true, "b": false},
			true,
		},
		{
			"test xor expression like 1 + 1",
			NewXORExpression(NewVarExpression("a"), NewVarExpression("b")),
			map[string]bool{"a": true, "b": true},
			false,
		},
	}

	runEvalTestCases(t, tests)
}

func TestNumberExpressionEval(t *testing.T) {
	tests := []evalTestCase{
		{
			"test number expression with value = 0",
			NewNumberExpression(0),
			nil,
			false,
		},
		{
			"test number expression with value = 1",
			NewNumberExpression(1),
			nil,
			true,
		},
	}

	runEvalTestCases(t, tests)
}
