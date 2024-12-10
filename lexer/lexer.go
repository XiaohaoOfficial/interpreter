package lexer

import (
	"interpreter/token"
)

type Lexer struct {
	_input        string
	_position     int
	_readPosition int
	_ch           byte
}

func New(input string) *Lexer {
	lexer := &Lexer{_input: input}
	lexer.readChar()

	return lexer
}

func (l *Lexer) readChar() {
	if l._readPosition >= len(l._input) {
		//注意！这里一定是0而不是'0'
		l._ch = 0
	} else {
		l._ch = l._input[l._readPosition]
	}
	l._position = l._readPosition
	l._readPosition += 1
}

var keywords = map[string]token.TokenType{
	"let":    token.LET,
	"fn":     token.FUNCTION,
	"if":     token.IF,
	"return": token.RETURN,
	"true":   token.TRUE,
	"false":  token.FALSE,
	"else":   token.ELSE,
}

func newToken(tpe token.TokenType, ch byte) token.Token {
	return token.Token{Type: tpe, Literal: string(ch)}
}
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	switch l._ch {
	case ':':
		tok = newToken(token.COLON, l._ch)
	case '=':
		l.readChar()
		if l._ch == '=' {
			l.readChar()
			tok.Literal = "=="
			tok.Type = token.EQ
		} else {
			tok = newToken(token.ASSIGN, '=')
		}
		return tok
	case ';':
		tok = newToken(token.SEMICOLON, l._ch)
	case '(':
		tok = newToken(token.LPAREN, l._ch)
	case ')':
		tok = newToken(token.RPAREN, l._ch)
	case '{':
		tok = newToken(token.LBRACE, l._ch)
	case '}':
		tok = newToken(token.RBRACE, l._ch)
	case ',':
		tok = newToken(token.COMMA, l._ch)
	case '+':
		tok = newToken(token.PLUS, l._ch)
	case '-':
		tok = newToken(token.MINUS, l._ch)
	case '*':
		tok = newToken(token.ASTERISK, l._ch)
	case '/':
		tok = newToken(token.SLASH, l._ch)
	case '[':
		tok = newToken(token.LBRACKET, l._ch)
	case ']':
		tok = newToken(token.RBRACKET, l._ch)
	case '<':
		l.readChar()
		if l._ch == '=' {
			tok.Literal = "<="
			tok.Type = token.LT_OR_EQ
			l.readChar()
		} else {
			tok = newToken(token.LT, '<')
		}
		return tok
	case '>':
		l.readChar()
		if l._ch == '=' {
			tok.Literal = ">="
			tok.Type = token.GT_OR_EQ
			l.readChar()
		} else {
			tok = newToken(token.GT, '>')
		}
		return tok
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	case '!':
		l.readChar()
		if l._ch == '=' {
			l.readChar()
			tok.Literal = "!="
			tok.Type = token.NOT_EQ
		} else {
			tok.Literal = "!"
			tok.Type = token.BANG
		}
		return tok
	case ' ':
		l.readChar()
		return l.NextToken()
	case '\n':
		l.readChar()
		return l.NextToken()
	case '\t':
		l.readChar()
		return l.NextToken()
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()

	default:
		//将关键字存储在一个map中，如果查找到key，则返回，查不到则设置为非法
		if isLetter(l._ch) {
			tok.Literal = l.readIden()
			if tpe, ok := keywords[tok.Literal]; ok {
				tok.Type = tpe
			} else {
				tok.Type = token.IDENT
			}
			return tok
		} else {
			if isNum(l._ch) {
				tok.Literal = l.readNum()
				tok.Type = token.INT
				return tok
			}
			tok = newToken(token.ILLEGAL, l._ch)

		}
	}
	l.readChar()
	return tok
}
func (l *Lexer) readIden() string {
	position := l._position
	if isLetter(l._ch) {
		for isLetter(l._ch) {
			l.readChar()
		}
	}
	return l._input[position:l._position]
}
func (l *Lexer) readNum() string {
	position := l._position
	for isNum(l._ch) {
		//fmt.Printf("char : %c isNum:%v", l._ch, isNum(l._ch))
		l.readChar()
	}

	return l._input[position:l._position]
}
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
func isNum(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
func (l *Lexer) readString() string {
	l.readChar()
	start := l._position
	for l._ch != '"' && l._ch != 0 {
		l.readChar()
	}
	str := l._input[start:l._position]
	return str
}
