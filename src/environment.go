// Ward Jaeger, CS 403
package main

// A location of variable storage, potentially enclosed in a parent environmnet
type Environment struct {
	enclosing *Environment
	values    map[string]any
}

// Get parent environment at a certain distance
func (e *Environment) ancestor(distance int) *Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}
	return environment
}

// Assign value to the current scope of variable name
func (e *Environment) assign(name Token, value any) {
	if _, prs := e.values[name.lexeme]; prs {
		e.values[name.lexeme] = value
		return
	}

	panic(RuntimeError{token: name, message: "Undefined variable '" + name.lexeme + "'."})
}

// Assign value to parent scope at certain distance
func (e *Environment) assignAt(distance int, name Token, value any) {
	e.ancestor(distance).values[name.lexeme] = value
}

// Define a new variable in this scope with a given initial value
func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

// Get value from most recent scope of variable name
func (e *Environment) get(name Token) any {
	if value, prs := e.values[name.lexeme]; prs {
		return value
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	panic(RuntimeError{token: name, message: "Undefined variable '" + name.lexeme + "'."})
}

// Get value of parent scope at certain distance
func (e *Environment) getAt(distance int, name string) any {
	return e.ancestor(distance).values[name]
}
