// Ward Jaeger, CS 403
package main

// Define "enum" type
type tokenType string

// "Enum" of all possible token types
// Go does not actually have enums, but constant values work as a substitute
const (
	// Single-character tokens.
	LEFT_PAREN    tokenType = "LEFT_PAREN"
	RIGHT_PAREN   tokenType = "RIGHT_PAREN"
	LEFT_BRACE    tokenType = "LEFT_BRACE"
	RIGHT_BRACE   tokenType = "RIGHT_BRACE"
	LEFT_BRACKET  tokenType = "LEFT_BRACKET"
	RIGHT_BRACKET tokenType = "RIGHT_BRACKET"
	COLON         tokenType = "COLON"
	COMMA         tokenType = "COMMA"
	DOT           tokenType = "DOT"
	QUESTION      tokenType = "QUESTION"
	SEMICOLON     tokenType = "SEMICOLON"

	// One or two character tokens.
	BANG          tokenType = "BANG"
	BANG_EQUAL    tokenType = "BANG_EQUAL"
	EQUAL         tokenType = "EQUAL"
	EQUAL_EQUAL   tokenType = "EQUAL_EQUAL"
	GREATER       tokenType = "GREATER"
	GREATER_EQUAL tokenType = "GREATER_EQUAL"
	LESS          tokenType = "LESS"
	LESS_EQUAL    tokenType = "LESS_EQUAL"
	MINUS         tokenType = "MINUS"
	MINUS_EQUAL   tokenType = "MINUS_EQUAL"
	MINUS_MINUS   tokenType = "MINUS_MINUS"
	PLUS          tokenType = "PLUS"
	PLUS_EQUAL    tokenType = "PLUS_EQUAL"
	PLUS_PLUS     tokenType = "PLUS_PLUS"
	SLASH         tokenType = "SLASH"
	SLASH_EQUAL   tokenType = "SLASH_EQUAL"
	STAR          tokenType = "STAR"
	STAR_EQUAL    tokenType = "STAR_EQUAL"

	// Literals.
	IDENTIFIER tokenType = "IDENTIFIER"
	STRING     tokenType = "STRING"
	NUMBER     tokenType = "NUMBER"

	// Keywords.
	AND    tokenType = "AND"
	BREAK  tokenType = "BREAK"
	CLASS  tokenType = "CLASS"
	ELSE   tokenType = "ELSE"
	FALSE  tokenType = "FALSE"
	FUN    tokenType = "FUN"
	FOR    tokenType = "FOR"
	IF     tokenType = "IF"
	LET    tokenType = "LET"
	NIL    tokenType = "NIL"
	OR     tokenType = "OR"
	RETURN tokenType = "RETURN"
	SUPER  tokenType = "SUPER"
	THIS   tokenType = "THIS"
	TRUE   tokenType = "TRUE"
	VAR    tokenType = "VAR"
	WHILE  tokenType = "WHILE"

	EOF tokenType = "EOF"
)
