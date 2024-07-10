package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type simplifyTestCase struct {
	name           string
	expr           Expression
	expectedResult Expression
}

func runSimplifyTestCases(t *testing.T, tests []simplifyTestCase) {
	assert := assert.New(t)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simplifiedExpr := test.expr.Simplify()
			assert.Equal(simplifiedExpr, test.expectedResult, test.name)
		})
	}
}

func TestNotExpressionSimplify(t *testing.T) {
	tests := []simplifyTestCase{
		{
			"test not expression simplify on !a",
			NewNotExpression(NewVarExpression("a")),
			NewNotExpression(NewVarExpression("a")),
		},
		{
			"test not expression simplify on !(a || b) --> De Morgan's law",
			NewNotExpression(NewOrExpression(NewVarExpression("a"), NewVarExpression("b"))),
			NewAndExpression(NewNotExpression(NewVarExpression("a")), NewNotExpression(NewVarExpression("b"))),
		},
		{
			"test not expression simplify on !(a && b) -> De Morgan's law",
			NewNotExpression(NewAndExpression(NewVarExpression("a"), NewVarExpression("b"))),
			NewOrExpression(NewNotExpression(NewVarExpression("a")), NewNotExpression(NewVarExpression("b"))),
		},
	}

	runSimplifyTestCases(t, tests)
}

func TestVarExpressionSimplify(t *testing.T) {
	tests := []simplifyTestCase{
		{
			"test var expression simplify",
			NewVarExpression("a"),
			NewVarExpression("a"),
		},
	}

	runSimplifyTestCases(t, tests)
}

func TestOrExpressionSimplify(t *testing.T) {
	tests := []simplifyTestCase{
		{
			"test or expression simplify on avb",
			NewOrExpression(NewVarExpression("a"), NewVarExpression("b")),
			NewOrExpression(NewVarExpression("a"), NewVarExpression("b")),
		},
		{
			"test or expression simplify on a || a --> Idempotence",
			NewOrExpression(NewVarExpression("a"), NewVarExpression("a")),
			NewVarExpression("a"),
		},
		{
			"test or expression simplify on a || false --> Identity",
			NewOrExpression(NewVarExpression("a"), NewNumberExpression(0)),
			NewVarExpression("a"),
		},
		{
			"test or expression simplify on false || a --> Identity",
			NewOrExpression(NewNumberExpression(0), NewVarExpression("a")),
			NewVarExpression("a"),
		},
		{
			"test or expression simplify on false || false",
			NewOrExpression(NewNumberExpression(0), NewNumberExpression(0)),
			NewNumberExpression(0),
		},
		{
			"test or expression simplify on true || true",
			NewOrExpression(NewNumberExpression(1), NewNumberExpression(1)),
			NewNumberExpression(1),
		},
		{
			"test or expression simplify on a || true --> Domination",
			NewOrExpression(NewVarExpression("a"), NewNumberExpression(1)),
			NewNumberExpression(1),
		},
		{
			"test or expression simplify on true || a --> Domination",
			NewOrExpression(NewNumberExpression(1), NewVarExpression("a")),
			NewNumberExpression(1),
		},
		{
			"test or expression simplify on a || !a --> Complementarity",
			NewOrExpression(NewVarExpression("a"), NewNotExpression(NewVarExpression("a"))),
			NewNumberExpression(1),
		},
		{
			"test or expression simplify on !a || a --> Complementarity",
			NewOrExpression(NewNotExpression(NewVarExpression("a")), NewVarExpression("a")),
			NewNumberExpression(1),
		},
		{
			"test or expression simplify on a || (a && b) --> Absorption",
			NewOrExpression(NewVarExpression("a"), NewAndExpression(NewVarExpression("a"), NewVarExpression("b"))),
			NewVarExpression("a"),
		},
	}

	runSimplifyTestCases(t, tests)
}

func TestAndExpressionSimplify(t *testing.T) {
	tests := []simplifyTestCase{
		{
			"test and expression simplify on a && b",
			NewAndExpression(NewVarExpression("a"), NewVarExpression("b")),
			NewAndExpression(NewVarExpression("a"), NewVarExpression("b")),
		},
		{
			"test and expression simplify on a && a --> Idempotence",
			NewAndExpression(NewVarExpression("a"), NewVarExpression("a")),
			NewVarExpression("a"),
		},
		{
			"test and expression simplify on a && true --> Identity",
			NewAndExpression(NewVarExpression("a"), NewNumberExpression(1)),
			NewVarExpression("a"),
		},
		{
			"test and expression simplify on true && a --> Identity",
			NewAndExpression(NewNumberExpression(1), NewVarExpression("a")),
			NewVarExpression("a"),
		},
		{
			"test and expression simplify on a && false --> Domination",
			NewAndExpression(NewVarExpression("a"), NewNumberExpression(0)),
			NewNumberExpression(0),
		},
		{
			"test and expression simplify on false && a --> Domination",
			NewAndExpression(NewNumberExpression(0), NewVarExpression("a")),
			NewNumberExpression(0),
		},
		{
			"test and expression simplify on a && !a --> Complementarity",
			NewAndExpression(NewVarExpression("a"), NewNotExpression(NewVarExpression("a"))),
			NewNumberExpression(0),
		},
		{
			"test and expression simplify on !a && a --> Complementarity",
			NewAndExpression(NewNotExpression(NewVarExpression("a")), NewVarExpression("a")),
			NewNumberExpression(0),
		},
		{
			"test and expression simplify on a && (a || b) --> Absorption",
			NewAndExpression(NewVarExpression("a"), NewOrExpression(NewVarExpression("a"), NewVarExpression("b"))),
			NewVarExpression("a"),
		},
		{
			"test and expression simplify on a&&a&&a",
			NewAndExpression(NewVarExpression("a"), NewAndExpression(NewVarExpression("a"), NewVarExpression("a"))),
			NewVarExpression("a"),
		},
	}

	runSimplifyTestCases(t, tests)
}

func TestImpliesExpressionSimplify(t *testing.T) {
	tests := []simplifyTestCase{
		{
			"test implies expression simplify on a->a",
			NewImpliesExpression(NewVarExpression("a"), NewVarExpression("a")),
			NewNumberExpression(1),
		},
		{
			"test implies expression simplify on a->0",
			NewImpliesExpression(NewVarExpression("a"), NewNumberExpression(0)),
			NewNotExpression(NewVarExpression("a")),
		},
	}

	runSimplifyTestCases(t, tests)
}

func TestXORExpressionSimplify(t *testing.T) {
	tests := []simplifyTestCase{
		{
			"test xor expression simplify on a + 0",
			NewXORExpression(NewVarExpression("a"), NewNumberExpression(0)),
			NewVarExpression("a"),
		},
		{
			"test xor expression simplify on 0 + a",
			NewXORExpression(NewNumberExpression(0), NewVarExpression("a")),
			NewVarExpression("a"),
		},
		{
			"test xor expression simplify on a + 1",
			NewXORExpression(NewVarExpression("a"), NewNumberExpression(1)),
			NewNotExpression(NewVarExpression("a")),
		},
		{
			"test xor expression simplify on 1 + a",
			NewXORExpression(NewNumberExpression(1), NewVarExpression("a")),
			NewNotExpression(NewVarExpression("a")),
		},
		{
			"test xor expression simplify on a + a",
			NewXORExpression(NewVarExpression("a"), NewVarExpression("a")),
			NewNumberExpression(0),
		},
		{
			"test xor expression simplify on a + (a+b)",
			NewXORExpression(NewVarExpression("a"), NewXORExpression(NewVarExpression("a"), NewVarExpression("b"))),
			NewVarExpression("b"),
		},
		{
			"test xor expression simplify on (a + b) + (a + c)",
			NewXORExpression(
				NewXORExpression(NewVarExpression("a"), NewVarExpression("b")),
				NewXORExpression(NewVarExpression("a"), NewVarExpression("c")),
			),
			NewXORExpression(NewVarExpression("b"), NewVarExpression("c")),
		},
	}

	runSimplifyTestCases(t, tests)
}
