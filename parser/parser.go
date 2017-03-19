package parser

import (
	"fmt"
	"strconv"

	"github.com/nomad-software/script/ast"
	"github.com/nomad-software/script/lexer"
	"github.com/nomad-software/script/precedence"
	"github.com/nomad-software/script/token"
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
	p.registerPrefixFn(token.BANG, p.parsePrefixExpression)
	p.registerPrefixFn(token.MINUS, p.parsePrefixExpression)

	p.registerInfixFn(token.PLUS, p.parseInfixExpression)
	p.registerInfixFn(token.MINUS, p.parseInfixExpression)
	p.registerInfixFn(token.SLASH, p.parseInfixExpression)
	p.registerInfixFn(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixFn(token.EQUAL, p.parseInfixExpression)
	p.registerInfixFn(token.NOT_EQUAL, p.parseInfixExpression)
	p.registerInfixFn(token.LT, p.parseInfixExpression)
	p.registerInfixFn(token.GT, p.parseInfixExpression)

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

func (p *Parser) addError(format string, a ...interface{}) {
	p.errors = append(p.errors, fmt.Sprintf(format, a...))
}

// Errors return all parsing errors.
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerPrefixFn(t token.Type, fn prefixFn) {
	p.prefixFns[t] = fn
}

func (p *Parser) registerInfixFn(t token.Type, fn infixFn) {
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
	p.addError("Expected token '%s', got '%s' instead", t, p.nextToken.Type)
	return false
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

	stmt.Expression = p.parseExpression(precedence.LOWEST)

	if p.nextToken.IsType(token.SEMICOLON) {
		p.advance()
	}

	return stmt
}

func (p *Parser) parseExpression(prec int) ast.Expression {
	prefix := p.prefixFns[p.curToken.Type]
	if prefix == nil {
		p.addError("no prefix parse function for %s found", p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.nextToken.IsType(token.SEMICOLON) && prec < p.nextToken.Precedence() {
		infix := p.infixFns[p.nextToken.Type]
		if infix == nil {
			return leftExp
		}

		p.advance()

		leftExp = infix(leftExp)
	}

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
		p.addError("could not parse %q as integer", p.curToken.Literal)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.advance()

	expression.Right = p.parseExpression(precedence.PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	prec := p.curToken.Precedence()
	p.advance()
	expression.Right = p.parseExpression(prec)

	return expression
}
