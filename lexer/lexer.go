package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/nomad-software/script/token"
)

const (
	// EOF represents the end of the input/file.
	EOF = '\uFFFF'
)

// New creates a new instance of the lexer channel.
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		Tokens: make(chan token.Token),
	}
	go l.run()
	return l
}

// Lexer is the instance of the lexer.
type Lexer struct {
	input  string           // The string being scanned.
	start  int              // Start position of this item.
	pos    int              // Current position in the input.
	width  int              // Width of the last rune read.
	Tokens chan token.Token // Channel for lexed tokens
}

type stateFn func(*Lexer) stateFn

func (l *Lexer) run() {
	for state := lex; state != nil; {
		state = state(l)
	}
	close(l.Tokens)
}

func (l *Lexer) read() string {
	return l.input[l.start:l.pos]
}

func (l *Lexer) unread() string {
	return l.input[l.pos:]
}

func (l *Lexer) emit(t token.Type) {
	l.Tokens <- token.Token{
		Type:    t,
		Literal: l.read(),
	}
	l.start = l.pos
}

func (l *Lexer) peek() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return EOF
	}
	r, _ = utf8.DecodeRuneInString(l.unread())
	return r
}

func (l *Lexer) advance() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return EOF
	}
	r, l.width = utf8.DecodeRuneInString(l.unread())
	l.pos += l.width
	return r
}

func (l *Lexer) acceptWhitespace() {
	for unicode.IsSpace(l.peek()) {
		l.advance()
	}
}

func (l *Lexer) discard() {
	l.start = l.pos
}

func (l *Lexer) error(text string) stateFn {
	l.Tokens <- token.Token{
		Type:    token.ILLEGAL,
		Literal: fmt.Sprintf("Parse error: %s %.50q", text, l.read()),
	}
	return nil
}

func lex(l *Lexer) stateFn {
	for {
		l.acceptWhitespace()
		l.discard()

		r := l.advance()

		switch string(r) {
		case token.ASSIGN:
			return lexAssign
		case token.ASTERISK:
			l.emit(token.ASTERISK)
		case token.BANG:
			return lexBang
		case token.COMMA:
			l.emit(token.COMMA)
		case token.GT:
			l.emit(token.GT)
		case token.LBRACE:
			l.emit(token.LBRACE)
		case token.LPAREN:
			l.emit(token.LPAREN)
		case token.LT:
			l.emit(token.LT)
		case token.MINUS:
			l.emit(token.MINUS)
		case token.PLUS:
			l.emit(token.PLUS)
		case token.RBRACE:
			l.emit(token.RBRACE)
		case token.RPAREN:
			l.emit(token.RPAREN)
		case token.SEMICOLON:
			l.emit(token.SEMICOLON)
		case token.SLASH:
			l.emit(token.SLASH)
		case token.DOUBLE_QUOTE:
			return lexString
		case token.EOF:
			return lexEOF
		default:
			if unicode.IsLetter(r) {
				return lexIdentifier
			} else if unicode.IsDigit(r) {
				return lexNumber
			} else {
				return l.error("Illegal token.")
			}
		}
	}
}

func lexAssign(l *Lexer) stateFn {
	if string(l.peek()) == token.ASSIGN {
		l.advance()
		l.emit(token.EQUAL)
	} else {
		l.emit(token.ASSIGN)
	}
	return lex
}

func lexBang(l *Lexer) stateFn {
	if string(l.peek()) == token.ASSIGN {
		l.advance()
		l.emit(token.NOT_EQUAL)
	} else {
		l.emit(token.BANG)
	}
	return lex
}

func lexIdentifier(l *Lexer) stateFn {
	for unicode.IsLetter(l.peek()) {
		l.advance()
	}
	l.emit(token.LookupType(l.read()))
	return lex
}

func lexNumber(l *Lexer) stateFn {
	for unicode.IsDigit(l.peek()) {
		l.advance()
	}
	l.emit(token.INT)
	return lex
}

func lexString(l *Lexer) stateFn {
	l.discard()
	for string(l.peek()) != token.DOUBLE_QUOTE {
		l.advance()
	}
	l.emit(token.STRING)
	l.advance()
	l.discard()
	return lex
}

func lexEOF(l *Lexer) stateFn {
	l.emit(token.EOF)
	return nil
}
