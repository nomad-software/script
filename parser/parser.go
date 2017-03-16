package parser

import (
	"github.com/nomad-software/script/ast"
	"github.com/nomad-software/script/lexer"
	"github.com/nomad-software/script/token"
)

// Parser is the parser itself.
type Parser struct {
	lexer     *lexer.Lexer
	curToken  token.Token
	nextToken token.Token
}

// New creates a new parser.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer: l,
	}
	p.advance()
	p.advance()
	return p
}

// Parse the tokens from the lexer.
func (p *Parser) Parse() *ast.Progam {
	prg := &ast.Progam{}
	prg.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			prg.Statements = append(prg.Statements, stmt)
		}
		p.advance()
	}
	return prg
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
	return false
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() ast.Statement {
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
