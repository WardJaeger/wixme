// Ward Jaeger, CS 403
package main

// Any possible statement in WIXME
type Stmt interface {
	accept(StmtVisitor) any
}

// A Visitor pattern interface for statements
type StmtVisitor interface {
	visitBlockStmt(*BlockStmt) any
	visitClassStmt(*ClassStmt) any
	visitExpressionStmt(*ExpressionStmt) any
	visitFunctionStmt(*FunctionStmt) any
	visitIfStmt(*IfStmt) any
	visitReturnStmt(*ReturnStmt) any
	visitVarStmt(*VarStmt) any
	visitWhileStmt(*WhileStmt) any
}

// A list of statements
type BlockStmt struct {
	statements []Stmt
}

func (b *BlockStmt) accept(visitor StmtVisitor) any {
	return visitor.visitBlockStmt(b)
}

// Declare a new class
type ClassStmt struct {
	name    Token
	methods []*FunctionStmt
}

func (c *ClassStmt) accept(visitor StmtVisitor) any {
	return visitor.visitClassStmt(c)
}

// Execute expression like a statement
type ExpressionStmt struct {
	expression Expr
}

func (e *ExpressionStmt) accept(visitor StmtVisitor) any {
	return visitor.visitExpressionStmt(e)
}

// Define a new function/method
type FunctionStmt struct {
	name   Token
	params []Token
	body   []Stmt
}

func (f *FunctionStmt) accept(visitor StmtVisitor) any {
	return visitor.visitFunctionStmt(f)
}

// If condition is true, execute thenBranch, otherwise elseBranch if it exists
type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (i *IfStmt) accept(visitor StmtVisitor) any {
	return visitor.visitIfStmt(i)
}

// Return a value from the current function
type ReturnStmt struct {
	keyword Token
	value   Expr
}

func (r *ReturnStmt) accept(visitor StmtVisitor) any {
	return visitor.visitReturnStmt(r)
}

// Define a new variable
type VarStmt struct {
	name        Token
	initializer Expr
}

func (v *VarStmt) accept(visitor StmtVisitor) any {
	return visitor.visitVarStmt(v)
}

// While condition is true, execute the body
type WhileStmt struct {
	condition Expr
	body      Stmt
}

func (w *WhileStmt) accept(visitor StmtVisitor) any {
	return visitor.visitWhileStmt(w)
}
