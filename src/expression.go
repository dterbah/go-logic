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

type AndExpression struct {
	left, right Expression
}

func NewAndExpression(left, right Expression) *AndExpression {
	return &AndExpression{left: left, right: right}
}

func (andExpr *AndExpression) Eval(variables map[string]bool) bool {
	return andExpr.left.Eval(variables) && andExpr.right.Eval(variables)
}
