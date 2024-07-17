package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type equalTestCase struct {
	name           string
	firstExpr      Expression
	secondExpr     Expression
	expectedResult bool
}

func runEqualTestCases(t *testing.T, tests []equalTestCase) {
	assert := assert.New(t)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(test.firstExpr.equal(test.secondExpr), test.expectedResult, test.name)
		})
	}
}

// Not Expression //
func TestNotExpressionEqual(t *testing.T) {
	tests := []equalTestCase{
		{"test not expression equal to another not expression", NewNotExpression(NewVarExpression("a")), NewNotExpression(NewVarExpression("a")), true},
		{"test not expression not equal to another expression", NewNotExpression(NewVarExpression("a")), NewVarExpression("b"), false},
	}

	runEqualTestCases(t, tests)
}

func TestVarExpressionEqual(t *testing.T) {
	tests := []equalTestCase{
		{"test var expression equal to another var expression", NewVarExpression("a"), NewVarExpression("a"), true},
		{"test var expression not equal to another expression", NewVarExpression("a"), NewNotExpression(NewVarExpression("a")), false},
	}

	runEqualTestCases(t, tests)
}

func TestOrExpressionEqual(t *testing.T) {
	tests := []equalTestCase{
		{
			"test or expression equal to another or expression",
			NewOrExpression(NewVarExpression("a"), NewVarExpression("b")),
			NewOrExpression(NewVarExpression("a"), NewVarExpression("b")),
			true,
		},
		{
			"test or expression not equal to another expression",
			NewOrExpression(NewVarExpression("a"), NewVarExpression("b")),
			NewVarExpression("a"),
			false,
		},
	}

	runEqualTestCases(t, tests)
}

func TestAndExpressionEqual(t *testing.T) {
	tests := []equalTestCase{
		{
			"test and expression equal to another and expression",
			NewAndExpression(NewVarExpression("a"), NewVarExpression("b")),
			NewAndExpression(NewVarExpression("a"), NewVarExpression("b")),
			true,
		},
		{
			"test and expression not equal to another expression",
			NewAndExpression(NewVarExpression("a"), NewVarExpression("b")),
			NewVarExpression("b"),
			false,
		},
	}

	runEqualTestCases(t, tests)
}

func TestImpliesExpressionEqual(t *testing.T) {
	tests := []equalTestCase{
		{
			"test implies expression equal to another implies expression",
			NewImpliesExpression(NewVarExpression("a"), NewVarExpression("b")),
			NewImpliesExpression(NewVarExpression("a"), NewVarExpression("b")),
			true,
		},
		{
			"test implies expression not equal to another expression",
			NewImpliesExpression(NewVarExpression("a"), NewVarExpression("b")),
			NewAndExpression(NewVarExpression("a"), NewVarExpression("c")),
			false,
		},
	}

	runEqualTestCases(t, tests)
}

func TestXORExpressionEqual(t *testing.T) {
	tests := []equalTestCase{
		{
			"test xor expression equal to another xor expression",
			NewXORExpression(NewVarExpression("a"), NewVarExpression("b")),
			NewXORExpression(NewVarExpression("a"), NewVarExpression("b")),
			true,
		},
		{
			"test xor expression not equal to another expression",
			NewXORExpression(NewVarExpression("a"), NewVarExpression("b")),
			NewAndExpression(NewVarExpression("a"), NewVarExpression("c")),
			false,
		},
	}

	runEqualTestCases(t, tests)
}

func TestNumberExpressionEqual(t *testing.T) {
	tests := []equalTestCase{
		{
			"test number expression equal to another number expression",
			NewNumberExpression(1), NewNumberExpression(1),
			true,
		},
		{
			"test number expression not equal to another expression",
			NewNumberExpression(1), NewNotExpression(NewNumberExpression(1)),
			false,
		},
	}

	runEqualTestCases(t, tests)
}

func TestEquivalenceExpressionEqual(t *testing.T) {
	tests := []equalTestCase{
		{
			"test equivalence equal to another equivalence expression",
			NewEquivalenceExpression(NewVarExpression("a"), NewVarExpression("b")),
			NewEquivalenceExpression(NewVarExpression("a"), NewVarExpression("b")),
			true,
		},
		{
			"test equivalence not equal to another expression",
			NewEquivalenceExpression(NewVarExpression("a"), NewVarExpression("b")),
			NewNumberExpression(1),
			false,
		},
	}

	runEqualTestCases(t, tests)
}
