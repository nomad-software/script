package token

import (
	"fmt"

	"github.com/nomad-software/script/precedence"
)

// Type represents the type of a token.
type Type string

// Keywords
const (
	ELSE     = "else"
	FALSE    = "false"
	FUNCTION = "fn"
	IF       = "if"
	LET      = "let"
	RETURN   = "return"
	TRUE     = "true"
)

// Operators
const (
	ASSIGN    = "="
	ASTERISK  = "*"
	BANG      = "!"
	COMMA     = ","
	EQUAL     = "=="
	GT        = ">"
	LBRACE    = "{"
	LPAREN    = "("
	LT        = "<"
	MINUS     = "-"
	NOT_EQUAL = "!="
	PLUS      = "+"
	RBRACE    = "}"
	RPAREN    = ")"
	SLASH     = "/"
)

// Data types
const (
	INT    = "int"
	STRING = "string"
)

// Miscellaneous
const (
	DOUBLE_QUOTE = "\""
	EOF          = "\uFFFF"
	IDENT        = "identifier"
	ILLEGAL      = "illegal"
	SEMICOLON    = ";"
)

// Token represents a unit of output from the lexer.
type Token struct {
	Type    Type
	Literal string
}

// IsType checks if this token is of a particular type.
func (t Token) IsType(other Type) bool {
	return t.Type == other
}

// String returns a string representation of a token.
func (t Token) String() string {
	return fmt.Sprintf("%-10s = %q", t.Type, t.Literal)
}

var keywords = map[string]Type{
	FUNCTION: FUNCTION,
	LET:      LET,
	TRUE:     TRUE,
	FALSE:    FALSE,
	IF:       IF,
	ELSE:     ELSE,
	RETURN:   RETURN,
}

// LookupType returns the token type for the passed identifier.
func LookupType(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

var precedences = map[Type]int{
	EQUAL:     precedence.EQUALS,
	NOT_EQUAL: precedence.EQUALS,
	GT:        precedence.LESSGREATER,
	LT:        precedence.LESSGREATER,
	MINUS:     precedence.SUM,
	PLUS:      precedence.SUM,
	ASTERISK:  precedence.PRODUCT,
	SLASH:     precedence.PRODUCT,
	LPAREN:    precedence.CALL,
}

// Precedence returns the precedence of a token's type.
func (t Token) Precedence() int {
	if p, ok := precedences[t.Type]; ok {
		return p
	}
	return precedence.LOWEST
}
