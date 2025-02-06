package parser

import (
	"fmt"
	"strconv"

	"github.com/dxtym/skibidi/ast"
	"github.com/dxtym/skibidi/lexer"
	"github.com/dxtym/skibidi/token"
)

// TODO: check parser for correctness
// saw some errors with if statements

// to give operator precedence using enums
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < >
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x
	CALL        // func(x)
	INDEX       // x[y]
)

// associate types with precedences
var precedences = map[token.TokenType]int{
	token.EQUAL:    EQUALS,
	token.NOTEQUAL: EQUALS,
	token.LESS:     LESSGREATER,
	token.MORE:     LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.MUL:      PRODUCT,
	token.DIV:      PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type (
	prefixFn func() ast.Expression
	infixFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l   *lexer.Lexer
	err []string

	currToken token.Token // current token
	nxtToken  token.Token // next token

	prefixFnMap map[token.TokenType]prefixFn
	infixFnMap  map[token.TokenType]infixFn
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, err: []string{}}

	// register prefix functions to token types
	p.prefixFnMap = make(map[token.TokenType]prefixFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseParen)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.IF, p.parseIfElseExpression)
	p.registerPrefix(token.FUNC, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACE, p.parseMapLiteral)
	p.registerPrefix(token.FOR, p.parseForExpression)

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
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

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

	p.NextToken()
	stmt.Value = p.parseExpression(LOWEST)
	for p.nxtToken.Type == token.SEMICOLON {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currToken}
	p.NextToken()

	stmt.Value = p.parseExpression(LOWEST)
	for p.nxtToken.Type == token.SEMICOLON {
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
	e := fmt.Sprintf("hold this l: %s", t)
	p.err = append(p.err, e)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixFnMap[p.currToken.Type]
	if prefix == nil {
		p.noPrefixFnError(p.currToken.Type)
		return nil
	}
	leftExp := prefix()

	// NOTE:
	// pratt parsing top-down approach
	// check out left and right binding powers
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

	val, err := strconv.Atoi(p.currToken.Literal)
	if err != nil {
		e := fmt.Sprintf("sassy baka: %s", p.currToken.Literal)
		p.err = append(p.err, e)
		return nil
	}

	lit.Value = val
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Literal}
}

// construct unary exp like <operator><exp>
func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}

	p.NextToken()
	exp.Right = p.parseExpression(PREFIX)
	return exp
}

// construct binary exp like <exp><operator><exp>
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
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
	e1 := fmt.Sprintf("slop: %s", t)
	e2 := fmt.Sprintf("kino: %s", p.nxtToken.Type)
	p.err = append(p.err, e1, e2)
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

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currToken, Value: p.currToken.Type == token.TRUE}
}

func (p *Parser) parseParen() ast.Expression {
	p.NextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfElseExpression() ast.Expression {
	exp := &ast.IfElseExpression{Token: p.currToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.NextToken()
	exp.Predicate = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Consequence = p.parseBlockStatement()
	if p.nxtToken.Type == token.ELSE {
		p.NextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		exp.Alternative = p.parseBlockStatement()
	}

	return exp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currToken}
	block.Statements = []ast.Statement{}

	p.NextToken()
	for p.currToken.Type != token.RBRACE && p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.NextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	exp := &ast.FunctionLiteral{Token: p.currToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	exp.Parameters = p.parseFunctionArguments()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	exp.Body = p.parseBlockStatement()

	return exp
}

func (p *Parser) parseFunctionArguments() []*ast.Identifier {
	idents := []*ast.Identifier{}
	if p.nxtToken.Type == token.RPAREN {
		p.NextToken()
		return idents
	}

	p.NextToken()
	ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	idents = append(idents, ident)

	for p.nxtToken.Type == token.COMMA {
		p.NextToken()
		p.NextToken()
		idents = append(idents, &ast.Identifier{
			Token: p.currToken,
			Value: p.currToken.Literal,
		})
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return idents
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.currToken, Function: fn}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	arr := &ast.ArrayLiteral{Token: p.currToken}
	arr.Elements = p.parseExpressionList(token.RBRACKET)
	return arr
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}
	if p.nxtToken.Type == end {
		p.NextToken()
		return nil
	}

	p.NextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.nxtToken.Type == token.COMMA {
		p.NextToken()
		p.NextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}
	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.currToken, Left: left}
	p.NextToken()

	exp.Index = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseMapLiteral() ast.Expression {
	mp := &ast.MapLiteral{Token: p.currToken}
	mp.Pairs = make(map[ast.Expression]ast.Expression)

	for p.nxtToken.Literal != token.RBRACE {
		p.NextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.NextToken()
		val := p.parseExpression(LOWEST)
		mp.Pairs[key] = val

		if p.nxtToken.Literal != token.RBRACE && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return mp
}

func (p *Parser) parseForExpression() ast.Expression {
	fl := &ast.ForExpression{Token: p.currToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.NextToken()
	fl.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	fl.Body = p.parseBlockStatement()
	return fl
}