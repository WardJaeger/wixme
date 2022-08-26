// Ward Jaeger, CS 403
package main

import (
	"fmt"
)

// Visitor pattern that evaluates an entire program of statements
type Interpreter struct {
	environment *Environment
	globals     *Environment
	locals      map[Expr]int
}

// Test for interface implementation
var _ ExprVisitor = &Interpreter{}
var _ StmtVisitor = &Interpreter{}

// Entry point for interpretation
func (i *Interpreter) interpret(statements []Stmt) {
	// Set up defered function to catch and report runtime errors
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(RuntimeError); ok {
				reportRuntime(err)
			} else {
				panic(r)
			}
		}
	}()

	for _, statement := range statements {
		i.execute(statement)
	}
}

// Pass interpreter to statements and expressions
func (i *Interpreter) execute(stmt Stmt) {
	stmt.accept(i)
}
func (i *Interpreter) evaluate(expr Expr) any {
	return expr.accept(i)
}

// Helper function for Interpreter that converts a value to a boolean
func isTruthy(value any) bool {
	return value != nil && value != false
}

// Execute a list of statements in a given environment
func (i *Interpreter) executeBlock(statements []Stmt, env *Environment) {
	previous := i.environment
	i.environment = env

	// Set up a defered function that returns the environment to its original state
	defer func() {
		i.environment = previous
	}()

	for _, statement := range statements {
		i.execute(statement)
	}
}

// Called by the Resolver to place variable references in the right depth
func (i *Interpreter) resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}

// Return variable at correct depth, or at global level
func (i *Interpreter) lookUpVariable(name Token, expr Expr) any {
	if distance, found := i.locals[expr]; found {
		return i.environment.getAt(distance, name.lexeme)
	} else {
		return i.globals.get(name)
	}
}

// Execute body of block statement in a new environment
func (i *Interpreter) visitBlockStmt(stmt *BlockStmt) any {
	newEnv := &Environment{enclosing: i.environment, values: map[string]any{}}
	i.executeBlock(stmt.statements, newEnv)
	return nil
}

// Define a new class, and define all its methods
func (i *Interpreter) visitClassStmt(stmt *ClassStmt) any {
	i.environment.define(stmt.name.lexeme, nil)
	methods := map[string]*Function{}
	for _, method := range stmt.methods {
		function := &Function{declaration: method, closure: i.environment,
			isInitializer: method.name.lexeme == "init"}
		methods[method.name.lexeme] = function
	}

	class := &Class{name: stmt.name.lexeme, methods: methods}
	i.environment.assign(stmt.name, class)
	return nil
}

// Evaluate expression and perform its side effects
func (i *Interpreter) visitExpressionStmt(stmt *ExpressionStmt) any {
	i.evaluate(stmt.expression)
	return nil
}

// Define a new function
func (i *Interpreter) visitFunctionStmt(stmt *FunctionStmt) any {
	function := &Function{declaration: stmt, closure: i.environment}
	i.environment.define(stmt.name.lexeme, function)
	return nil
}

// If condition is true, execute thenBranch, otherwise elseBranch if it exists
func (i *Interpreter) visitIfStmt(stmt *IfStmt) any {
	if isTruthy(i.evaluate(stmt.condition)) {
		i.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		i.execute(stmt.elseBranch)
	}
	return nil
}

// Throw a Return value up the call stack to be caught by function call
func (i *Interpreter) visitReturnStmt(stmt *ReturnStmt) any {
	if stmt.value != nil {
		panic(Return{value: i.evaluate(stmt.value)})
	}

	panic(Return{value: nil})
}

// Define (default to nil) a new variable in the current scope
func (i *Interpreter) visitVarStmt(stmt *VarStmt) any {
	var value any
	if stmt.initializer != nil {
		value = i.evaluate(stmt.initializer)
	}

	i.environment.define(stmt.name.lexeme, value)
	return nil
}

// While condition is true, execute the body
func (i *Interpreter) visitWhileStmt(stmt *WhileStmt) any {
	for isTruthy(i.evaluate(stmt.condition)) {
		i.execute(stmt.body)
	}
	return nil
}

// Assign variable to new value
func (i *Interpreter) visitAssignExpr(expr *AssignExpr) any {
	value := i.evaluate(expr.value)

	if distance, found := i.locals[expr]; found {
		i.environment.assignAt(distance, expr.name, value)
	} else {
		i.globals.assign(expr.name, value)
	}

	return value
}

// Perform (arithmetic/comparison/concatenation) operation on two values
func (i *Interpreter) visitBinaryExpr(expr *BinaryExpr) any {
	left := i.evaluate(expr.left)
	right := i.evaluate(expr.right)

	switch expr.operator.tokenType {
	case GREATER:
		// Greater than
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l > r
			}
		}
		panic(RuntimeError{token: expr.operator, message: "Operands must be numbers."})

	case GREATER_EQUAL:
		// Greater than or equal to
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l >= r
			}
		}
		panic(RuntimeError{token: expr.operator, message: "Operands must be numbers."})

	case LESS:
		// Less than
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l < r
			}
		}
		panic(RuntimeError{token: expr.operator, message: "Operands must be numbers."})

	case LESS_EQUAL:
		// Less than or equal to
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l <= r
			}
		}
		panic(RuntimeError{token: expr.operator, message: "Operands must be numbers."})

	case BANG_EQUAL:
		// Not equal
		return !compare(left, right)

	case EQUAL_EQUAL:
		// Equal
		return compare(left, right)

	case MINUS:
		fallthrough
	case MINUS_EQUAL:
		fallthrough
	case MINUS_MINUS:
		// Subtraction
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l - r
			}
		}
		panic(RuntimeError{token: expr.operator, message: "Operands must be numbers."})

	case PLUS:
		fallthrough
	case PLUS_EQUAL:
		fallthrough
	case PLUS_PLUS:
		// Addition
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l + r
			}
		}
		// Concatenation
		if l, ok := left.(Sequence); ok {
			if r, ok := right.(Sequence); ok {
				if l.isString == r.isString {
					newList := append(append([]any{}, l.list...), r.list...)
					return Sequence{list: newList, isString: l.isString}
				}
			}
		}
		panic(RuntimeError{token: expr.operator, message: "Operands must be two numbers, two strings, or two lists."})

	case SLASH:
		fallthrough
	case SLASH_EQUAL:
		// Division
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				if r == 0 {
					// panic(RuntimeError{token: expr.operator, message: "Division by zero."})
				}
				return l / r
			}
		}
		panic(RuntimeError{token: expr.operator, message: "Operands must be numbers."})

	case STAR:
		fallthrough
	case STAR_EQUAL:
		// Multiplication
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l * r
			}
		}
		panic(RuntimeError{token: expr.operator, message: "Operands must be numbers."})
	}

	// Unreachable
	panic(RuntimeError{token: expr.operator, message: "Unrecognized binary operator."})
}

// Helper function for Interpreter that compares simple values or Sequences
// Functions, Classes, and Instances are passed around by pointer, so they do not need extra handling
func compare(left any, right any) bool {
	// slices are not comparable, so Sequences must be handled separately
	if l, ok := left.(Sequence); ok {
		if r, ok := right.(Sequence); ok {
			if l.isString != r.isString || l.size() != r.size() {
				// left and right are incomparable Sequences (different size or types)
				return false
			}
			for i := 0; i < l.size(); i++ {
				if !compare(l.list[i], r.list[i]) {
					// left and right are comparable Sequences, some elements are not equal
					return false
				}
				// left and right are comparable Sequences, all elements are equal
				return true
			}
		}
		// left is Sequence, right is not
		return false
	} else if _, ok := right.(Sequence); ok {
		// right is Sequence, left is not
		return false
	}

	// left and right are not Sequences, safe to equal
	return left == right
}

// Perform a call on a Callable
func (i *Interpreter) visitCallExpr(expr *CallExpr) any {
	callee := i.evaluate(expr.callee)

	arguments := []any{}
	for _, argument := range expr.arguments {
		arguments = append(arguments, i.evaluate(argument))
	}

	if callable, ok := callee.(Callable); ok {
		// Throw error if the arity doesn't match
		if len(arguments) != callable.arity() {
			panic(RuntimeError{
				token: expr.paren,
				message: "Expected " +
					fmt.Sprint(callable.arity()) + " arguments but got " +
					fmt.Sprint(len(arguments)) + ".",
			})
		}

		// Native functions don't have access to tokens...
		if _, ok := callable.(*Native); ok {
			// ...so set up a defered function to add tokens to Runtime errors
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(RuntimeError); ok {
						panic(RuntimeError{token: expr.paren, message: err.message})
					} else {
						panic(r)
					}
				}
			}()
		}

		return callable.call(i, arguments)
	}

	panic(RuntimeError{token: expr.paren,
		message: "Can only call functions and classes."})
}

// Get value of instance property
func (i *Interpreter) visitGetExpr(expr *GetExpr) any {
	object := i.evaluate(expr.object)
	if instance, ok := object.(*Instance); ok {
		return instance.get(expr.name)
	}

	panic(RuntimeError{token: expr.name,
		message: "Only instances have properties."})
}

// Parentheses
func (i *Interpreter) visitGroupingExpr(expr *GroupingExpr) any {
	return i.evaluate(expr.expression)
}

// Get index or slice copy of a Sequence
func (i *Interpreter) visitIndexExpr(expr *IndexExpr) any {
	indexee := i.evaluate(expr.indexee)

	// Only try indexing on a Sequence
	if sequence, ok := indexee.(Sequence); ok {
		start := i.evaluate(expr.start)
		// Only continue indexing if the start is a number or is omitted
		if startF, ok := start.(float64); ok || start == nil {
			startI := 0
			if ok {
				startI = int(startF)
				// Negative indexing
				if startI < 0 {
					startI += sequence.size()
				}
			}

			// Index operation, because there is no stop index
			if expr.stop == nil {
				if startI < 0 || startI >= sequence.size() {
					// Out of range case
					return Sequence{
						list:     []any{},
						isString: sequence.isString,
					}
				} else if sequence.isString {
					// String case
					return Sequence{
						list:     []any{sequence.list[startI]},
						isString: true,
					}
				} else {
					// List case
					return sequence.list[startI]
				}
			}

			// Slice operation
			if startI < 0 {
				startI = 0
			}
			stop := i.evaluate(expr.stop)
			// Only continue slicing if the stop is a number or is omitted
			if stopF, ok := stop.(float64); ok || stop == nil {
				stopI := sequence.size()
				if ok {
					stopI = int(stopF)
					if stopI < 0 {
						stopI += sequence.size()
					}
					if stopI > sequence.size() {
						stopI = sequence.size()
					}
				}

				if stopI <= startI {
					// Out of range case
					return Sequence{
						list:     []any{},
						isString: sequence.isString,
					}
				} else {
					// Normal case
					// Get shallow copy of the sequence (Instance is copied by reference)
					return Sequence{list: append([]any{}, sequence.list[startI:stopI]...),
						isString: sequence.isString}
				}
			}
		}

		panic(RuntimeError{token: expr.bracket,
			message: "Indices must be numbers."})
	}

	panic(RuntimeError{token: expr.bracket,
		message: "Can only index strings and lists."})
}

// Create a new list
func (i *Interpreter) visitListExpr(expr *ListExpr) any {
	elements := []any{}
	for _, element := range expr.elements {
		elements = append(elements, i.evaluate(element))
	}
	return Sequence{list: elements, isString: false}
}

// A literal value that needs no additional evaluation
func (i *Interpreter) visitLiteralExpr(expr *LiteralExpr) any {
	return expr.value
}

// A logical control-flow operation on two values
func (i *Interpreter) visitLogicalExpr(expr *LogicalExpr) any {
	left := i.evaluate(expr.left)

	switch expr.operator.tokenType {
	case AND:
		if !isTruthy(left) {
			return left
		}
	case OR:
		if isTruthy(left) {
			return left
		}
	default:
		// Unreachable
		panic(RuntimeError{token: expr.operator, message: "Unrecognized logical operator."})
	}

	return i.evaluate(expr.right)
}

// Replace an element of a Sequence at a given index
func (i *Interpreter) visitReplaceExpr(expr *ReplaceExpr) any {
	indexee := i.evaluate(expr.indexee)

	// Only try indexing on a Sequence
	if sequence, ok := indexee.(Sequence); ok {
		value := i.evaluate(expr.value)

		index := i.evaluate(expr.index)
		// Only continue indexing if the index is a number
		if indexF, ok := index.(float64); ok {
			indexI := int(indexF)
			if indexI < 0 {
				indexI += sequence.size()
			}
			if indexI < 0 || indexI >= sequence.size() {
				panic(RuntimeError{token: expr.bracket,
					message: "Index out of range."})
			}

			if sequence.isString {
				// Only replace character if the value is a string of length 1
				if char, ok := value.(Sequence); ok &&
					char.isString && char.size() == 1 {
					sequence.list[indexI] = char.list[0]
					return value
				}

				panic(RuntimeError{token: expr.bracket,
					message: "Replace value must be string of length 1."})
			} else {
				// Replace list element no matter what
				sequence.list[indexI] = value
				return value
			}
		}

		panic(RuntimeError{token: expr.bracket,
			message: "Index must be a number."})
	}

	panic(RuntimeError{token: expr.bracket,
		message: "Can only index strings and lists."})
}

// Set value of instance field
func (i *Interpreter) visitSetExpr(expr *SetExpr) any {
	object := i.evaluate(expr.object)

	if instance, ok := object.(*Instance); ok {
		value := i.evaluate(expr.value)
		instance.set(expr.name, value)
		return value
	}

	panic(RuntimeError{token: expr.name,
		message: "Only instances have fields."})
}

// If condition is true, return trueValue, otherwise falseValue
func (i *Interpreter) visitTernaryExpr(expr *TernaryExpr) any {
	if isTruthy(i.evaluate(expr.condition)) {
		return i.evaluate(expr.trueValue)
	} else {
		return i.evaluate(expr.falseValue)
	}
}

// Special object referring to the current instance
func (i *Interpreter) visitThisExpr(expr *ThisExpr) any {
	return i.lookUpVariable(expr.Token, expr)
}

// Perform an operation on one value
func (i *Interpreter) visitUnaryExpr(expr *UnaryExpr) any {
	right := i.evaluate(expr.operand)

	switch expr.operator.tokenType {
	case BANG:
		return !isTruthy(right)
	case MINUS:
		if r, ok := right.(float64); ok {
			return -r
		}
		panic(RuntimeError{token: expr.operator, message: "Operand must be a number."})
	case PLUS:
		if r, ok := right.(float64); ok {
			return r
		}
		panic(RuntimeError{token: expr.operator, message: "Operand must be a number."})
	}

	// Unreachable
	panic(RuntimeError{token: expr.operator, message: "Unrecognized unary operator."})
}

// Variable name
func (i *Interpreter) visitVariableExpr(expr *VariableExpr) any {
	return i.lookUpVariable(expr.Token, expr)
}
