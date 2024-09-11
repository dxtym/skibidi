package lexer

import "github.com/dxtym/monke/token"

type Lexer struct {
	input string
	pos   int // current pos pointing to char
	nxt   int // next pos after current pos
	char  byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// TODO: cover unicode
func (l *Lexer) readChar() {
	if l.nxt >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.nxt]
	}

	l.pos = l.nxt
	l.nxt++
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.char {
	case '=':
		tok = token.NewToken(token.ASSIGN, l.char)
	case '+':
		tok = token.NewToken(token.PLUS, l.char)
	case '(':
		tok = token.NewToken(token.LBRACKET, l.char)
	case ')':
		tok = token.NewToken(token.RBRACKET, l.char)
	case '{':
		tok = token.NewToken(token.LBRACE, l.char)
	case '}':
		tok = token.NewToken(token.RBRACE, l.char)
	case ';':
		tok = token.NewToken(token.SEMICOLON, l.char)
	case ',':
		tok = token.NewToken(token.COMMA, l.char)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdent()
			tok.Type = token.LookUpIdent(tok.Literal) // check if keyword
			return tok
		} else if isInteger(l.char) {
			tok.Literal = l.readIdent()
			tok.Type = token.INT
			return tok
		} else {
			tok = token.NewToken(token.ILLEGAL, l.char)
		}
	}

	l.readChar() // advance to next char
	return tok
}

func (l *Lexer) readIdent() string {
	start := l.pos
	for isLetter(l.char) || isInteger(l.char) {
		l.readChar()
	}
	return l.input[start:l.pos] // l.pos no longer letter or integer
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func isInteger(char byte) bool {
	return '0' <= char && char <= '9'
}
