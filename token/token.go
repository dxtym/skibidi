package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

// available token types
const (
	LET   = "LET"
	IDENT = "IDENT"
	FUNC  = "FUNC"
	TRUE = "TRUE"
	FALSE = "FALSE"
	IF = "IF"
	ELSE = "ELSE"
	RETURN = "RETURN"

	INT    = "INT"
	ASSIGN = "="
	PLUS   = "+"
	NOT = "!"
	MINUS = "-"
	DIV = "/"
	MUL = "*"
	LESS = "<"
	MORE = ">"
	EQUAL = "=="
	NOTEQUAL = "!="

	COMMA     = ","
	SEMICOLON = ";"
	LBRACKET  = "("
	RBRACKET  = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
)

var keywords = map[string]TokenType{
	"let":  LET,
	"func": FUNC,
	"true": TRUE,
	"false": FALSE,
	"if": IF,
	"else": ELSE,
	"return": RETURN,
}

func NewToken(ttype TokenType, char byte) Token {
	return Token{
		Type:    ttype,
		Literal: string(char),
	}
}

func LookUpIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT // default to IDENT (x, e.g.)
}