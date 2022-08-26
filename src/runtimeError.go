// Ward Jaeger, CS 403
package main

// Struct to indicate an error during interpretation
type RuntimeError struct {
	token   Token
	message string
}
