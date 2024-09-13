package lexer

import (
	"testing"

	"github.com/dxtym/monke/token"
)

func TestNextToken(t *testing.T) {
	input := `let add = func(x, y) { x + y }; 
	!-/*<>
	if (5 < 10) { return true; } else { return false; }
	10 == 10
	9 != 10`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNC, "func"},
		{token.LBRACKET, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RBRACKET, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.NOT, "!"},
		{token.MINUS, "-"},
		{token.DIV, "/"},
		{token.MUL, "*"},
		{token.LESS, "<"},
		{token.MORE, ">"},
		{token.IF, "if"},
		{token.LBRACKET, "("},
		{token.INT, "5"},
		{token.LESS, "<"},
		{token.INT, "10"},
		{token.RBRACKET, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQUAL, "=="},
		{token.INT, "10"},
		{token.INT, "9"},
		{token.NOTEQUAL, "!="},
		{token.INT, "10"},
		{token.EOF, ""},
	}

	l := NewLexer(input) // creates new lexer
	for i, tt := range tests {
		token := l.NextToken() // gets next token
		if token.Type != tt.expectedType {
			t.Fatalf("test[%d] - wrong type. want=%q, got=%q", i, tt.expectedType, token.Type)
		}
		if token.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - wrong literal. want=%q, got=%q", i, tt.expectedLiteral, token.Literal)
		}
	}
}
