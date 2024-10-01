package parser

import (
	"fmt"
	"testing"

	"github.com/dxtym/maymun/ast"
	"github.com/dxtym/maymun/lexer"
	"github.com/dxtym/maymun/token"
)

// TODO: standardize errors

func checkParser(t *testing.T, p *Parser) {
	err := p.Errors()
	if len(err) == 0 {
		return 
	}

	t.Errorf("parser has %d errors", len(err))
	for _, e := range err {
		t.Errorf("parser error: %q", e)
	}
	t.FailNow()
}

func testIdentifier(t *testing.T, exp ast.Expression, val string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != val {
		t.Errorf("ident.Value not %s. got=%s", val, ident.Value)
		return false
	}
	if ident.TokenLiteral() != val {
		t.Errorf("ident.TokenLiteral not %s. got=%s", val, ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, val bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != val {
		t.Errorf("bo.Value not equal to %t. got=%t", val, bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", val) {
		t.Errorf("bo.TokenLiteral() not equal to %t. got=%s", val, bo.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	default:
		t.Errorf("type of exp not handled. got=%T", exp)
		return false
	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, op string, left, right any) bool {
	ope, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T", exp)
		return false
	}

	if !testLiteralExpression(t, ope.Left, left) {
		return false
	}
	if ope.Operator != op {
		t.Errorf("ope.Operator not equal to %s. got=%s", op, ope.Operator)
		return false
	}
	if !testLiteralExpression(t, ope.Right, right) {
		return false
	}

	return true
}

// TODO: move input to text file
func TestLetStatements(t *testing.T) {
	input := "let x = 1; let y = 2;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)
	if program == nil {
		t.Fatal("Parse() returned nil")
	}
	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements doesn't contain 2 statements. got=%d", len(program.Statements)) 
	}

	tests := []struct{
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, identifier string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Fatalf("stmt.TokenLiteral() not 'let'. got=%q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Fatalf("cannot convert to *ast.LetStatement. got=%T", stmt)
		return false
	}

	if letStmt.Name.Value != identifier {
		t.Fatalf("letStmt.Name.Value not equal to %q. got=%s", identifier, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != identifier {
		t.Fatalf("letStmt.TokenLiteral() not equal to %q. got=%s", identifier, letStmt.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatement(t *testing.T) {
	input := "return 1; return add(10);"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)
	if program == nil {
		t.Fatal("Parse() returned nil")
	}
	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements doesn't contain 2 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("cannot convert to *ast.LetStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not equal to 'return'. got=%s", returnStmt.TokenLiteral())
		}
	}
}

func TestString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "a"},
					Value: "a",
				},
				Value: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "b"},
					Value: "b",
				},
			},
		},
	}

	if program.String() != "let a = b;" {
		t.Fatalf("program.String() wrong. got=%q", program.String())
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has not enough arguments. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression not *ast.Identifier. got=%T", stmt.Expression)
	}
	
	if ident.Value != "foobar" {
		t.Fatalf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteral(t *testing.T) {
	input := "5;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has not enough arguments. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Fatalf("literal.Value not equal to 5. got=%d", literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Fatalf("literal.TokenLiteral() not equal to '5'. got=%q", literal.TokenLiteral())
	}
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct{
		in string
		op string
		val int64
	}{
		{"-1;", "-", 1},
		{"!2;", "!", 2},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.in)
		p := NewParser(l)
		program := p.Parse()
		checkParser(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements has not enough arguments. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression not *ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.op {
			t.Fatalf("exp.Operator not equal to %s. got=%s", tt.op, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.val) {
			return
		}
	}
} 

func TestInfixExpression(t *testing.T) {
	tests := []struct{
		in string
		left int64
		op string
		right int64
	}{
		{"1 + 1;", 1, "+", 1},
		{"1 - 1;", 1, "-", 1},
		{"1 * 1;", 1, "*", 1},
		{"1 / 1;", 1, "/", 1},
		{"1 > 1;", 1, ">", 1},
		{"1 < 1;", 1, "<", 1},
		{"1 == 1;", 1, "==", 1},
		{"1 != 1;", 1, "!=", 1},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.in)
		p := NewParser(l)
		program := p.Parse()
		checkParser(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements has not enough arguments. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt.Expression not *ast.InfixExpression. got=%T", stmt.Expression)
		}

		if !testLiteralExpression(t, exp.Left, tt.left) {
			return
		}
		if exp.Operator != tt.op {
			t.Fatalf("exp.Operator not equal to %s. got=%s", tt.op, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.left) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, val int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("literal not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integer.Value != val {
		t.Errorf("integer.Value not equal to %d. got=%d", val, integer.Value)
		return false
	}
	if integer.TokenLiteral() != fmt.Sprintf("%d", val) {
		t.Errorf("integer.TokenLiteral not equal to %d. got=%s", val, integer.TokenLiteral())
		return false
	}

	return true
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct{
		in string
		out string
	}{
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"1 + 2 * 3 + 4 / 5 - 6",
			"(((1 + (2 * 3)) + (4 / 5)) - 6)",
		},
		{
			"2 > 1 == 3 < 4",
			"((2 > 1) == (3 < 4))",
		},
		{
			"3 + 4 * -5 != 3 * -1 + 4 * 5",
			"((3 + (4 * (-5))) != ((3 * (-1)) + (4 * 5)))",
		},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.in)
		p := NewParser(l)
		program := p.Parse()
		checkParser(t, p)

		actual := program.String()
		if actual != tt.out {
			t.Fatalf("expected=%q, got=%q", tt.out, actual)
		}
	}
}

func TestBoolean(t *testing.T) {
	tests := []struct{
		in string
		out bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.in)
		p := NewParser(l)
		program := p.Parse()
		checkParser(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements has not enough arguments. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		bo, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("stmt.Expression not *ast.Boolran. got=%T", stmt.Expression)
		}

		if bo.Value != tt.out {
			t.Fatalf("bo.Value not %t. got=%t", tt.out, bo.Value)
		}
		if bo.TokenLiteral() != tt.in[:len(tt.in)-1] {
			t.Fatalf("bo.TokenLiteral not %s. got=%s", tt.in[:len(tt.in)-1], bo.TokenLiteral())
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := "if (x < y) { x };"
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has not enough arguments. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IfElseExpression)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.Boolran. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Predicate, "<", "x", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}
	cons, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, cons.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y };"
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has not enough arguments. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IfElseExpression)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.Boolran. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Predicate, "<", "x", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("exp.Consequence.Statements has not enough arguments. got=%d\n", len(exp.Consequence.Statements))
	}
	cons, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Consequence.Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, cons.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Fatalf("exp.Alternative.Statements has not enough arguments. got=%d\n", len(exp.Alternative.Statements))
	}
	alt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("exp.Alternative.Statements[0] is not ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alt.Expression, "y") {
		return
	}
}

func TestFunctionLiteral(t *testing.T) {
	input := "func(x, y) {x + y};"
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has not enough arguments. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("smt.Expression not *ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(fn.Arguments) != 2 {
		t.Fatalf("fn.Arguments has not enough arguments. got=%d", len(fn.Arguments))
	}

	testLiteralExpression(t, fn.Arguments[0], "x")
	testLiteralExpression(t, fn.Arguments[1], "y") 

	if len(fn.Body.Statements) != 1 {
		t.Fatalf("fn.Body.Statements has not enough arguments. got=%d", len(fn.Body.Statements))
	}

	body, ok := fn.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("fn.Body.Expression[0] not *ast.ExperssionStatement. got=%T", fn.Body.Statements[0])
	}

	testInfixExpression(t, body.Expression, "+", "x", "y")
}

func TestFunctionArgumentParsing(t *testing.T) {
	tests := []struct{
		in string
		out []string
	}{
		{"func() {}", []string{}},
		{"func(x) {}", []string{"x"}},
		{"func(x, y, z) {}", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.in)
		p := NewParser(l)
		program := p.Parse()
		checkParser(t, p)

		stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
		fn, _ := stmt.Expression.(*ast.FunctionLiteral)

		if len(fn.Arguments) != len(tt.out) {
			t.Fatalf("fn.Arguments has not enough arguments. got=%d", len(fn.Arguments))
		}

		for i := range fn.Arguments {
			testLiteralExpression(t, fn.Arguments[i], tt.out[i])
		}
	}
}