// Ward Jaeger, CS 403
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

// Boolean that keeps track of whether an error has occured anywhere, from parsing to interpretation
var hadError = false

// Interpreter class that gets initialized throughout the main file
var mainInterpreter = Interpreter{}

// Entry point for the entire class
func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: wixme [script]")
		os.Exit(1)
	} else {
		setUpInterpreter(&mainInterpreter)

		if len(os.Args) == 2 {
			runFile(os.Args[1])
		} else {
			runPrompt()
		}
	}
}

// Sets up the interpreter with fresh environments and native functions
func setUpInterpreter(interpreter *Interpreter) {
	interpreter.environment = &Environment{values: map[string]any{}}
	interpreter.globals = interpreter.environment

	interpreter.globals.define("clock", &Native{
		arityFunc: func() int { return 0 },
		callFunc: func(_ *Interpreter, _ []any) any {
			return float64(time.Now().UnixNano()) / 1000000000
		},
	})
	interpreter.globals.define("len", &Native{
		arityFunc: func() int { return 1 },
		callFunc: func(_ *Interpreter, args []any) any {
			if sequence, ok := args[0].(Sequence); ok {
				return float64(len(sequence.list))
			}
			panic(RuntimeError{message: "Expect string or list."})
		},
	})
	interpreter.globals.define("print", &Native{
		arityFunc: func() int { return 1 },
		callFunc: func(_ *Interpreter, args []any) any {
			fmt.Println(stringify(args[0], false))
			return nil
		},
	})
	interpreter.globals.define("toNumber", &Native{
		arityFunc: func() int { return 1 },
		callFunc: func(_ *Interpreter, args []any) any {
			if sequence, ok := args[0].(Sequence); ok && sequence.isString {
				str := stringify(sequence, false)
				i := 0

				// Throw error if conversion fails for any reason
				defer func() {
					if r := recover(); r != nil {
						panic(RuntimeError{message: "Invalid format."})
					}
				}()

				// Check for sign
				sign := 1.0
				if str[0] == '+' {
					i++
				} else if str[0] == '-' {
					sign = -1.0
					i++
				}

				// Check for whole part
				if !isDigit(str[i]) {
					panic(0)
				}
				i++

				// Loop through rest of string
				fraction := false
				for i < len(str) {
					if !fraction && str[i] == '.' {
						fraction = true
						i++
					}
					if !isDigit(str[i]) {
						panic(0)
					}
					i++
				}

				value, _ := strconv.ParseFloat(str, 64)
				return value * sign
			}
			panic(RuntimeError{message: "Expect string."})
		},
	})
	interpreter.globals.define("toString", &Native{
		arityFunc: func() int { return 1 },
		callFunc: func(_ *Interpreter, args []any) any {
			stringified := stringify(args[0], false)
			str := []any{}
			for i := 0; i < len(stringified); i++ {
				str = append(str, stringified[i])
			}
			return Sequence{list: str, isString: true}
		},
	})

	interpreter.locals = map[Expr]int{}
}

// Run on the input from a given file
func runFile(filename string) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Could not open file " + filename)
		os.Exit(1)
	}
	run(src)

	if hadError {
		os.Exit(1)
	}
}

// Run in interactive mode from the terminal
func runPrompt() {
	stdin := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for stdin.Scan() {
		line := stdin.Bytes()
		run(line)
		hadError = false
		fmt.Print("> ")
	}
}

// Using a given source, do parse through interpret
func run(source []byte) {
	scanner := Scanner{source: source, startChar: 0, currChar: 0, line: 1}
	tokens := scanner.scanTokens()
	parser := Parser{tokens: tokens, current: 0}
	statements := parser.parse()

	// Stop if there was a syntax error.
	if hadError {
		return
	}

	resolver := Resolver{interpreter: &mainInterpreter}
	resolver.resolve(statements)

	// Stop if there was a resolution error.
	if hadError {
		return
	}

	mainInterpreter.interpret(statements)
}

// Report a given lexeme (mainly for scanner)
func reportLexeme(line int, col int, lexeme string, message string) {
	report(line, col, " at '"+lexeme+"'", message)
}

// Report a given token (mainly for parser)
func reportToken(token Token, message string) {
	if token.tokenType == EOF {
		report(token.line, token.col, " at EOF", message)
	} else {
		report(token.line, token.col, " at '"+token.lexeme+"'", message)
	}
}

// Report a given Runtime error (mainly for interpreter)
func reportRuntime(err RuntimeError) {
	report(err.token.line, err.token.col, " at '"+err.token.lexeme+"' during runtime", err.message)
}

// Output an error message, noting the line and column
func report(line int, col int, where string, message string) {
	fmt.Println("[line " + strconv.Itoa(line) + ", col " + strconv.Itoa(col) +
		"] Error" + where + ": " + message)
	hadError = true
}
