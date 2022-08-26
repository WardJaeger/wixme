// Ward Jaeger, CS 403
package main

// Any possible expression in WIXME
type Expr interface {
	accept(ExprVisitor) any
}

// A Visitor pattern interface for expressions
type ExprVisitor interface {
	visitAssignExpr(*AssignExpr) any
	visitBinaryExpr(*BinaryExpr) any
	visitCallExpr(*CallExpr) any
	visitGetExpr(*GetExpr) any
	visitIndexExpr(*IndexExpr) any
	visitGroupingExpr(*GroupingExpr) any
	visitListExpr(*ListExpr) any
	visitLiteralExpr(*LiteralExpr) any
	visitLogicalExpr(*LogicalExpr) any
	visitReplaceExpr(*ReplaceExpr) any
	visitSetExpr(*SetExpr) any
	visitTernaryExpr(*TernaryExpr) any
	visitThisExpr(*ThisExpr) any
	visitUnaryExpr(*UnaryExpr) any
	visitVariableExpr(*VariableExpr) any
}

// Assign an existing variable to a new value
type AssignExpr struct {
	name  Token
	value Expr
}

func (a *AssignExpr) accept(visitor ExprVisitor) any {
	return visitor.visitAssignExpr(a)
}

// Perform (arithmetic/comparison/concatenation) operation on two values
type BinaryExpr struct {
	left     Expr
	operator Token
	right    Expr
}

func (b *BinaryExpr) accept(visitor ExprVisitor) any {
	return visitor.visitBinaryExpr(b)
}

// Perform a call on a Callable
type CallExpr struct {
	callee    Expr
	arguments []Expr
	paren     Token
}

func (c *CallExpr) accept(visitor ExprVisitor) any {
	return visitor.visitCallExpr(c)
}

// Get value of instance property
type GetExpr struct {
	object Expr
	name   Token
}

func (g *GetExpr) accept(visitor ExprVisitor) any {
	return visitor.visitGetExpr(g)
}

// In parentheses, for correct evaluation order
type GroupingExpr struct {
	expression Expr
}

func (g *GroupingExpr) accept(visitor ExprVisitor) any {
	return visitor.visitGroupingExpr(g)
}

// Get index or slice copy of a Sequence
type IndexExpr struct {
	indexee Expr
	start   Expr
	stop    Expr
	bracket Token
}

func (c *IndexExpr) accept(visitor ExprVisitor) any {
	return visitor.visitIndexExpr(c)
}

// Create a new list
type ListExpr struct {
	elements []Expr
	bracket  Token
}

func (l *ListExpr) accept(visitor ExprVisitor) any {
	return visitor.visitListExpr(l)
}

// A literal value that needs no additional evaluation
type LiteralExpr struct {
	value any
}

func (l *LiteralExpr) accept(visitor ExprVisitor) any {
	return visitor.visitLiteralExpr(l)
}

// A logical control-flow operation on two values
type LogicalExpr struct {
	left     Expr
	operator Token
	right    Expr
}

func (l *LogicalExpr) accept(visitor ExprVisitor) any {
	return visitor.visitLogicalExpr(l)
}

// Replace an element of a Sequence at a given index
type ReplaceExpr struct {
	indexee Expr
	index   Expr
	bracket Token
	value   Expr
}

func (r *ReplaceExpr) accept(visitor ExprVisitor) any {
	return visitor.visitReplaceExpr(r)
}

// Set value of instance field
type SetExpr struct {
	object Expr
	name   Token
	value  Expr
}

func (s *SetExpr) accept(visitor ExprVisitor) any {
	return visitor.visitSetExpr(s)
}

// If condition is true, return trueValue, otherwise falseValue
type TernaryExpr struct {
	condition  Expr
	operator   Token
	trueValue  Expr
	falseValue Expr
}

func (t *TernaryExpr) accept(visitor ExprVisitor) any {
	return visitor.visitTernaryExpr(t)
}

// Special object referring to the current instance
type ThisExpr struct {
	Token
}

func (t *ThisExpr) accept(visitor ExprVisitor) any {
	return visitor.visitThisExpr(t)
}

// Perform an operation on one value
type UnaryExpr struct {
	operator Token
	operand  Expr
}

func (u *UnaryExpr) accept(visitor ExprVisitor) any {
	return visitor.visitUnaryExpr(u)
}

// Variable name
type VariableExpr struct {
	Token
}

func (v *VariableExpr) accept(visitor ExprVisitor) any {
	return visitor.visitVariableExpr(v)
}
