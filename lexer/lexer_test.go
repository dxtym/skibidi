package lexer

import (
	"testing"

	"github.com/dxtym/monke/token"
)

func TestNextToken(t *testing.T) {
	input := "let add = func(x, y) { x + y };"
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
