package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

// available token types
const (
	LET    = "LET"
	IDENT  = "IDENT"
	FUNC   = "FUNC"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RETURN"

	INT    = "INT"
	STRING = "STRING"

	ASSIGN   = "="
	PLUS     = "+"
	NOT      = "!"
	MINUS    = "-"
	DIV      = "/"
	MUL      = "*"
	LESS     = "<"
	MORE     = ">"
	EQUAL    = "=="
	NOTEQUAL = "!="

	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"

	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
)

// TODO: change to uzbek
var keywords = map[string]TokenType{
	"amogus": LET,
	"brainrot":   FUNC,
	"kino": TRUE,
	"slop": FALSE,
	"hawk":   IF,
	"tuah":   ELSE,
	"rizz": RETURN,
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
