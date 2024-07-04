package logic

type Expression interface {
	Eval(variables map[string]bool) bool
}

// Not Expression API
type NotExpression struct {
	expr Expression
}

func NewNotExpression(expr Expression) *NotExpression {
	return &NotExpression{expr: expr}
}

func (notExpr *NotExpression) Eval(variables map[string]bool) bool {
	return !notExpr.expr.Eval(variables)
}

// Var Expression API
type VarExpression struct {
	variable string
}

func NewVarExpression(variable string) *VarExpression {
	return &VarExpression{variable: variable}
}

func (varExpr *VarExpression) Eval(variables map[string]bool) bool {
	return variables[varExpr.variable]
}

// Or Expression API
type OrExpression struct {
	left, right Expression
}

func NewOrExpression(left, right Expression) *OrExpression {
	return &OrExpression{left: left, right: right}
}

func (orExpr *OrExpression) Eval(variables map[string]bool) bool {
	return orExpr.left.Eval(variables) || orExpr.right.Eval(variables)
}

// And Expression API
type AndExpression struct {
	left, right Expression
}

func NewAndExpression(left, right Expression) *AndExpression {
	return &AndExpression{left: left, right: right}
}

func (andExpr *AndExpression) Eval(variables map[string]bool) bool {
	return andExpr.left.Eval(variables) && andExpr.right.Eval(variables)
}

// Implies Expression API
type ImpliesExpression struct {
	left, right Expression
}

func NewImpliesExpression(left, right Expression) *ImpliesExpression {
	return &ImpliesExpression{left: left, right: right}
}

func (impliesExpr *ImpliesExpression) Eval(variables map[string]bool) bool {
	left := impliesExpr.left.Eval(variables)
	right := impliesExpr.right.Eval(variables)

	return !left || right
}

// XOR Expression API
type XORExpression struct {
	left, right Expression
}

func NewXORExpression(left, right Expression) *XORExpression {
	return &XORExpression{left: left, right: right}
}

func (xorExpr *XORExpression) Eval(variables map[string]bool) bool {
	return xorExpr.left.Eval(variables) != xorExpr.right.Eval(variables)
}
