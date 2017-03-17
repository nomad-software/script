package parser

import (
	"fmt"

	"github.com/nomad-software/script/ast"
	"github.com/nomad-software/script/lexer"
	"github.com/nomad-software/script/token"
)

// Parser is the parser itself.
type Parser struct {
	lexer     *lexer.Lexer
	curToken  token.Token
	nextToken token.Token
	errors    []string
}

// New creates a new parser.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}
	p.advance()
	p.advance()
	return p
}

// Parse the tokens from the lexer.
func (p *Parser) Parse() *ast.Progam {
	prg := &ast.Progam{}
	prg.Statements = []ast.Statement{}

	for !p.curToken.IsType(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			prg.Statements = append(prg.Statements, stmt)
		}
		p.advance()
	}
	return prg
}

// Errors return all parsing errors.
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) advance() {
	p.curToken = p.nextToken
	p.nextToken, _ = <-p.lexer.Tokens
}

func (p *Parser) expect(t token.Type) bool {
	if p.nextToken.IsType(t) {
		p.advance()
		return true
	}
	p.addError(t)
	return false
}

func (p *Parser) addError(t token.Type) {
	msg := fmt.Sprintf("Expected token '%s', got '%s' instead", t, p.nextToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{
		Token: p.curToken,
	}

	if !p.expect(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expect(token.ASSIGN) {
		return nil
	}

	for !p.curToken.IsType(token.SEMICOLON) {
		p.advance()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.curToken,
	}

	p.advance()

	for !p.curToken.IsType(token.SEMICOLON) {
		p.advance()
	}

	return stmt
}
