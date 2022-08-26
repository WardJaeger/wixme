// Ward Jaeger, CS 403
package main

// Any object that can be called with parentheses
type Callable interface {
	toString() string                                   // string representation
	arity() int                                         // number of paremeters
	call(interpreter *Interpreter, arguments []any) any // functionality of the call
}
