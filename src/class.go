// Ward Jaeger, CS 403
package main

// A bundle of data with methods that operate on it, created by an initizalizer
type Class struct {
	name    string
	methods map[string]*Function
}

// Test for interface implementation
var _ Callable = &Class{}

func (c *Class) toString() string {
	return c.name
}

func (c *Class) arity() int {
	// Return arity of initializer if it exists, or 0 otherwise
	if initializer := c.findMethod("init"); initializer != nil {
		return initializer.arity()
	}
	return 0
}

func (c *Class) call(interpreter *Interpreter, arguments []any) any {
	// Create a new instance of this class, and initialize it
	instance := &Instance{Class: c, fields: map[string]any{}}
	if initializer := c.findMethod("init"); initializer != nil {
		initializer.bind(instance).call(interpreter, arguments)
	}

	return instance
}

func (c *Class) findMethod(name string) *Function {
	// Return a method with the given name
	if function, found := c.methods[name]; found {
		return function
	}

	return nil
}
