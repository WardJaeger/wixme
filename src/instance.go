// Ward Jaeger, CS 403
package main

// A particular instantiation of a class
type Instance struct {
	*Class
	fields map[string]any
}

// String representation
func (i *Instance) toString() string {
	return i.name + " instance"
}

// Get property or bound method
func (i *Instance) get(name Token) any {
	if property, found := i.fields[name.lexeme]; found {
		return property
	}

	if method := i.findMethod(name.lexeme); method != nil {
		return method.bind(i)
	}

	panic(RuntimeError{token: name,
		message: "Undefined property '" + name.lexeme + "'."})
}

// Set property
func (i *Instance) set(name Token, value any) {
	i.fields[name.lexeme] = value
}
