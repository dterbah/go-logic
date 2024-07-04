package logic

import (
	"fmt"
	"strings"
)

const DOT_FORMAT = "\"%s\" -> \"%s\";\n"

type Expression interface {
	Eval(variables map[string]bool) bool
	String() string
	ToDot(builder *strings.Builder, parentID string)
	Simplify() Expression
	equal(expr Expression) bool
}

// Not Expression API
type NotExpression struct {
	expr Expression
}

func NewNotExpression(expr Expression) *NotExpression {
	return &NotExpression{expr: expr}
}

func (notExpr NotExpression) equal(expr Expression) bool {
	if value, ok := expr.(*NotExpression); ok {
		return value.expr.equal(notExpr.expr)
	}

	return false
}

func (notExpr *NotExpression) Eval(variables map[string]bool) bool {
	return !notExpr.expr.Eval(variables)
}

func (notExpr *NotExpression) Simplify() Expression {
	return notExpr
}

func (notExprt NotExpression) String() string {
	return fmt.Sprintf("! %s", notExprt.expr)
}

func (notExpr *NotExpression) ToDot(builder *strings.Builder, parentID string) {
	nodeID := fmt.Sprintf("not_%p", notExpr)
	builder.WriteString(fmt.Sprintf("\"%s\" [label=\"NOT\"];\n", nodeID))
	if parentID != "" {
		builder.WriteString(fmt.Sprintf(DOT_FORMAT, parentID, nodeID))
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

func (varExpr VarExpression) equal(expr Expression) bool {
	if value, ok := expr.(*VarExpression); ok {
		return value.variable == varExpr.variable
	}

	return false
}

func (varExpr *VarExpression) Eval(variables map[string]bool) bool {
	return variables[varExpr.variable]
}

func (varExpr *VarExpression) Simplify() Expression {
	return varExpr
}

func (varExpr VarExpression) String() string {
	return varExpr.variable
}

func (varExpr *VarExpression) ToDot(builder *strings.Builder, parentID string) {
	nodeID := fmt.Sprintf("var_%s", varExpr.variable)
	builder.WriteString(fmt.Sprintf("\"%s\" [label=\"%s\"];\n", nodeID, varExpr.variable))
	if parentID != "" {
		builder.WriteString(fmt.Sprintf(DOT_FORMAT, parentID, nodeID))
	}
}

// Or Expression API
type OrExpression struct {
	left, right Expression
}

func NewOrExpression(left, right Expression) *OrExpression {
	return &OrExpression{left: left, right: right}
}

func (orExpr OrExpression) equal(expr Expression) bool {
	if value, ok := expr.(*OrExpression); ok {
		return value.left.equal(orExpr.left) && value.right.equal(orExpr.right)
	}

	return false
}

func (orExpr *OrExpression) Eval(variables map[string]bool) bool {
	return orExpr.left.Eval(variables) || orExpr.right.Eval(variables)
}

func (orExpr *OrExpression) Simplify() Expression {
	// Idempotence
	if orExpr.right.equal(orExpr.left) {
		return orExpr.right
	}

	// Identity
	if value, ok := orExpr.right.(*NumberExpression); ok && value.value == 0 {
		return orExpr.left
	}

	if value, ok := orExpr.left.(*NumberExpression); ok && value.value == 0 {
		return orExpr.right
	}

	// Domination
	if value, ok := orExpr.right.(*NumberExpression); ok && value.value == 1 {
		return NewNumberExpression(1)
	}

	if value, ok := orExpr.left.(*NumberExpression); ok && value.value == 1 {
		return NewNumberExpression(1)
	}

	return NewOrExpression(orExpr.left.Simplify(), orExpr.right.Simplify())
}

func (orExpr OrExpression) String() string {
	return fmt.Sprintf("%s v %s", orExpr.left, orExpr.right)
}

func (orExpr *OrExpression) ToDot(builder *strings.Builder, parentID string) {
	nodeID := fmt.Sprintf("or_%p", orExpr)
	builder.WriteString(fmt.Sprintf("\"%s\" [label=\"OR\"];\n", nodeID))
	if parentID != "" {
		builder.WriteString(fmt.Sprintf(DOT_FORMAT, parentID, nodeID))
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

func (andExpr AndExpression) equal(expr Expression) bool {
	if value, ok := expr.(*AndExpression); ok {
		return value.left.equal(andExpr.left) && value.right.equal(andExpr.right)
	}

	return false
}

func (andExpr *AndExpression) Eval(variables map[string]bool) bool {
	return andExpr.left.Eval(variables) && andExpr.right.Eval(variables)
}

func (andExpr *AndExpression) Simplify() Expression {
	// Idempotence
	if andExpr.right.equal(andExpr.left) {
		return andExpr.right
	}

	// Identity
	if value, ok := andExpr.right.(*NumberExpression); ok && value.value == 1 {
		return andExpr.left
	}

	if value, ok := andExpr.left.(*NumberExpression); ok && value.value == 1 {
		return andExpr.right
	}

	// Domination
	if value, ok := andExpr.right.(*NumberExpression); ok && value.value == 0 {
		return NewNumberExpression(0)
	}

	if value, ok := andExpr.left.(*NumberExpression); ok && value.value == 0 {
		return NewNumberExpression(0)
	}

	return NewAndExpression(andExpr.left.Simplify(), andExpr.right.Simplify()).Simplify()
}

func (andExpr AndExpression) String() string {
	return fmt.Sprintf("%s ^ %s", andExpr.left, andExpr.right)
}

func (andExpr *AndExpression) ToDot(builder *strings.Builder, parentID string) {
	nodeID := fmt.Sprintf("and_%p", andExpr)
	builder.WriteString(fmt.Sprintf("\"%s\" [label=\"AND\"];\n", nodeID))
	if parentID != "" {
		builder.WriteString(fmt.Sprintf(DOT_FORMAT, parentID, nodeID))
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

func (impliesExpression ImpliesExpression) equal(expr Expression) bool {
	if value, ok := expr.(*ImpliesExpression); ok {
		return value.left.equal(impliesExpression.left) && value.right.equal(impliesExpression.right)
	}

	return false
}

func (impliesExpr *ImpliesExpression) Eval(variables map[string]bool) bool {
	left := impliesExpr.left.Eval(variables)
	right := impliesExpr.right.Eval(variables)

	return !left || right
}

func (impliesExpr *ImpliesExpression) Simplify() Expression {
	return impliesExpr
}

func (impliesExpr ImpliesExpression) String() string {
	return fmt.Sprintf("%s -> %s", impliesExpr.left, impliesExpr.right)
}

func (impliesExpr *ImpliesExpression) ToDot(builder *strings.Builder, parentID string) {
	nodeID := fmt.Sprintf("implies_%p", impliesExpr)
	builder.WriteString(fmt.Sprintf("\"%s\" [label=\"IMPLIES\"];\n", nodeID))
	if parentID != "" {
		builder.WriteString(fmt.Sprintf(DOT_FORMAT, parentID, nodeID))
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

func (xorExpr XORExpression) equal(expr Expression) bool {
	if value, ok := expr.(*XORExpression); ok {
		return value.left.equal(xorExpr.left) && value.right.equal(xorExpr.right)
	}

	return false
}

func (xorExpr *XORExpression) Eval(variables map[string]bool) bool {
	return xorExpr.left.Eval(variables) != xorExpr.right.Eval(variables)
}

func (xorExpr *XORExpression) Simplify() Expression {
	return xorExpr
}

func (xorExpr XORExpression) String() string {
	return fmt.Sprintf("%s ⊕ %s", xorExpr.left, xorExpr.right)
}

func (xorExpr *XORExpression) ToDot(builder *strings.Builder, parentID string) {
	nodeID := fmt.Sprintf("xor_%p", xorExpr)
	builder.WriteString(fmt.Sprintf("\"%s\" [label=\"XOR\"];\n", nodeID))
	if parentID != "" {
		builder.WriteString(fmt.Sprintf(DOT_FORMAT, parentID, nodeID))
	}
	xorExpr.left.ToDot(builder, nodeID)
	xorExpr.right.ToDot(builder, nodeID)
}

// Number expression API
type NumberExpression struct {
	Expression
	value int
}

func NewNumberExpression(value int) *NumberExpression {
	return &NumberExpression{value: value}
}

func (nbrExpr NumberExpression) equal(expr Expression) bool {
	if value, ok := expr.(*NumberExpression); ok {
		return value.value == nbrExpr.value
	}

	return false
}

func (nbrExpr *NumberExpression) Eval(variables map[string]bool) bool {
	return nbrExpr.value == 1
}

func (nbrExpr *NumberExpression) Simplify() Expression {
	return nbrExpr
}

func (nbrExpr NumberExpression) String() string {
	return fmt.Sprintf("%d", nbrExpr.value)
}

func (nbrExpr *NumberExpression) ToDot(builder *strings.Builder, parentID string) {
	nodeID := fmt.Sprintf("number_%d", nbrExpr.value)
	builder.WriteString(fmt.Sprintf("\"%s\" [label=\"%d\"];\n", nodeID, nbrExpr.value))
	if parentID != "" {
		builder.WriteString(fmt.Sprintf(DOT_FORMAT, parentID, nodeID))
	}
}
