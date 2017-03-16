package token

import "fmt"

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

// Other
const (
	EOF       = "\uFFFF"
	IDENT     = "identifier"
	ILLEGAL   = "illegal"
	INT       = "int"
	SEMICOLON = ";"
)

type Token struct {
	Type    Type
	Literal string
}

// IsType checks if this token is of a particular type.
func (t Token) IsType(other Type) bool {
	return t.Type == other
}

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

func LookupType(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
