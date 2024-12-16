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
	COLON     = ":"
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

// TODO: redefine to brainrot -> mew
var keywords = map[string]TokenType{
	"amogus": LET,
	"cook":   FUNC,
	"fax":    TRUE,
	"cap":    FALSE,
	"hawk":   IF,
	"tuah":   ELSE,
	"rizz":   RETURN,
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
