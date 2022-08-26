// Ward Jaeger, CS 403
package main

// Data type that can be indexed/sliced/concatenated
type Sequence struct {
	list     []any
	isString bool // Whether the Sequence is a string or a list
}

// Easy access to number of elements
func (s *Sequence) size() int {
	return len(s.list)
}
