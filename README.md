# **WIXME**

WIXME is Ward Jaeger's programming language, created as the final project for CS 403. "WIXME" is a portmanteau of my own name "Ward" and the common programming comment "FIXME", suggesting at this language's purpose for quick prototyping. (My girlfriend Zoe Griffin came up with the name, and wanted me to give her credit for it.) This repository, coded in GoLang, serves as the interpreter. I have been developing the repository on Go 1.18, although it may also work on a previous version.

# Justification

WIXME was inspired in part by an experience I had in an internship interview earlier this semester. As I was writing pseudocode to solve some dynamic programming problem, one of the interviewers asked me a question: "Is this code in C? There is some syntax here I don't recognize." That made me pause for a moment and evaluate exactly what I was writing. My pseudocode was in fact quite close to C, but it also drew inspiration from Python (in regards to indexing and slicing) and JavaScript (in regards to omitted semicolons). None of this was a conscious design choice at the time; the syntax just felt quite natural to me.

From this experience, I realized that I didn't know of any good language for prototyping. On the one hand, plain and unstandardized pseudocode is simply untestable, since it cannot be run. But on the other hand, the most common lightweight languages each have various design features that limit their readablity or ease of use when doing rapid prototyping. C is rigid in its requirement of semicolons, uses static typing, and provides minimal support for list operations; Python relies heavily on the proper use of whitespace; and Javascript has strange Automatic Semicolon Insertion rules and often fails silently. To be sure, these languages have their strengths, but not one is perfect.

I decided for this project to solidify my intuitions about what a language built for prototyping should look like. The result is WIXME, which is a general-purpose, imperative, dynamically-typed, object-oriented language. It attempts to solves the issues found in alternative languages, combining simplicity and flexibility with a wide range of functionality. Specifically, it meets the following minimum requirements I decided were necessary for this purpose:

- **C-style syntax:** The easily recognizable syntax is shared by many common languages.
- **Compound assignment operators:** Reassignments that modify numbers are short and sweet.
- **Dynamic typing:** By not having to specify types, variable declarations look cleaner and more uniform.
- **List operations:** Lists of values can be concatenated, indexed, and sliced, all with simple syntax.
- **String operations:** Strings are treated like lists, with all the same functionality. Common escape sequences are alsoincluded.
- **Optional semicolons:** While semicolons are still allowed as statement terminators, the parser will handle most situations where one is omitted.

# Language Guide

WIXME is based heavily on Lox, having many of the same features (excluding class inheritance). Lox happend to solve many of the problems noted above, but even it has a some holes that need to be filled. Thus, WIXME has more than a few distinguishing features, which are listed here. The most drastic difference is the introduction of the list data type, which will be explained soon.

## Multiline comments

In addition to single-line comments, WIXME supports multiline comments. These begin with a forward slash and an asterisk, and they end with an asterisk and a forward-slash. Like Swift,
multiline comments can be nested.

    /*
        This is a
        multiline comment.
        /*
            This is a nested comment.
        */
    */

## Statement termination

In Lox, statements such as variable declarations, expression statements, and return statements must always be terminated with a semicolon. This is not the case in WIXME. Following convention from JavaScript, such statements may also be implicitly terminated with a line break, a right brace, or the end of the file. The only exception to this rule is the `for` clauses, in which the initializer statement and the loop condition must each be followed by a semicolon, and they cannot be implicitly terminated.

    var str = "Hi"; var sum = 0    // Valid
    for (var i = 0                 // Invalid
        i < 10; i = i+1) {         // Valid
        sum = sum + i }            // Valid
    var x = 6       x = x/3        // Invalid
    Foo(str, x)                    // Valid

## Print function

Printing is not its own statement. Rather, it is a native function that takes a single parameter, outputting with a newline.

    print("Hello, world!")

## Lists

In addition to the four data types of Lox (Booleans, numbers, strings, and nil), WIXME has a fifth data type: lists. Similar to how they appear in Python, a list is an ordered, zero-indexed, changable, and heterogenous collection of elements. A list is written as a pair of square brackets containing comma-separated values. Lists have three primary operations:

- **Indexing** returns the element at a given numerical index. It is notated by placing the index in square brackets after an expression that evaluates to the list. If a negative index is used, the element is selected from the end of a list, rather than the beginning. This index is truncated toward zero, regardless of sign.
  
- **Slicing** returns a continuous sublist from a given start index (inclusive, defaults to zero) to a given end index (exclusive, defaults to length of list). It is notated by separating the indices by a colon and placing them in square brackets after an expression that evaluates to a list. Either or both of the indices can be omitted (or set to nil), which casues them to be evaluated as their default values. Slicing also supports negative indexes. Unlike Python, there is no step parameter, because I decided that it was largely unnecessary.

- **Concatenation** combines two existing lists and returns it as a new list. It is notated by the infix operator `+`.

```
var myList = [13.5, "Hi", nil, true, -20]
print(myList[-4])                            // Hi
print(myList[:3])                            // [13.5, "Hi", nil]
print([5, 6, 7, 8][1:3] + myList[2:4])       // [6, 7, nil, true]
```

An index of a list can also be used as an assignment target, changing the element in the original list.

```
var myList2 = [0, 1, 29, 3, 4]
myList2[2] = 2
print(myList2)                    // [0, 1, 2, 3, 4]
```

## String operations

String values can be indexed, sliced, and concatenated as if they were lists of individual characters. Note that indexing a string will return a string of length 1, because WIXME has no character data type. Similarly, an index of a string can only be assigned to a string of length 1.

```
var message = "Hello, world!"
print(message[:5] + message[-1])    // Hello!
message[-1] = "?"
print(message)                      // Hello, world?
```

Lists and strings are strictly distinct types, and they cannot be concatenated together.

## Escape sequences

In a string, the backslash acts as an escape character. There are four valid escape sequences, each of which is parsed as a string of length 1.

- `\n` indicates a newline.
- `\t` indicates a horizontal tab.
- `\\` indicates a backslash.
- `\"` indicates a double quote.

```
print("Hello, \"world\".")       // Hello, "world".
print("How are you?\n\tGood.)    /* How are you?
                                        Good.          */
```

## References

Lists, strings, and instances are passed around as references to locations in memory. By definition, whenever the value of one of these references is modified, the values of the rest are modified too. This includes when a reference is passed as an argument to a function or is indexed from a list.

```
class Foo {
    init(name) {
        this.name = name
    }
}
var fooList = [Foo("Frodo"), Foo("Gollum")]
var myFoo = fooList[1]
myFoo.name = "Sam"
print(fooList[1].name)                         // Sam
```

## Compound assignment operations

Compound assigment operators are a shorthand for updating a variable by performing basic arithmetic or concatenation operations on it. The four compound assignment operators are `+=`, `-=`, `*=`, and `/=`.

```
var number = 10
number += 2         // 12
number /= 4         // 3
number -= 10        // -7
number *= -2        // 14
var name = "Joe"
name += " Dart"     // Joe Dart
```

Related to these are the increment and decrement operations, notated by the postfix operators `++` and `--`, respectively. These reassign a variable by adding one or subtracting one, but with much higher precedence than regular assignment.

```
var number = 10
print(number++)        // 11
print(number-- / 5)    // 2
```

## Ternary operation

The syntax and the functionality of the ternary operator is identical to its appearance in C and in Javascript. A condition is followed by a question mark, a first value, a colon, and a second value.

```
var responseTime = 132
print(responseTime < 100 ? "Completed" : "Timeout")    // Timeout
```

## Unary plus operation

In symmetry with the unary minus operation, WIXME defines the unary plus operation. It can only act on a number, but it returns the number unchanged. It is notated by the prefix operator `+`.

```
print(+100)    // 100
```

## Special numbers

There are a few special values that numbers can take on.

- `0` and `-0` represent signed zero. It retains this sign when doing operations and printing.
- `+Inf` and `-Inf` represent numbers with a sufficiently large magnitude. This can occur, for example, when dividing a nonzero number by zero. Although these infinite values can be operated on, they should generally be used only to indicate of a loss of precision, rather than actual numbers.
- `NaN` represents any number with an indeterminate value. This occurs when dividing zero by zero, or when operating on infinite values in certain mathematically ambiguous ways. Like the infinite values, it is only intended to indicate of a loss of precision, rather than an actual number.

## Other native functions

Beyond `clock` (which is in Lox) and `print` (described above), WIXME has a few other native functions.

- `len` takes a list or a string and returns its length. If the argument is not a list or a string, a runtime error is thrown.
- `toNumber` takes a string and converts it to a number. If the argument is not a string or is not in the format of a number literal, a runtime error is thrown. 
- `toString` takes a single argument and converts it into its string representation, which is how it would look when printed.

```
var list = [0, 1, 2, 3, 4]
print("My list: " + toString(list))        // My list: [0, 1, 2, 3, 4]
print("Length: " + toString(len(list)))    // Length: 5
print(toNumber("7.5") / toNumber("4"))     // 1.875
```

# Grammar

## Syntax

Below is the complete grammatical syntax for WIXME in EBNF form. It follows the principle of maximal munch, which resolves all cases of ambiguity.

```
program         → declaration* EOF

declaration     → classDecl
                | funDecl
                | varDecl TERMINATOR
                | statement

classDecl       → "class" IDENTIFIER "{" function* "}"

funDecl         → "fun" function

function        → IDENTIFIER "(" parameters? ")" block

parameters      → IDENTIFIER ( "," IDENTIFIER )*

varDecl         → "var" IDENTIFIER ( "=" expression )?

statement       → exprStmt TERMINATOR
                | forStmt
                | ifStmt
                | returnStmt TERMINATOR
                | whileStmt
                | block

exprStmt        → expression

forStmt         → "for" "(" ( varDecl | exprStmt )? ";"
                    expression? ";"
                    exprStmt? ")" statement

ifStmt          → "if" "(" expression ")" statement
                    ( "else" statement )?

returnStmt      → "return" expression?

whileStmt       → "while" "(" expression ")" statement

block           → "{" declaration* "}"

expression      → assignment

assignment      → target ( "=" | "-=" | "+=" | "*=" | /=" )
                    assignment
                | ternary

target          → ( postfix "." )? IDENTIFIER
                | postfix "[" expression "]"

ternary         → logic_or ( "?" logic_or ":" logic_or )*

logic_or        → logic_and ( "or" logic_and )*

logic_and       → equality ( "and" equality )*

equality        → comparison ( ( "!=" | "==" ) comparison )*

comparison      → term ( ( ">" | ">=" | "<" | "<=" ) term )*

term            → factor ( ( "-" | "+" ) factor )*

factor          → prefix ( ( "/" | "*" ) prefix )*

prefix          → ( "!" | "-" ) prefix | increment

increment       → target ( "++" | "--" ) | postfix

postfix         → primary ( "." IDENTIFIER
                    | "(" arguments? ")" | "[" index "]" )*

primary         → "true" | "false" | "nil"
                | NUMBER | STRING
                | "(" expression ")"
                | "[" arguments? "]"
                | IDENTIFIER

arguments       → expression ( "," expression )*

index           → expression? ":" expression?
                | expression
```

## Terminals

Excluding the literal text values, this grammar includes five terminal symbols.

- `EOF` is the end-of-file token, which is added by the scanner after the entire source file is read.

- `IDENTIFIER` is a symbol name literal, which must start with an alphabetical character (or an underscore) and be followed by any number of alphanumeric characters. This token may not be one of the reserved words.

- `NUMBER` is a number literal, which must be in decimal notation and may be optionally preceded by '+' or '-'. If a decimal point is included, it must be preceded and followed by valid digits.

- `STRING` is a literal that represents a sequence of characters, set off by double quotation marks. This token must terminate on the same line it is started.

- `TERMINATOR` marks the termination of certain statements. This terminal is unique, insofar as it does not always need to be matched by an actual token. The parser will match this terminal in any of the following cases, in decreasing order of precedence:
  1. When the next token occurs after a newline
  2. When the next token is a semicolon (which gets consumed)
  3. When the next token is a right brace, indicating the close of a block
  4. When the next token is the end-of-file token

# Examples

In addition to the various snippits of code above, three example files have been created to demontrate the functionality of WIXME in detail. Each of these three files also has a corresponding results text file, which records the output produced by running them in the interpreter.

## test.wxm

The first file, *test.wxm*, is essentially a test file that performs unit tests on the language. While by no means extensive, it goes through all of the basic functionality and a few edge cases. The expected output is section headers, each followed by rows of `true`; if any row contains something other than `true`, then a test has failed. This file is quite similar to *test.lox* from my second project, but it includes additional cases for the features unique to WIXME.

## interview.wxm

The second file, *interview.wxm*, implements a few solutions to the same problem posed to me in my interview. The problem statement went something like this: *Suppose you are given a string of digits, like "215634". Each digit or pair of digits can be converted into a alphabetical character, where 1→a, 2→b, ..., 26→z. For example, the string "11" can be converted into "aa" or "k". Write a function that counts the number of unique ways to convert the string of digits into a string of characters.* A few solutions of varying quality are given, and they are compared using a class. See the file for more details.

## coins.wxm

The third file, *coins.wxm*, contains a function that implements a dynamic solution to the coin change problem. The problem is this: *Given an unlimited supply of coins of given denominations, find the minimum number of coins required to get a desired value.* The only expected output of the file is rows of `true`; if any row contains `false`, then the function calculated the wrong solution for a certain input.

# Compilation Instructions

A fully functional interpreter is provided for WIXME, implementing all of the features described above. To compile and run the interpreter, you can use the various commands in the makefile:

- `make build` will compile all the *.go* files into a single executable file.
- `make run` will run the generated executable.
- `make` will perform the actions of both `make build` and `make run`.
- `make clean` will delete the generated executable.
- `make brc` will perform the actions of `make build`, `make run`, and `make clean`.

By default, the command `make run` will run WIXME in interactive mode, where code can be inputted directly through the command line. To specify a target file for the interpreter, set the environment variable `FILE`. For example,

```
make run FILE=test.wxm
```

I don't like having to remember terminal commands, so the makefile was the next best option.
