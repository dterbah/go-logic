package logic

import (
	"fmt"
	"strings"
)

type Expression interface {
	Eval(variables map[string]bool) bool
	String() string
	ToDot(builder *strings.Builder, parentID string)
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

func (notExprt NotExpression) String() string {
	return fmt.Sprintf("! %s", notExprt.expr)
}

func (notExpr *NotExpression) ToDot(builder *strings.Builder, parentID string) {
	nodeID := fmt.Sprintf("not_%p", notExpr)
	builder.WriteString(fmt.Sprintf("\"%s\" [label=\"NOT\"];\n", nodeID))
	if parentID != "" {
		builder.WriteString(fmt.Sprintf("\"%s\" -> \"%s\";\n", parentID, nodeID))
	}
	notExpr.expr.ToDot(builder, nodeID)
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

func (varExpr VarExpression) String() string {
	return varExpr.variable
}

func (varExpr *VarExpression) ToDot(builder *strings.Builder, parentID string) {
	nodeID := fmt.Sprintf("var_%s", varExpr.variable)
	builder.WriteString(fmt.Sprintf("\"%s\" [label=\"%s\"];\n", nodeID, varExpr.variable))
	if parentID != "" {
		builder.WriteString(fmt.Sprintf("\"%s\" -> \"%s\";\n", parentID, nodeID))
	}
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

func (orExpr OrExpression) String() string {
	return fmt.Sprintf("%s v %s", orExpr.left, orExpr.right)
}

func (orExpr *OrExpression) ToDot(builder *strings.Builder, parentID string) {
	nodeID := fmt.Sprintf("or_%p", orExpr)
	builder.WriteString(fmt.Sprintf("\"%s\" [label=\"OR\"];\n", nodeID))
	if parentID != "" {
		builder.WriteString(fmt.Sprintf("\"%s\" -> \"%s\";\n", parentID, nodeID))
	}
	orExpr.left.ToDot(builder, nodeID)
	orExpr.right.ToDot(builder, nodeID)
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

func (andExpr AndExpression) String() string {
	return fmt.Sprintf("%s ^ %s", andExpr.left, andExpr.right)
}

func (andExpr *AndExpression) ToDot(builder *strings.Builder, parentID string) {
	nodeID := fmt.Sprintf("and_%p", andExpr)
	builder.WriteString(fmt.Sprintf("\"%s\" [label=\"AND\"];\n", nodeID))
	if parentID != "" {
		builder.WriteString(fmt.Sprintf("\"%s\" -> \"%s\";\n", parentID, nodeID))
	}
	andExpr.left.ToDot(builder, nodeID)
	andExpr.right.ToDot(builder, nodeID)
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

func (impliesExpr ImpliesExpression) String() string {
	return fmt.Sprintf("%s -> %s", impliesExpr.left, impliesExpr.right)
}

func (impliesExpr *ImpliesExpression) ToDot(builder *strings.Builder, parentID string) {
	nodeID := fmt.Sprintf("implies_%p", impliesExpr)
	builder.WriteString(fmt.Sprintf("\"%s\" [label=\"IMPLIES\"];\n", nodeID))
	if parentID != "" {
		builder.WriteString(fmt.Sprintf("\"%s\" -> \"%s\";\n", parentID, nodeID))
	}
	impliesExpr.left.ToDot(builder, nodeID)
	impliesExpr.right.ToDot(builder, nodeID)
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

func (xorExpr XORExpression) String() string {
	return fmt.Sprintf("%s ⊕ %s", xorExpr.left, xorExpr.right)
}

func (xorExpr *XORExpression) ToDot(builder *strings.Builder, parentID string) {
	nodeID := fmt.Sprintf("xor_%p", xorExpr)
	builder.WriteString(fmt.Sprintf("\"%s\" [label=\"XOR\"];\n", nodeID))
	if parentID != "" {
		builder.WriteString(fmt.Sprintf("\"%s\" -> \"%s\";\n", parentID, nodeID))
	}
	xorExpr.left.ToDot(builder, nodeID)
	xorExpr.right.ToDot(builder, nodeID)
}
