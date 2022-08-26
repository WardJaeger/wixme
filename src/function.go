// Ward Jaeger, CS 403
package main

// User-defined function
type Function struct {
	declaration   *FunctionStmt
	closure       *Environment
	isInitializer bool
}

// Test for interface implementation
var _ Callable = &Function{}

// Return a new function that is bound to a specific instance
func (f *Function) bind(instance *Instance) *Function {
	currEnvironment := &Environment{enclosing: f.closure,
		values: map[string]any{}}
	currEnvironment.define("this", instance)
	return &Function{declaration: f.declaration,
		closure: currEnvironment, isInitializer: f.isInitializer}
}

func (n *Function) toString() string {
	return "<fn " + n.declaration.name.lexeme + ">"
}

func (f *Function) arity() int {
	return len(f.declaration.params)
}

func (f *Function) call(i *Interpreter, arguments []any) (returnValue any) {
	currEnvironment := &Environment{enclosing: f.closure, values: map[string]any{}}
	for i, param := range f.declaration.params {
		currEnvironment.define(param.lexeme, arguments[i])
	}

	// Set up a defered function to catch a return value from the body
	defer func() {
		if r := recover(); r != nil {
			if caughtValue, ok := r.(Return); ok {
				if f.isInitializer {
					returnValue = f.closure.getAt(0, "this")
				} else {
					returnValue = caughtValue.value
				}
			} else {
				panic(r)
			}
		}
	}()

	i.executeBlock(f.declaration.body, currEnvironment)

	if f.isInitializer {
		return f.closure.getAt(0, "this")
	}
	return nil
}
