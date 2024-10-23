package lexer

import (
	"testing"

	"github.com/dxtym/maymun/token"
)

func TestNextToken(t *testing.T) {
	input := `deylik ayirish = amal(x, y) { x - y }; deylik a = "foo bar";`

	tests := []struct {
		got  token.TokenType
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
		{token.MINUS, "-"},
		{token.IDENT, "y"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "deylik"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.STRING, "foo bar"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := NewLexer(input) // creates new lexer
	for _, tt := range tests {
		token := l.NextToken() // gets next token
		if token.Type != tt.got {
			t.Errorf("token.Type %q ga teng emas: %q", tt.got, token.Type)
		}
		if token.Literal != tt.want {
			t.Errorf("token.Literal %q ga teng emas: %q", tt.want, token.Literal)
		}
	}
}
