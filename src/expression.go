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
	if value, ok := notExpr.expr.(*OrExpression); ok {
		// De Morgan's law: !(a || b) => !a && !b
		return NewAndExpression(NewNotExpression(value.left), NewNotExpression(value.right)).Simplify()
	}

	if value, ok := notExpr.expr.(*AndExpression); ok {
		// De Morgan's law: !(a && b) => !a || !b
		return NewOrExpression(NewNotExpression(value.left), NewNotExpression(value.right)).Simplify()
	}

	if value, ok := notExpr.expr.(*NumberExpression); ok {
		if value.value == 0 {
			return NewNumberExpression(1)
		} else {
			return NewNumberExpression(0)
		}
	}

	if value, ok := notExpr.expr.(*NotExpression); ok {
		return value.expr.Simplify()
	}

	return notExpr
}

func (notExprt NotExpression) String() string {
	return fmt.Sprintf("!%s", notExprt.expr)
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
	// Idempotence: a || a = a
	if orExpr.left.equal(orExpr.right) {
		return orExpr.left.Simplify()
	}

	// Identity: a || false = a
	if value, ok := orExpr.right.(*NumberExpression); ok && value.value == 0 {
		return orExpr.left.Simplify()
	}

	if value, ok := orExpr.left.(*NumberExpression); ok && value.value == 0 {
		return orExpr.right.Simplify()
	}

	// Domination: a || true = true
	if value, ok := orExpr.right.(*NumberExpression); ok && value.value == 1 {
		return NewNumberExpression(1)
	}

	if value, ok := orExpr.left.(*NumberExpression); ok && value.value == 1 {
		return NewNumberExpression(1)
	}

	// Complementarity: a || !a = true
	if value, ok := orExpr.right.(*NotExpression); ok && value.expr.equal(orExpr.left) {
		return NewNumberExpression(1)
	}

	if value, ok := orExpr.left.(*NotExpression); ok && value.expr.equal(orExpr.right) {
		return NewNumberExpression(1)
	}

	// Absorption: a || (a && b) = a
	if value, ok := orExpr.right.(*AndExpression); ok && orExpr.left.equal(value.left) {
		return orExpr.left.Simplify()
	}

	return &OrExpression{
		left:  orExpr.left.Simplify(),
		right: orExpr.right.Simplify(),
	}
}

func (orExpr OrExpression) String() string {
	return fmt.Sprintf("%sv%s", orExpr.left, orExpr.right)
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
	// Idempotence: a && a = a
	if andExpr.left.equal(andExpr.right) {
		return andExpr.left.Simplify()
	}

	// Identity: a && true = a
	if value, ok := andExpr.right.(*NumberExpression); ok && value.value == 1 {
		return andExpr.left.Simplify()
	}

	if value, ok := andExpr.left.(*NumberExpression); ok && value.value == 1 {
		return andExpr.right.Simplify()
	}

	// Domination: a && false = false
	if andExpr.isDomination() {
		return NewNumberExpression(0)
	}

	// Complementarity: a && !a = false
	if andExpr.isComplementarity() {
		return NewNumberExpression(0)
	}

	// Absorption: a && (a || b) = a
	if value, ok := andExpr.right.(*OrExpression); ok && andExpr.left.equal(value.left) {
		return andExpr.left.Simplify()
	}

	left := andExpr.left.Simplify()
	right := andExpr.right.Simplify()

	if left.equal(right) {
		return left
	}

	return &AndExpression{
		left:  andExpr.left.Simplify(),
		right: andExpr.right.Simplify(),
	}
}

func (andExpr AndExpression) isComplementarity() bool {
	if value, ok := andExpr.right.(*NotExpression); ok && value.expr.equal(andExpr.left) {
		return true
	} else if value, ok := andExpr.left.(*NotExpression); ok && value.expr.equal(andExpr.right) {
		return true
	}

	return false
}

func (andExpr AndExpression) isDomination() bool {
	if value, ok := andExpr.right.(*NumberExpression); ok && value.value == 0 {
		return true
	}

	if value, ok := andExpr.left.(*NumberExpression); ok && value.value == 0 {
		return true
	}

	return false
}

func (andExpr AndExpression) String() string {
	return fmt.Sprintf("%s^%s", andExpr.left, andExpr.right)
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
	return NewOrExpression(NewNotExpression(impliesExpr.left), impliesExpr.right).Simplify()
}

func (impliesExpr ImpliesExpression) String() string {
	return fmt.Sprintf("%s->%s", impliesExpr.left, impliesExpr.right)
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
	// a + 0 --> a, 0 + a --> a
	if value, ok := xorExpr.left.(*NumberExpression); ok {
		if value.value == 0 {
			return xorExpr.right
		}

		return NewNotExpression(xorExpr.right).Simplify()
	}

	if value, ok := xorExpr.right.(*NumberExpression); ok {
		if value.value == 0 {
			return xorExpr.left
		}

		return NewNotExpression(xorExpr.left).Simplify()
	}

	// a + a --> 0
	if xorExpr.right.equal(xorExpr.left) {
		return NewNumberExpression(0)
	}

	// a + (a+b) --> b
	if value, ok := xorExpr.right.(*XORExpression); ok {
		if xorExpr.left.equal(value.left) {
			return value.right.Simplify()
		}

		// (a+b)+(a+c)
		if left, leftOK := xorExpr.left.(*XORExpression); leftOK {
			if left.left.equal(value.left) {
				return NewXORExpression(left.right, value.right).Simplify()
			}
		}
	}

	return xorExpr
}

func (xorExpr XORExpression) String() string {
	return fmt.Sprintf("%s⊕%s", xorExpr.left, xorExpr.right)
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

// Equivalence expression API
type EquivalenceExpression struct {
	Expression
	left, right Expression
}

func NewEquivalenceExpression(left, right Expression) *EquivalenceExpression {
	return &EquivalenceExpression{left: left, right: right}
}

func (equivalenceExpr EquivalenceExpression) equal(expr Expression) bool {
	if value, ok := expr.(*EquivalenceExpression); ok {
		return value.left.equal(equivalenceExpr.left) && value.right.equal(equivalenceExpr.right)
	}

	return false
}

func (equivalenceExpr EquivalenceExpression) Eval(variables map[string]bool) bool {
	return equivalenceExpr.left.Eval(variables) == equivalenceExpr.right.Eval(variables)
}

func (equivalenceExpr *EquivalenceExpression) Simplify() Expression {
	if left, ok := equivalenceExpr.left.(*NumberExpression); ok {
		if left.value == 1 {
			// 1 <-> B => B
			return equivalenceExpr.right.Simplify()
		} else if left.value == 0 {
			// 0 <-> B => !B
			return NewNotExpression(equivalenceExpr.right).Simplify()
		}
	}

	if right, ok := equivalenceExpr.right.(*NumberExpression); ok {
		if right.value == 1 {
			// A <-> 1 => A
			return equivalenceExpr.left.Simplify()
		} else if right.value == 0 {
			// A <-> 0 => !A
			return NewNotExpression(equivalenceExpr.left).Simplify()
		}
	}

	expr := NewOrExpression(
		NewAndExpression(equivalenceExpr.left, equivalenceExpr.right),
		NewAndExpression(NewNotExpression(equivalenceExpr.left), NewNotExpression(equivalenceExpr.right)),
	).Simplify()

	return expr
}

func (equivalenceExpression EquivalenceExpression) String() string {
	return fmt.Sprintf("%s<->%s", equivalenceExpression.left, equivalenceExpression.right)
}

func (equivalenceExpr *EquivalenceExpression) ToDot(builder *strings.Builder, parentID string) {
	nodeID := fmt.Sprintf("equ_%p", equivalenceExpr)
	builder.WriteString(fmt.Sprintf("\"%s\" [label=\"EQU\"];\n", nodeID))
	if parentID != "" {
		builder.WriteString(fmt.Sprintf(DOT_FORMAT, parentID, nodeID))
	}
	equivalenceExpr.left.ToDot(builder, nodeID)
	equivalenceExpr.right.ToDot(builder, nodeID)
}
