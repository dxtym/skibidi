package lexer

import (
	"testing"

	"github.com/dxtym/maymun/token"
)

func TestNextToken(t *testing.T) {
	input := "deylik ayirish = amal(x, y) { x - y };"

	tests := []struct {
		got    token.TokenType
		want string
	}{
		{token.LET, "deylik"},
		{token.IDENT, "ayirish"},
		{token.ASSIGN, "="},
		{token.FUNC, "amal"},
		{token.LBRACKET, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RBRACKET, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "-"},
		{token.IDENT, "y"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := NewLexer(input) // creates new lexer
	for _, tt := range tests {
		token := l.NextToken() // gets next token
		if token.Type != tt.got {
			t.Fatalf("token.Type not equal to %s: got=%s", tt.got, token.Type)
		}
		if token.Literal != tt.want {
			t.Fatalf("token.Literal not equal to %s: got=%s", tt.want, token.Literal)
		}
	}
}
