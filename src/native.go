// Ward Jaeger, CS 403
package main

import "fmt"

// WIXME native functions, defined in main.go
type Native struct {
	arityFunc func() int
	callFunc  func(interpreter *Interpreter, arguments []any) any
}

// Test for interface implementation
var _ Callable = &Native{}

func (n *Native) toString() string {
	return "<native fn>"
}

func (n *Native) arity() int {
	return n.arityFunc()
}

func (n *Native) call(interpreter *Interpreter, arguments []any) any {
	return n.callFunc(interpreter, arguments)
}

// Helper function for native function that converts objects into strings
// Nested strings include quotes, isolated strings do not
func stringify(value any, withQuotes bool) string {
	if value == nil {
		return "nil"
	} else if callable, ok := value.(Callable); ok {
		return callable.toString()
	} else if sequence, ok := value.(Sequence); ok {
		// Sequences need to be recursively constructed
		if sequence.isString {
			runes := ""
			for _, element := range sequence.list {
				runes += string(element.(byte))
			}
			if withQuotes {
				return "\"" + runes + "\""
			} else {
				return runes
			}
		} else {
			elements := ""
			for j, element := range sequence.list {
				if j != 0 {
					elements = elements + ", "
				}
				elements = elements + stringify(element, true)
			}
			return "[" + elements + "]"
		}
	} else {
		// Everything else gets converted to Go's default string representation
		return fmt.Sprint(value)
	}
}
