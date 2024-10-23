package parser

import (
	"fmt"
	"testing"

	"github.com/dxtym/maymun/ast"
	"github.com/dxtym/maymun/lexer"
	"github.com/dxtym/maymun/token"
)

func checkParser(t *testing.T, p *Parser) {
	err := p.Errors()
	if len(err) == 0 {
		return
	}

	t.Errorf("parser %d errors", len(err))
	for _, e := range err {
		t.Errorf("parser error: %q", e)
	}
	t.FailNow()
}

func testIdentifier(t *testing.T, exp ast.Expression, val string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier: got=%T", exp)
		return false
	}

	if ident.Value != val {
		t.Errorf("ident.Value not equal to %s: got=%s", val, ident.Value)
		return false
	}
	if ident.TokenLiteral() != val {
		t.Errorf("ident.TokenLiteral not equal to %s: got=%s", val, ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, val bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean: got=%T", exp)
		return false
	}

	if bo.Value != val {
		t.Errorf("bo.Value not equal to %t: got=%t", val, bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", val) {
		t.Errorf("bo.TokenLiteral not equal to %t: got=%s", val, bo.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	default:
		t.Errorf("unexpected type: got=%T", exp)
		return false
	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, op string, left, right any) bool {
	ope, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp not *ast.OperatorExpression: got=%T", exp)
		return false
	}

	if !testLiteralExpression(t, ope.Left, left) {
		return false
	}
	if ope.Operator != op {
		t.Errorf("ope.Operator not equal to %s: got=%s", op, ope.Operator)
		return false
	}
	if !testLiteralExpression(t, ope.Right, right) {
		return false
	}

	return true
}

func TestProgram(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "deylik"},
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

	if program.String() != "deylik a = b;" {
		t.Errorf("program.String not equal to %s: got=%s", "let a = b;", program.String())
	}
}

// TODO: move input to text file
func TestLetStatements(t *testing.T) {
	input := "deylik x = 1; deylik y = y;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)

	if len(program.Statements) != 2 {
		t.Errorf("program.Statements must be 2 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		want string
	}{
		{"x"},
		{"y"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.want) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, identifier string) bool {
	if stmt.TokenLiteral() != "deylik" {
		t.Errorf("stmt.TokenLiteral not %s: got=%s", "let", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt not *ast.LetStatement: got=%T", stmt)
		return false
	}

	if letStmt.Name.Value != identifier {
		t.Errorf("letStmt.Name.Value not equal to %s: got=%s", identifier, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != identifier {
		t.Errorf("letStmt.TokenLiteral not equal to %s: got=%s", identifier, letStmt.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatement(t *testing.T) {
	input := "qaytar 1; qaytar ortirish(10);"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)

	if len(program.Statements) != 2 {
		t.Errorf("program.Statements must be 2 statements: got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		rtrn, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.LetStatement: got=%T", stmt)
			continue
		}
		if rtrn.TokenLiteral() != "qaytar" {
			t.Errorf("rtrn.TokenLiteral not equal to %s: got=%s", "return", rtrn.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements must be 1 statement: got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] not *ast.ExpressionStatement: got=%T", program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Errorf("stmt.Expression not *ast.Identifier: got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not equal to %s: got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not equal to %s: got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteral(t *testing.T) {
	input := "5;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements must be 1 statement: got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] not *ast.ExpressionStatement: got=%T", program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("stmt.Expression not *ast.IntegerLiteral: got=%T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not equal to %d: got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not equal to %s: got=%s", "5", literal.TokenLiteral())
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"hello, world";`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements must be 1 statement: got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] not *ast.ExpressionStatement: got=%T", program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Errorf("stmt.Expression not *ast.StringLiteral: got=%T", stmt.Expression)
	}

	if literal.Value != "hello, world" {
		t.Errorf("literal.Value not equal to %s: got=%s", "hello, world", literal.Value)
	}
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		got  string
		op   string
		want int
	}{
		{"-1;", "-", 1},
		{"!2;", "!", 2},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.got)
		p := NewParser(l)
		program := p.Parse()
		checkParser(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("program.Statements must be 1 statement: got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("program.Statements[0] not *ast.ExpressionStatement: got=%T", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Errorf("stmt.Expression not *ast.PrefixExpression: got=%T", stmt.Expression)
		}

		if exp.Operator != tt.op {
			t.Errorf("exp.Operator not equal to %s: got=%s", tt.op, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.want) {
			return
		}
	}
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		got   string
		left  int
		op    string
		right int
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
		l := lexer.NewLexer(tt.got)
		p := NewParser(l)
		program := p.Parse()
		checkParser(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("program.Statements must be 1 statement: got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("program.Statements[0] not *ast.ExpressionStatement: got=%T", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Errorf("stmt.Expression not *ast.InfixExpression: got=%T", stmt.Expression)
		}

		if !testLiteralExpression(t, exp.Left, tt.left) {
			return
		}
		if exp.Operator != tt.op {
			t.Errorf("exp.Operator not equal to %s: got=%s", tt.op, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.left) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, val int) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("literal not *ast.IntegerLiteral: got=%T", il)
		return false
	}
	if integer.Value != val {
		t.Errorf("integer.Value not equal to %d: got=%d", val, integer.Value)
		return false
	}
	if integer.TokenLiteral() != fmt.Sprintf("%d", val) {
		t.Errorf("integer.TokenLiteral not equal to %d: got=%s", val, integer.TokenLiteral())
		return false
	}

	return true
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		got  string
		want string
	}{
		{
			"3 > 5 == salbiy",
			"((3 > 5) == salbiy)",
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
		l := lexer.NewLexer(tt.got)
		p := NewParser(l)
		program := p.Parse()
		checkParser(t, p)

		actual := program.String()
		if actual != tt.want {
			t.Errorf("program.String not equal to %s: got=%s", tt.want, actual)
		}
	}
}

func TestBoolean(t *testing.T) {
	tests := []struct {
		got  string
		want bool
	}{
		{"ijobiy;", true},
		{"salbiy;", false},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.got)
		p := NewParser(l)
		program := p.Parse()
		checkParser(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("program.Statements must be 1 statement: got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("program.Statements[0] not *ast.ExpressionStatement: got=%T", program.Statements[0])
		}
		bo, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Errorf("stmt.Expression not *ast.Boolean: got=%T", stmt.Expression)
		}

		if bo.Value != tt.want {
			t.Errorf("bo.Value not equal to %t: got=%t", tt.want, bo.Value)
		}
		if bo.TokenLiteral() != tt.got[:len(tt.got)-1] {
			t.Errorf("bo.TokenLiteral not equal to %s: got=%s", tt.got[:len(tt.got)-1], bo.TokenLiteral())
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := "agar (x < y) { x };"
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements must be 1 statement: got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] not *ast.ExpressionStatement: got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IfElseExpression)
	if !ok {
		t.Errorf("stmt.Expression not *ast.Boolean: got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Predicate, "<", "x", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("exp.Consequence.Statements must be 1 statement: got=%d\n", len(exp.Consequence.Statements))
	}
	cons, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Statements[0] not *ast.ExpressionStatement: got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, cons.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements not equal to nil. got=%v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "agar (x < y) { x } yana { y };"
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements must be 1 statement: got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] not *ast.ExpressionStatement: got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IfElseExpression)
	if !ok {
		t.Errorf("stmt.Expression not *ast.Boolean: got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Predicate, "<", "x", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("exp.Consequence.Statements must be 1 statement: got=%d\n", len(exp.Consequence.Statements))
	}
	cons, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("exp.Consequence.Statements[0] not *ast.ExpressionStatement: got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, cons.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements must be 1 statement: got=%d\n", len(exp.Alternative.Statements))
	}
	alt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("exp.Alternative.Statements[0] not *ast.ExpressionStatement: got=%T", exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alt.Expression, "y") {
		return
	}
}

func TestFunctionLiteral(t *testing.T) {
	input := "amal(x, y) {x + y};"
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements must be 1 statement: got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] not *ast.ExpressionStatement: got=%T", program.Statements[0])
	}
	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Errorf("smt.Expression not *ast.FunctionLiteral: got=%T", stmt.Expression)
	}

	if len(fn.Parameters) != 2 {
		t.Errorf("fn.Parameters must be 2 statements: got=%d", len(fn.Parameters))
	}

	testLiteralExpression(t, fn.Parameters[0], "x")
	testLiteralExpression(t, fn.Parameters[1], "y")

	if len(fn.Body.Statements) != 1 {
		t.Errorf("fn.Body.Statements must be 1 statement: got=%d", len(fn.Body.Statements))
	}

	body, ok := fn.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("fn.Body.Expression[0] not *ast.ExperssionStatement: got=%T", fn.Body.Statements[0])
	}

	testInfixExpression(t, body.Expression, "+", "x", "y")
}

func TestFunctionArgumentParsing(t *testing.T) {
	tests := []struct {
		got  string
		want []string
	}{
		{"amal() {}", []string{}},
		{"amal(x) {}", []string{"x"}},
		{"amal(x, y, z) {}", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.got)
		p := NewParser(l)
		program := p.Parse()
		checkParser(t, p)

		stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
		fn, _ := stmt.Expression.(*ast.FunctionLiteral)

		if len(fn.Parameters) != len(tt.want) {
			t.Errorf("fn.Parameters not equal to %d: got=%d", len(tt.want), len(fn.Parameters))
		}

		for i := range fn.Parameters {
			testLiteralExpression(t, fn.Parameters[i], tt.want[i])
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := "bajarish(1 + 2, 3 / 1);"
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements must be 1 statement: got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] not *ast.ExpressionStatement: got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Errorf("smt.Expression not *ast.CallExpression: got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "bajarish") {
		t.Errorf("exp.Function not eqaul to %s: got=%s", "add", exp.Function)
	}

	if len(exp.Arguments) != 2 {
		t.Errorf("exp.Arguments must be 2 statements: got=%d", len(exp.Arguments))
	}
	testInfixExpression(t, exp.Arguments[0], "+", 1, 2)
	testInfixExpression(t, exp.Arguments[1], "/", 3, 1)
}
