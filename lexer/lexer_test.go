package lexer

import (
	"testing"

	"github.com/dxtym/maymun/token"
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
	for _, tt := range tests {
		token := l.NextToken() // gets next token
		if token.Type != tt.expectedType {
			t.Fatalf("token.Type not equal to %s: got=%s", tt.expectedType, token.Type)
		}
		if token.Literal != tt.expectedLiteral {
			t.Fatalf("token.Literal not equal to %s: got=%s", tt.expectedLiteral, token.Literal)
		}
	}
}
