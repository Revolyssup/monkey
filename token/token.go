package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

const (
	//keywords
	LET      = "LET"
	FUNCTION = "FUNCTION"

	//Operators
	PLUS   = "+"
	ASSIGN = "="
	//delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LEFT_BRACKET  = "("
	RIGHT_BRACKET = ")"
	LEFT_BRACE    = "{"
	RIGHT_BRACE   = "}"

	//identifier
	IDENTIFIER = "IDENT"

	//literal
	INTEGER = "INT"

	//special
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
)

//To check if given token is keyword or an identifier
func IdentOrKeyword(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
