package parser

import (
	"testing"

	"github.com/nomad-software/script/ast"
	"github.com/nomad-software/script/lexer"
	"github.com/nomad-software/script/token"
)

type test struct {
	name string
}

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foo = 1337;
`
	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.Parse()
	checkParserErrors(t, parser)

	if program == nil {
		t.Fatal("Parse() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Not enough statements, only got %d", len(program.Statements))
	}

	tests := []test{
		{"x"},
		{"y"},
		{"foo"},
	}

	for x, test := range tests {
		stmt := program.Statements[x]
		testLetStatement(t, stmt, test.name)
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 1337;
`
	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.Parse()
	checkParserErrors(t, parser)

	if len(program.Statements) != 3 {
		t.Fatalf("Not enough statements, only got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		cast, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("Wrong stmt type: %T", stmt)
		}

		if cast.TokenLiteral() != token.RETURN {
			t.Errorf("Wrong literal, expecting %q, got %q", token.RETURN, cast.TokenLiteral())
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) {
	if stmt.TokenLiteral() != token.LET {
		t.Errorf("Wrong stmt literal, expecting %q, got %q", token.LET, stmt.TokenLiteral())
	}

	cast, ok := stmt.(*ast.LetStatement)

	if !ok {
		t.Errorf("Wrong stmt type: %T", stmt)
	}

	if cast.Name.Value != name {
		t.Errorf("Wrong name, expecting %q, got %q", name, cast.Name.Value)
	}

	if cast.Name.TokenLiteral() != name {
		t.Errorf("Wrong literal, expecting %q, got %q", name, cast.Name.TokenLiteral())
	}
}
