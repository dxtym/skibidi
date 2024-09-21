package parser

import (
	"testing"

	"github.com/dxtym/monke/ast"
	"github.com/dxtym/monke/lexer"
	"github.com/dxtym/monke/token"
)

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
		t.Fatalf("program.Statements doesn't contain 2 statements. got=%d", len(program.Statements)) // TODO: modify to 3 after parsing expressions
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