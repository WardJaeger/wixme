// Ward Jaeger, CS 403
package main

// Define "enum" type
type functionType int

// "Enum" for where we are in regards to function declarations
// Go does not actually have enums, but constant values work as a substitute
const (
	NONE functionType = iota
	FUNCTION
	INITIALIZER
	METHOD
)

// Visitor pattern that resolves references for identifiers
// Updates the interpreter as it goes
type Resolver struct {
	interpreter     *Interpreter
	scopes          []map[string]bool
	currentFunction functionType
	inClass         bool
}

// Test for interface implementation
var _ ExprVisitor = &Resolver{}
var _ StmtVisitor = &Resolver{}

// Entry point for resolution
func (r *Resolver) resolve(statements []Stmt) {
	for _, statement := range statements {
		r.resolveStmt(statement)
	}
}

// Pass resolver to statements and expressions
func (r *Resolver) resolveStmt(stmt Stmt) {
	stmt.accept(r)
}
func (r *Resolver) resolveExpr(expr Expr) {
	expr.accept(r)
}

// Creates an additional scope one level deeper
func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, map[string]bool{})
}

// Removes the most recent scope
func (r *Resolver) endScope() {
	r.scopes = r.scopes[0 : len(r.scopes)-1]
}

// Mark a variable as newly declared in the current scope
func (r *Resolver) declare(name Token) {
	if length := len(r.scopes); length != 0 {
		if _, found := r.scopes[length-1][name.lexeme]; found {
			reportToken(name, "Already a variable with this name in this scope.")
		}

		r.scopes[length-1][name.lexeme] = false
	}
}

// Mark a variable as already defined in the current scope
func (r *Resolver) define(name Token) {
	if length := len(r.scopes); length != 0 {
		r.scopes[length-1][name.lexeme] = true
	}
}

// Resolves a local variable to the correct depth
func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, found := r.scopes[i][name.lexeme]; found {
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

// Defines the parameters in a new scope, and resolves the body
func (r *Resolver) resolveFunction(function *FunctionStmt, ftype functionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = ftype
	r.beginScope()

	for _, param := range function.params {
		r.declare(param)
		r.define(param)
	}
	r.resolve(function.body)

	r.endScope()
	r.currentFunction = enclosingFunction
}

// Resolve the statements
func (r *Resolver) visitBlockStmt(stmt *BlockStmt) any {
	r.beginScope()
	r.resolve(stmt.statements)
	r.endScope()
	return nil
}

// Defines the class, defines "this" in a new scope, and resolves the methods
func (r *Resolver) visitClassStmt(stmt *ClassStmt) any {
	enclosingClass := r.inClass
	r.inClass = true

	r.declare(stmt.name)
	r.define(stmt.name)

	r.beginScope()
	r.scopes[len(r.scopes)-1]["this"] = true

	for _, method := range stmt.methods {
		declaration := METHOD
		if method.name.lexeme == "init" {
			declaration = INITIALIZER
		}

		r.resolveFunction(method, declaration)
	}

	r.endScope()

	r.inClass = enclosingClass
	return nil
}

// Resolves the expression
func (r *Resolver) visitExpressionStmt(stmt *ExpressionStmt) any {
	r.resolveExpr(stmt.expression)
	return nil
}

// Defines and resolves the function
func (r *Resolver) visitFunctionStmt(stmt *FunctionStmt) any {
	r.declare(stmt.name)
	r.define(stmt.name)

	r.resolveFunction(stmt, FUNCTION)
	return nil
}

// Resolves the condition and branches
func (r *Resolver) visitIfStmt(stmt *IfStmt) any {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.thenBranch)
	if stmt.elseBranch != nil {
		r.resolveStmt(stmt.elseBranch)
	}
	return nil
}

// Checks for location errors and resolves return value
func (r *Resolver) visitReturnStmt(stmt *ReturnStmt) any {
	if r.currentFunction == NONE {
		reportToken(stmt.keyword, "Can't return from top-level code.")
	}

	if stmt.value != nil {
		if r.currentFunction == INITIALIZER {
			reportToken(stmt.keyword,
				"Can't return a value from an initializer.")
		}

		r.resolveExpr(stmt.value)
	}
	return nil
}

// Declares, then resolves the value, then defines
func (r *Resolver) visitVarStmt(stmt *VarStmt) any {
	r.declare(stmt.name)
	if stmt.initializer != nil {
		r.resolveExpr(stmt.initializer)
	}
	r.define(stmt.name)
	return nil
}

// Resolves the condition and the body
func (r *Resolver) visitWhileStmt(stmt *WhileStmt) any {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.body)
	return nil
}

// Resolves the value, then resolves the variable
func (r *Resolver) visitAssignExpr(expr *AssignExpr) any {
	r.resolveExpr(expr.value)
	r.resolveLocal(expr, expr.name)
	return nil
}

// Resolves both operands
func (r *Resolver) visitBinaryExpr(expr *BinaryExpr) any {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)
	return nil
}

// Resolves the callee and the arguments
func (r *Resolver) visitCallExpr(expr *CallExpr) any {
	r.resolveExpr(expr.callee)

	for _, argument := range expr.arguments {
		r.resolveExpr(argument)
	}

	return nil
}

// Resolves the object
func (r *Resolver) visitGetExpr(expr *GetExpr) any {
	r.resolveExpr(expr.object)
	return nil
}

// Resolves the expression
func (r *Resolver) visitGroupingExpr(expr *GroupingExpr) any {
	r.resolveExpr(expr.expression)
	return nil
}

// Resolves the indexee and the indices
func (r *Resolver) visitIndexExpr(expr *IndexExpr) any {
	r.resolveExpr(expr.indexee)
	r.resolveExpr(expr.start)
	if expr.stop != nil {
		r.resolveExpr(expr.stop)
	}
	return nil
}

// Resolves each element
func (r *Resolver) visitListExpr(expr *ListExpr) any {
	for _, element := range expr.elements {
		r.resolveExpr(element)
	}
	return nil
}

// Does nothing, literals need no resolution
func (*Resolver) visitLiteralExpr(*LiteralExpr) any {
	return nil
}

// Resolves both operands
func (r *Resolver) visitLogicalExpr(expr *LogicalExpr) any {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)
	return nil
}

// Resolves indexee, index, and value
func (r *Resolver) visitReplaceExpr(expr *ReplaceExpr) any {
	r.resolveExpr(expr.indexee)
	r.resolveExpr(expr.index)
	r.resolveExpr(expr.value)
	return nil
}

// Resolves the object and the value
func (r *Resolver) visitSetExpr(expr *SetExpr) any {
	r.resolveExpr(expr.object)
	r.resolveExpr(expr.value)
	return nil
}

// Resolves the condition and both values
func (r *Resolver) visitTernaryExpr(expr *TernaryExpr) any {
	r.resolveExpr(expr.condition)
	r.resolveExpr(expr.trueValue)
	r.resolveExpr(expr.falseValue)
	return nil
}

// Checks for location error and resolves "this"
func (r *Resolver) visitThisExpr(expr *ThisExpr) any {
	if !r.inClass {
		reportToken(expr.Token,
			"Can't use 'this' outside of a class.")
		return nil
	}

	r.resolveLocal(expr, expr.Token)
	return nil
}

// Resolves operand
func (r *Resolver) visitUnaryExpr(expr *UnaryExpr) any {
	r.resolveExpr(expr.operand)
	return nil
}

// Checks if variable has not been defined yet, then resolves
func (r *Resolver) visitVariableExpr(expr *VariableExpr) any {
	if length := len(r.scopes); length != 0 {
		if defined, found := r.scopes[length-1][expr.lexeme]; found && !defined {
			reportToken(expr.Token, "Can't read local variable in its own initializer.")
		}
	}

	r.resolveLocal(expr, expr.Token)
	return nil
}
