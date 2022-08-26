// Ward Jaeger, CS 403
package main

import "strconv"

// Converts a list of tokens into an AST
type Parser struct {
	tokens  []Token // Tokens to parse
	current int     // Index of current token
}

// Entry point to begin parsing tokens
func (p *Parser) parse() []Stmt {
	statements := []Stmt{}
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}

// Any type of declaration or statement
func (p *Parser) declaration() Stmt {
	// Set up a deferred function that handles Parse errors
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(ParseError); ok {
				p.synchronize()
			} else {
				panic(r)
			}
		}
	}()

	if p.match(CLASS) {
		return p.classDeclaration()
	}
	if p.match(FUN) {
		return p.function("function")
	}
	if p.match(VAR) {
		// Check for terminator after variable declaration
		defer p.terminator("Expect terminator after variable declaration.")
		return p.varDeclaration()
	}
	return p.statement()
}

// Declare a new class
func (p *Parser) classDeclaration() *ClassStmt {
	name := p.consume(IDENTIFIER, "Expect class name.")
	p.consume(LEFT_BRACE, "Expect '{' before class body.")

	methods := []*FunctionStmt{}
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		methods = append(methods, p.function("method"))
	}

	p.consume(RIGHT_BRACE, "Expect '}' after class body.")

	return &ClassStmt{name: name, methods: methods}
}

// Parse some sort of function
func (p *Parser) function(kind string) *FunctionStmt {
	name := p.consume(IDENTIFIER, "Expect "+kind+" name.")
	p.consume(LEFT_PAREN, "Expect '(' after "+kind+" name.")
	parameters := []Token{}
	if !p.check(RIGHT_PAREN) {
		parameters = append(parameters, p.consume(IDENTIFIER, "Expect paramter name."))
		for p.match(COMMA) {
			parameters = append(parameters, p.consume(IDENTIFIER, "Expect parameter name."))
		}
	}
	p.consume(RIGHT_PAREN, "Expect ')' after parameters.")

	p.consume(LEFT_BRACE, "Expect '{' before "+kind+" body.")
	body := p.block()
	return &FunctionStmt{name: name, params: parameters, body: body}
}

// Declare a new variable
func (p *Parser) varDeclaration() *VarStmt {
	name := p.consume(IDENTIFIER, "Expect variable name.")

	var initializer Expr
	if p.match(EQUAL) {
		initializer = p.expression()
	}

	return &VarStmt{name: name, initializer: initializer}
}

// Get some other kind of statement
func (p *Parser) statement() Stmt {
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(RETURN) {
		// Check for terminator after return statement
		defer p.terminator("Expect terminator after return value.")
		return p.returnStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(LEFT_BRACE) {
		return &BlockStmt{statements: p.block()}
	}

	// Check for terminator after expression statement
	defer p.terminator("Expect terminator after expression.")
	return p.expressionStatement()
}

// For statement, just syntactic sugar
func (p *Parser) forStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer Stmt
	if p.match(VAR) {
		initializer = p.varDeclaration()
	} else if !p.check(SEMICOLON) {
		initializer = p.expressionStatement()
	}
	// Note: Semicolon is necessary here
	p.consume(SEMICOLON, "Expect ';' after initializer statement.")

	var condition Expr
	if !p.check(SEMICOLON) {
		condition = p.expression()
	}
	// Note: Semicolon is necessary here
	p.consume(SEMICOLON, "Expect ';' after loop condition.")

	var increment Stmt
	if !p.check(RIGHT_PAREN) {
		increment = p.expressionStatement()
	}
	p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")

	body := p.statement()

	// Body always ends with the increment
	if increment != nil {
		body = &BlockStmt{statements: []Stmt{body, increment}}
	}

	// Body is executed when condition is true
	if condition == nil {
		condition = &LiteralExpr{value: true}
	}
	body = &WhileStmt{condition: condition, body: body}

	// Initializer kicks the whole thing off
	if initializer != nil {
		body = &BlockStmt{statements: []Stmt{initializer, body}}
	}

	return body
}

// If statement
func (p *Parser) ifStatement() *IfStmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch = p.statement()
	}

	return &IfStmt{condition: condition, thenBranch: thenBranch, elseBranch: elseBranch}
}

// Return statement
func (p *Parser) returnStatement() (returnValue *ReturnStmt) {
	keyword := p.previous()
	currToken := p.current

	// Set up a deferred function that catches a parse error
	// If a valid expression is not found, return nil
	// This ensures maximal munch for return statments
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(ParseError); ok {
				p.current = currToken
				returnValue = &ReturnStmt{keyword: keyword, value: nil}
			} else {
				panic(r)
			}
		}
	}()

	return &ReturnStmt{keyword: keyword, value: p.expression()}
}

// While statement
func (p *Parser) whileStatement() *WhileStmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	body := p.statement()

	return &WhileStmt{condition: condition, body: body}
}

// Parse a list of statements
func (p *Parser) block() []Stmt {
	statements := []Stmt{}

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

// Convert any expression into a statement
func (p *Parser) expressionStatement() *ExpressionStmt {
	return &ExpressionStmt{expression: p.expression()}
}

// Expression (reduces immediately to assignment)
func (p *Parser) expression() Expr {
	return p.assignment()
}

// Assignment
func (p *Parser) assignment() Expr {
	expr := p.ternary()

	if p.match(EQUAL, MINUS_EQUAL, PLUS_EQUAL, SLASH_EQUAL, STAR_EQUAL) {
		equals := p.previous()
		value := p.assignment()
		// Compound assignment operators get expanded
		if equals.tokenType != EQUAL {
			value = &BinaryExpr{left: expr, operator: equals, right: value}
		}
		return p.finishAssignment(expr, equals, value)
	}

	return expr
}

// Verify assignment target, and complete assignment expression
func (*Parser) finishAssignment(target Expr, operator Token, value Expr) Expr {
	// Check for correct assignment target, as indicated by the grammar
	if name, ok := target.(*VariableExpr); ok {
		return &AssignExpr{name: name.Token, value: value}
	} else if get, ok := target.(*GetExpr); ok {
		return &SetExpr{object: get.object, name: get.name, value: value}
	} else if index, ok := target.(*IndexExpr); ok && index.stop == nil {
		return &ReplaceExpr{indexee: index.indexee, index: index.start,
			bracket: index.bracket, value: value}
	}

	reportToken(operator, "Invalid assignment target.")
	return target
}

// Ternary operator
func (p *Parser) ternary() Expr {
	expr := p.or()

	for p.match(QUESTION) {
		operator := p.previous()
		trueValue := p.or()
		p.consume(COLON, "Expect ':' after expression.")
		falseValue := p.or()
		expr = &TernaryExpr{condition: expr, operator: operator,
			trueValue: trueValue, falseValue: falseValue}
	}

	return expr
}

// Logical or expression
func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = &LogicalExpr{left: expr, operator: operator, right: right}
	}

	return expr
}

// Logical and expression
func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		expr = &LogicalExpr{left: expr, operator: operator, right: right}
	}

	return expr
}

// Equality operation
func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &BinaryExpr{left: expr, operator: operator, right: right}
	}

	return expr
}

// Comparison operation
func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &BinaryExpr{left: expr, operator: operator, right: right}
	}

	return expr
}

// Addition, subtraction, concatenation
func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &BinaryExpr{left: expr, operator: operator, right: right}
	}

	return expr
}

// Multiplication, division
func (p *Parser) factor() Expr {
	expr := p.prefix()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.prefix()
		expr = &BinaryExpr{left: expr, operator: operator, right: right}
	}

	return expr
}

// Unary prefix operations
func (p *Parser) prefix() Expr {
	if p.match(BANG, MINUS, PLUS) {
		operator := p.previous()
		right := p.prefix()
		return &UnaryExpr{operator: operator, operand: right}
	}

	return p.increment()
}

// Increment/decrement postfix operations
func (p *Parser) increment() Expr {
	expr := p.postfix()

	if p.match(MINUS_MINUS, PLUS_PLUS) {
		operator := p.previous()
		value := &BinaryExpr{left: expr, operator: operator,
			right: &LiteralExpr{value: 1.0}}
		return p.finishAssignment(expr, operator, value)
	}

	return expr
}

// Instance getting, calls, and list indexing/slicing
func (p *Parser) postfix() Expr {
	expr := p.primary()

	for {
		if p.match(DOT) {
			name := p.consume(IDENTIFIER,
				"Expect property name after '.'.")
			expr = &GetExpr{object: expr, name: name}
		} else if p.match(LEFT_PAREN) {
			args := []Expr{}
			if !p.check(RIGHT_PAREN) {
				args = p.arguments()
			}
			paren := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
			expr = &CallExpr{callee: expr, arguments: args, paren: paren}
		} else if p.match(LEFT_BRACKET) {
			expr = p.finishIndex(expr)
		} else {
			break
		}
	}

	return expr
}

// Lowest level, matches to many literal values
func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return &LiteralExpr{value: false}
	} else if p.match(TRUE) {
		return &LiteralExpr{value: true}
	} else if p.match(NIL) {
		return &LiteralExpr{value: nil}
	} else if p.match(NUMBER) {
		value, _ := strconv.ParseFloat(p.previous().lexeme, 64)
		return &LiteralExpr{value: value}
	}

	if p.match(STRING) {
		// Convert to Sequence, and remove the quotation marks
		prev := p.previous()
		str := []any{}
		for i := 1; i < len(prev.lexeme)-1; i++ {
			if prev.lexeme[i] == '\\' {
				switch prev.lexeme[i+1] {
				case 'n':
					str = append(str, byte('\n'))
				case 't':
					str = append(str, byte('\t'))
				case '"':
					str = append(str, byte('"'))
				case '\\':
					str = append(str, byte('\\'))
				default:
					reportToken(prev, "Contains invalid escape sequence '\\"+string(prev.lexeme[i+1])+"'.")
				}
				i++
			} else {
				str = append(str, prev.lexeme[i])
			}
		}
		return &LiteralExpr{value: Sequence{list: str, isString: true}}
	}

	if p.match(THIS) {
		return &ThisExpr{Token: p.previous()}
	}

	if p.match(IDENTIFIER) {
		return &VariableExpr{p.previous()}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return &GroupingExpr{expression: expr}
	}

	if p.match(LEFT_BRACKET) {
		elems := []Expr{}
		if !p.check(RIGHT_BRACKET) {
			elems = p.arguments()
		}
		bracket := p.consume(RIGHT_BRACKET, "Expect ']' after elements.")
		return &ListExpr{elements: elems, bracket: bracket}
	}

	reportToken(p.peek(), "Expect expression.")
	panic(ParseError{})
}

// A list of comma-separated values
func (p *Parser) arguments() []Expr {
	arguments := []Expr{p.expression()}
	for p.match(COMMA) {
		arguments = append(arguments, p.expression())
	}
	return arguments
}

// Indexing or slicing arguments
func (p *Parser) finishIndex(indexee Expr) *IndexExpr {
	var start Expr
	if p.check(COLON) {
		start = &LiteralExpr{value: nil}
	} else {
		start = p.expression()
		if p.check(RIGHT_BRACKET) {
			return &IndexExpr{indexee: indexee, start: start, stop: nil, bracket: p.advance()}
		}
	}

	p.consume(COLON, "Expect ':' or ']' after start index.")

	var stop Expr
	if p.check(RIGHT_BRACKET) {
		stop = &LiteralExpr{value: nil}
	} else {
		stop = p.expression()
	}

	bracket := p.consume(RIGHT_BRACKET, "Expect ']' after stop index.")

	return &IndexExpr{indexee: indexee, start: start, stop: stop, bracket: bracket}
}

// Look at next token
func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

// Look at previous token
func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

// Check if next token is of a given type
func (p *Parser) check(ttype tokenType) bool {
	return p.peek().tokenType == ttype
}

// Check if next token is EOF
func (p *Parser) isAtEnd() bool {
	return p.check(EOF)
}

// Move forward and return the passed token
func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// Advance only if the token is of some given types
func (p *Parser) match(types ...tokenType) bool {
	for _, ttype := range types {
		if p.check(ttype) {
			p.advance()
			return true
		}
	}

	return false
}

// Consume a token of a given type, or throw an error
func (p *Parser) consume(ttype tokenType, errorMessage string) Token {
	if p.check(ttype) {
		return p.advance()
	}

	reportToken(p.peek(), errorMessage)
	panic(ParseError{})
}

// Check that a terminator exists (and consume if it's a semicolon), or throw an error
func (p *Parser) terminator(errorMessage string) {
	if p.previous().line != p.peek().line || p.match(SEMICOLON) ||
		p.check(RIGHT_BRACE) || p.isAtEnd() {
		return
	}

	reportToken(p.peek(), errorMessage)
	panic(ParseError{})
}

// Synchronize the parser to a recognizable state (keyword or past semicolon)
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().tokenType == SEMICOLON {
			return
		}

		switch p.peek().tokenType {
		case CLASS:
			return
		case FUN:
			return
		case VAR:
			return
		case FOR:
			return
		case IF:
			return
		case WHILE:
			return
		case RETURN:
			return
		}

		p.advance()
	}
}
