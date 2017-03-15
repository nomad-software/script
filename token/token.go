package token

import (
	"fmt"
)

type Type string

const (
	ASSIGN    = "="
	ASTERISK  = "*"
	BANG      = "!"
	COMMA     = ","
	ELSE      = "ELSE"
	EOF       = "\uFFFF"
	EQUAL     = "=="
	FALSE     = "FALSE"
	FUNCTION  = "FUNCTION"
	GT        = ">"
	IDENT     = "IDENT"
	IF        = "IF"
	ILLEGAL   = "ILLEGAL"
	INT       = "INT"
	LBRACE    = "{"
	LET       = "LET"
	LPAREN    = "("
	LT        = "<"
	MINUS     = "-"
	NOT_EQUAL = "!="
	PLUS      = "+"
	RBRACE    = "}"
	RETURN    = "RETURN"
	RPAREN    = ")"
	SEMICOLON = ";"
	SLASH     = "/"
	TRUE      = "TRUE"
)

type Token struct {
	Type    Type
	Literal string
}

func (t Token) String() string {
	return fmt.Sprintf("%-8s = %q", t.Type, t.Literal)
}

var keywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func Lookup(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
