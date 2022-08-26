// Ward Jaeger, CS 403
package main

// Converts a list of bytes into a list of tokens
type Scanner struct {
	source    []byte  // Bytes to scan
	tokens    []Token // Tokens created
	startChar int     // Starting index of current token
	currChar  int     // Index of current byte
	line      int     // Line number
	colStart  int     // Index of first byte in the line
}

// List of keywords and the tokens that they evaluate to
var keywords = map[string]tokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"let":    LET,
	"nil":    NIL,
	"or":     OR,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

// Entry point for scanning
func (s *Scanner) scanTokens() []Token {
	for !s.isAtEnd() {
		s.scanToken()
	}

	s.startChar = s.currChar
	s.tokens = append(s.tokens, Token{tokenType: EOF, line: s.line,
		col: s.getCol()})
	return s.tokens
}

// Generate the next token
func (s *Scanner) scanToken() {
	s.startChar = s.currChar

	c := s.advance()
	switch c {
	// Ignore all whitespace
	case ' ':
	case '\r':
	case '\t':
	case '\n':

	// Single character tokens
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case '[':
		s.addToken(LEFT_BRACKET)
	case ']':
		s.addToken(RIGHT_BRACKET)
	case ':':
		s.addToken(COLON)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '?':
		s.addToken(QUESTION)
	case ';':
		s.addToken(SEMICOLON)

	// Multi-character tokens
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '-':
		if s.match('=') {
			s.addToken(MINUS_EQUAL)
		} else if s.match('-') {
			s.addToken(MINUS_MINUS)
		} else {
			s.addToken(MINUS)
		}
	case '+':
		if s.match('=') {
			s.addToken(PLUS_EQUAL)
		} else if s.match('+') {
			s.addToken(PLUS_PLUS)
		} else {
			s.addToken(PLUS)
		}
	case '/':
		if s.match('/') {
			// Single-line comment, advance to end of line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			// Multiline comment, begin recursive scan
			s.multilineComment()
		} else if s.match('=') {
			s.addToken(SLASH_EQUAL)
		} else {
			s.addToken(SLASH)
		}
	case '*':
		if s.match('=') {
			s.addToken(STAR_EQUAL)
		} else {
			s.addToken(STAR)
		}
	case '"':
		s.string()

	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			reportLexeme(s.line, s.getCol(), string(c), "Unexpected character "+string(c)+".")
		}
	}
}

// Recursively scan multiline comment
func (s *Scanner) multilineComment() {
	s.startChar = s.currChar - 2
	startLine := s.line
	startCol := s.getCol()

	for !s.isAtEnd() {
		switch s.advance() {
		case '/':
			if s.match('*') {
				// Nested multiline comment, recurse deeper
				s.multilineComment()
			}
		case '*':
			if s.match('/') {
				// Close the current comment
				return
			}
		}
	}

	reportLexeme(startLine, startCol, "/*", "Unterminated multiline comment.")
}

// Generate a token of a given confirmed type
func (s *Scanner) addToken(ttype tokenType) {
	// Select lexeme from the source
	text := string(s.source[s.startChar:s.currChar])
	s.tokens = append(s.tokens, Token{tokenType: ttype, lexeme: text,
		line: s.line, col: s.getCol()})
}

// Scan all characters associated with a string and generate a token
func (s *Scanner) string() {
	for true {
		c := s.advance()
		if c == '"' {
			break
		} else if c == '\\' {
			s.match('"')
		} else if c == '\n' || s.isAtEnd() {
			reportLexeme(s.line, s.getCol(), "\"", "Unterminated string.")
			return
		}
	}

	s.addToken(STRING)
}

// Scan all characters associated with a number and generate a token
func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part.
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	s.addToken(NUMBER)
}

// Scan all characters associated with an identifier and generate a token
func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	lexeme := string(s.source[s.startChar:s.currChar])

	// Check if the lexeme is a keyword, otherwise set type to IDENTIFIER
	ttype, pres := keywords[lexeme]
	if pres == false {
		ttype = IDENTIFIER
	}
	s.addToken(ttype)
}

// Check if there are no more bytes to scan
func (s *Scanner) isAtEnd() bool {
	return s.currChar >= len(s.source)
}

// Look at next byte
func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return byte(0)
	}
	return s.source[s.currChar]
}

// Look ahead two bytes
func (s *Scanner) peekNext() byte {
	if s.currChar+1 >= len(s.source) {
		return byte(0)
	}
	return s.source[s.currChar+1]
}

// Move forward and return the passed byte, also noting line breaks
func (s *Scanner) advance() byte {
	r := s.peek()
	s.currChar++
	if r == '\n' {
		s.colStart = s.currChar
		s.line++
	}
	return r
}

// Advance only for a given byte, also noting line breaks
func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() || s.source[s.currChar] != expected {
		return false
	}

	s.advance()
	return true
}

// Get the column of the current token
func (s *Scanner) getCol() int {
	return s.startChar - s.colStart + 1
}

// Helper function, check if byte is digit
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// Helper function, check if byte is alphabetical or underscore
func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

// Helper function, check if byte is alphabetical, underscore, or digit
func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}
