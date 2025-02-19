package parser

import (
	"fmt"
	"testing"

	"github.com/dxtym/skibidi/ast"
	"github.com/dxtym/skibidi/lexer"
	"github.com/dxtym/skibidi/token"
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
	b, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean: got=%T", exp)
		return false
	}

	if b.Value != val {
		t.Errorf("b.Value not equal to %t: got=%t", val, b.Value)
		return false
	}
	if b.TokenLiteral() != fmt.Sprintf("%t", val) {
		t.Errorf("b.TokenLiteral not equal to %t: got=%s", val, b.TokenLiteral())
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
				Token: token.Token{Type: token.LET, Literal: "amogus"},
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

	if program.String() != "amogus a = b;" {
		t.Errorf("program.String not equal to %s: got=%s", "amogus a = b;", program.String())
	}
}

// TODO: move got to text file
func TestLetStatements(t *testing.T) {
	got := "amogus x = 1; amogus y = y;"

	l := lexer.NewLexer(got)
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
	if stmt.TokenLiteral() != "amogus" {
		t.Errorf("stmt.TokenLiteral not %s: got=%s", "amogus", stmt.TokenLiteral())
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
	got := "rizz 1; rizz add(10);"

	l := lexer.NewLexer(got)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)

	if len(program.Statements) != 2 {
		t.Errorf("program.Statements must be 2 statements: got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		r, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.LetStatement: got=%T", stmt)
			continue
		}
		if r.TokenLiteral() != "rizz" {
			t.Errorf("r.TokenLiteral not equal to %s: got=%s", "return", r.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	got := "foobar;"

	l := lexer.NewLexer(got)
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
	got := "5;"

	l := lexer.NewLexer(got)
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
	got := `"hello, world";`

	l := lexer.NewLexer(got)
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
			"3 > 5 == cap",
			"((3 > 5) == cap)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"2 > 1 == 3 < 4",
			"((2 > 1) == (3 < 4))",
		},
		{
			"1 + 2 * 3 + 4 / 5 - 6",
			"(((1 + (2 * 3)) + (4 / 5)) - 6)",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
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
		{"fax;", true},
		{"cap;", false},
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
		b, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Errorf("stmt.Expression not *ast.Boolean: got=%T", stmt.Expression)
		}

		if b.Value != tt.want {
			t.Errorf("b.Value not equal to %t: got=%t", tt.want, b.Value)
		}
		if b.TokenLiteral() != tt.got[:len(tt.got)-1] {
			t.Errorf("b.TokenLiteral not equal to %s: got=%s", tt.got[:len(tt.got)-1], b.TokenLiteral())
		}
	}
}

func TestIfExpression(t *testing.T) {
	got := "hawk (x < y) { x };"
	l := lexer.NewLexer(got)
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

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements not equal to nil. got=%v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	got := "hawk (x < y) { x } tuah { y };"
	l := lexer.NewLexer(got)
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
	got := "cook(x, y) {x + y};"
	l := lexer.NewLexer(got)
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
		t.Errorf("stmt.Expression not *ast.FunctionLiteral: got=%T", stmt.Expression)
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
		{"cook() {}", []string{}},
		{"cook(x) {}", []string{"x"}},
		{"cook(x, y, z) {}", []string{"x", "y", "z"}},
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
	got := "do(1 + 2, 3 / 1);"
	l := lexer.NewLexer(got)
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
		t.Errorf("stmt.Expression not *ast.CallExpression: got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "do") {
		t.Errorf("exp.Function not eqaul to %s: got=%s", "add", exp.Function)
	}

	if len(exp.Arguments) != 2 {
		t.Errorf("exp.Arguments must be 2 statements: got=%d", len(exp.Arguments))
	}
	testInfixExpression(t, exp.Arguments[0], "+", 1, 2)
	testInfixExpression(t, exp.Arguments[1], "/", 3, 1)
}

func TestArrayLiteral(t *testing.T) {
	got := "[1, 2 * 2, 3 + 3];"
	l := lexer.NewLexer(got)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] not *ast.ExpressionStatement: got=%T", program.Statements[0])
	}
	arr, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Errorf("stmt.Expression not *ast.IndexExpression: got=%T", stmt.Expression)
	}

	if len(arr.Elements) != 3 {
		t.Errorf("arr.Elements must be 2 statements: got=%d", len(arr.Elements))
	}

	testIntegerLiteral(t, arr.Elements[0], 1)
	testInfixExpression(t, arr.Elements[1], "*", 2, 2)
	testInfixExpression(t, arr.Elements[2], "+", 3, 3)
}

func TestIndexExpression(t *testing.T) {
	got := "arr[1 + 1];"
	l := lexer.NewLexer(got)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] not *ast.ExpressionStatement: got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Errorf("stmt.Expression not *ast.IndexExpression: got=%T", stmt.Expression)
	}

	testIdentifier(t, exp.Left, "arr")
	testInfixExpression(t, exp.Index, "+", 1, 1)
}

func TestMapLiteral(t *testing.T) {
	got := `{"foo": "bar", "baz": "qux"};`
	l := lexer.NewLexer(got)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] not *ast.ExpressionStatement: got=%T", program.Statements[0])
	}
	mp, ok := stmt.Expression.(*ast.MapLiteral)
	if !ok {
		t.Errorf("stmt.Expression not *ast.MapLiteral: got=%T", stmt.Expression)
	}

	if len(mp.Pairs) != 2 {
		t.Errorf("mp.Pairs must be 2 statements: got=%d", len(mp.Pairs))
	}

	want := map[string]string{"foo": "bar", "baz": "qux"}
	for key, val := range mp.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key not *ast.StringLiteral: got=%T", key)
		}

		value, ok := val.(*ast.StringLiteral)
		if !ok {
			t.Errorf("val not *ast.StringLiteral: got=%T", val)
		}

		if value.String() != want[literal.String()] {
			t.Errorf("value not equal to %s: got=%s", want[literal.String()], value.String())
		}
	}
}

func TestForExpression(t *testing.T) {
	got := "mew (fax) { rizz 1; }"
	l := lexer.NewLexer(got)
	p := NewParser(l)
	program := p.Parse()
	checkParser(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] not *ast.ExpressionStatement: got=%T", program.Statements[0])
	}
	loop, ok := stmt.Expression.(*ast.ForExpression)
	if !ok {
		t.Errorf("stmt.Expression not *ast.ForExpression: got=%T", stmt.Expression)
	}

	_, ok = loop.Condition.(*ast.Boolean)
	if !ok {
		t.Errorf("loop.Condition not *ast.Boolean: got=%T", loop.Condition)
	}

	if len(loop.Body.Statements) != 1 {
		t.Errorf("loop.Body.Statements must be 1 statement: got=%d", len(loop.Body.Statements))
	}
	_, ok = loop.Body.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Errorf("loop.Body.Statements[0] not *ast.ReturnStatement: got=%T", loop.Body.Statements[0])
	}
}