// Ward Jaeger, CS 403
package main

// Struct to contain information about a scanned lexeme
type Token struct {
	tokenType tokenType // Category of token
	lexeme    string    // Actual string
	line      int       // Line it was found on
	col       int       // Column it started on
}
