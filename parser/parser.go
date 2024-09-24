package parser

import (
	"fmt"
	"strconv"

	"github.com/dxtym/monke/ast"
	"github.com/dxtym/monke/lexer"
	"github.com/dxtym/monke/token"
)

// to give operator precedence using enums
const (
	_ int = iota
	LOWEST
	EQUALS // ==
	LESSGREATER // < >
	SUM // +
	PRODUCT // *
	PREFIX // -x 
	CALL // func(x)
)

// associate types with precedences
var precedences = map[token.TokenType]int{
	token.EQUAL: EQUALS,
	token.NOTEQUAL: EQUALS,
	token.LESS: LESSGREATER,
	token.MORE: LESSGREATER,
	token.PLUS: SUM,
	token.MINUS: SUM,
	token.MUL: PRODUCT,
	token.DIV: PRODUCT,
}

type (
	prefixFn func() ast.Expression
	infixFn func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer
	err []string

	currToken token.Token // current token
	nxtToken token.Token // next token

	prefixFnMap map[token.TokenType]prefixFn
	infixFnMap map[token.TokenType]infixFn
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, err: []string{}}

	// register prefix functions to token types
	p.prefixFnMap = make(map[token.TokenType]prefixFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	// register infix functions to token types
	p.infixFnMap = make(map[token.TokenType]infixFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.MUL, p.parseInfixExpression)
	p.registerInfix(token.DIV, p.parseInfixExpression)
	p.registerInfix(token.LESS, p.parseInfixExpression)
	p.registerInfix(token.MORE, p.parseInfixExpression)
	p.registerInfix(token.EQUAL, p.parseInfixExpression)
	p.registerInfix(token.NOTEQUAL, p.parseInfixExpression)
	
	p.NextToken()
	p.NextToken()
	return p
}

// move curr and nxt pointers to tokens
func (p *Parser) NextToken() {
	p.currToken = p.nxtToken
	p.nxtToken = p.l.NextToken()
}

func (p *Parser) Errors() []string {
	return p.err
}

// acts like a stack machine of tokens
func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.NextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// NOTE:
// why ast.Statement accepts ast.LetStatement? because they
// both implement TokenLiteral() method, so ast.Statement
// serves as a general interface for ast.LetStatement
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for p.currToken.Type != token.SEMICOLON {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currToken}
	p.NextToken()
	
	for p.currToken.Type != token.SEMICOLON {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}
	stmt.Expression = p.parseExpression(LOWEST)
	
	if p.nxtToken.Type == token.SEMICOLON {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) noPrefixFnError(t token.TokenType) {
	e := fmt.Sprintf("no prefix function found for %s", t)
	p.err = append(p.err, e)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixFnMap[p.currToken.Type]
	if prefix == nil {
		p.noPrefixFnError(p.currToken.Type)
		return nil
	}
	leftExp := prefix()

	// pratt parsing technique
	for p.nxtToken.Type != token.SEMICOLON && precedence < p.peekPrecedence() {
		infix := p.infixFnMap[p.nxtToken.Type]
		if infix == nil {
			return leftExp
		}

		p.NextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}


func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currToken}
	
	val, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		e := fmt.Sprintf("caanot convert %q to int", p.currToken.Literal)
		p.err = append(p.err, e)
		return nil
	}
	
	lit.Value = val
	return lit
}

// construct unary exp like <operator><exp>
func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token: p.currToken,
		Operator: p.currToken.Literal,
	}

	p.NextToken()
	exp.Right = p.parseExpression(PREFIX)
	return exp
}

// construct binary exp like <exp><operator><exp>
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token: p.currToken,
		Operator: p.currToken.Literal,
		Left: left,
	}

	precedence := p.currPrecedence()
	p.NextToken()
	exp.Right = p.parseExpression(precedence)

	return exp
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.nxtToken.Type == t {
		p.NextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	e := fmt.Sprintf("next token got=%s, expected=%s", t, p.currToken.Type)
	p.err = append(p.err, e)
}

func (p *Parser) registerPrefix(tt token.TokenType, fn prefixFn) {
	p.prefixFnMap[tt] = fn
}

func (p *Parser) registerInfix(tt token.TokenType, fn infixFn) {
	p.infixFnMap[tt] = fn
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.nxtToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}
	return LOWEST
}