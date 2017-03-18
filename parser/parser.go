package parser

import (
	"fmt"
	"strconv"

	"github.com/nomad-software/script/ast"
	"github.com/nomad-software/script/lexer"
	"github.com/nomad-software/script/token"
)

const (
	LOWEST      = iota
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

type prefixFn func() ast.Expression
type infixFn func(ast.Expression) ast.Expression

// New creates a new parser.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:     l,
		errors:    []string{},
		prefixFns: make(map[token.Type]prefixFn),
		infixFns:  make(map[token.Type]infixFn),
	}

	p.registerPrefixFn(token.IDENT, p.parseIdentifier)
	p.registerPrefixFn(token.INT, p.parseIntegerLiteral)

	p.advance()
	p.advance()
	return p
}

// Parser is the parser itself.
type Parser struct {
	lexer     *lexer.Lexer
	curToken  token.Token
	nextToken token.Token
	errors    []string
	prefixFns map[token.Type]prefixFn
	infixFns  map[token.Type]infixFn
}

// Parse the tokens from the lexer.
func (p *Parser) Parse() *ast.Program {
	prg := &ast.Program{}
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

func (p *Parser) registerPrefixFn(t token.Type, fn prefixFn) {
	p.prefixFns[t] = fn
}

func (p *Parser) registerinfixFn(t token.Type, fn infixFn) {
	p.infixFns[t] = fn
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
		return p.parseExpressionStatement()
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

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token: p.curToken,
	}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.nextToken.IsType(token.SEMICOLON) {
		p.advance()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	fn := p.prefixFns[p.curToken.Type]

	if fn == nil {
		return nil
	}

	leftExp := fn()

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{
		Token: p.curToken,
	}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}
