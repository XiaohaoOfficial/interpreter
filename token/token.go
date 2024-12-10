package token

// 定义了表达式中，变量，关键词，操作符
type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF" //END OF FILE
	STRING  = "STRING"
	IDENT   = "IDENT" //i j foo
	INT     = "INT"   //1 2 3 123

	BANG      = "!"
	ASSIGN    = "="
	PLUS      = "+"
	MINUS     = "-"
	SLASH     = "/"
	ASTERISK  = "*"
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"

	LT       = "<"
	GT       = ">"
	LT_OR_EQ = "<="
	GT_OR_EQ = ">="
	EQ       = "=="
	NOT_EQ   = "!="

	IF       = "if"
	ELSE     = "else"
	RETURN   = "RETURN"
	TRUE     = "true"
	FALSE    = "false "
	FUNCTION = "FUNCTION"
	LET      = "LET"
)
