package lexer

import (
	"testing"

	"github.com/nomad-software/script/token"
)

type test struct {
	typ     token.Type
	literal string
}

func TestLexingSingleCharacters(t *testing.T) {
	input := `=*!,>{(<-+});/`

	tests := []test{
		{token.ASSIGN, "="},
		{token.ASTERISK, "*"},
		{token.BANG, "!"},
		{token.COMMA, ","},
		{token.GT, ">"},
		{token.LBRACE, "{"},
		{token.LPAREN, "("},
		{token.LT, "<"},
		{token.MINUS, "-"},
		{token.PLUS, "+"},
		{token.RBRACE, "}"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.SLASH, "/"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for i, test := range tests {
		tok := <-lexer.Tokens

		if tok.Type != test.typ {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, test.typ, tok.Type)
		}

		if tok.Literal != test.literal {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, test.literal, tok.Literal)
		}
	}
}

func TestLexingVariables(t *testing.T) {
	input := `let five = 5;`

	tests := []test{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for i, test := range tests {
		tok := <-lexer.Tokens

		if tok.Type != test.typ {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, test.typ, tok.Type)
		}

		if tok.Literal != test.literal {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, test.literal, tok.Literal)
		}
	}
}

func TestLexingEdgeCases(t *testing.T) {
	input := `
!-/*5;
5 < 10 > 5;

10 == 10;
10 != 9;`

	tests := []test{
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.INT, "10"},
		{token.EQUAL, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQUAL, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},

		{token.EOF, ""},
	}

	lexer := New(input)

	for i, test := range tests {
		tok := <-lexer.Tokens

		if tok.Type != test.typ {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, test.typ, tok.Type)
		}

		if tok.Literal != test.literal {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, test.literal, tok.Literal)
		}
	}
}

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
"foobar"
"foo bar"
`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQUAL, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQUAL, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		// {token.LBRACKET, "["},
		// {token.INT, "1"},
		// {token.COMMA, ","},
		// {token.INT, "2"},
		// {token.RBRACKET, "]"},
		// {token.SEMICOLON, ";"},
		// {token.LBRACE, "{"},
		// {token.STRING, "foo"},
		// {token.COLON, ":"},
		// {token.STRING, "bar"},
		// {token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := <-l.Tokens

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
