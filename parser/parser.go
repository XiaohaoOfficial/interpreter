package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	ASSIGN
	EQUALS     //==
	LESSGRATER // > or <
	SUM        // +
	PRODUCT    // *
	PREFIX     // ++
	CALL       // myFunction(a,b)
	INDEX
)

var precedences = map[token.TokenType]int{
	token.ASSIGN:   ASSIGN,
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGRATER,
	token.GT:       LESSGRATER,
	token.LT_OR_EQ: LESSGRATER,
	token.GT_OR_EQ: LESSGRATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

// 语法分析器
type Parser struct {
	_lexer          *lexer.Lexer
	_curToken       token.Token
	_peekToken      token.Token
	_errors         []string
	_prefixParseFns map[token.TokenType]prefixParseFn
	_infixParseFns  map[token.TokenType]infixParseFn
}

// 向前缀函数和中缀函数map中注册方法
func (p *Parser) registerPrefixParseFn(tpe token.TokenType, fn prefixParseFn) {
	p._prefixParseFns[tpe] = fn
}
func (p *Parser) registerInfixParseFn(tpe token.TokenType, fn infixParseFn) {
	p._infixParseFns[tpe] = fn
}
func (p *Parser) nextToken() {
	p._curToken = p._peekToken
	p._peekToken = p._lexer.NextToken()
}
func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{_lexer: lexer}

	p._prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefixParseFn(token.IDENT, p.parseIdentifier)
	p.registerPrefixParseFn(token.INT, p.parseIntergerLiberal)
	p.registerPrefixParseFn(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixParseFn(token.BANG, p.parsePrefixExpression)
	p.registerPrefixParseFn(token.TRUE, p.parseBoolean)
	p.registerPrefixParseFn(token.FALSE, p.parseBoolean)
	p.registerPrefixParseFn(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefixParseFn(token.IF, p.parseIfExpression)
	p.registerPrefixParseFn(token.FUNCTION, p.parseFunctionalLiteral)
	p.registerPrefixParseFn(token.STRING, p.parseStringLiteral)
	p.registerPrefixParseFn(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefixParseFn(token.LBRACE, p.parseHashingLiteral)

	p._infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfixParseFn(token.PLUS, p.parseInfixExpression)
	p.registerInfixParseFn(token.MINUS, p.parseInfixExpression)
	p.registerInfixParseFn(token.SLASH, p.parseInfixExpression)
	p.registerInfixParseFn(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixParseFn(token.EQ, p.parseInfixExpression)
	p.registerInfixParseFn(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfixParseFn(token.LT, p.parseInfixExpression)
	p.registerInfixParseFn(token.GT, p.parseInfixExpression)
	p.registerInfixParseFn(token.LT_OR_EQ, p.parseInfixExpression)
	p.registerInfixParseFn(token.GT_OR_EQ, p.parseInfixExpression)
	p.registerInfixParseFn(token.LBRACKET, p.parseIndexExpression)
	p.registerInfixParseFn(token.ASSIGN, p.parseInfixExpression)
	p.registerInfixParseFn(token.LPAREN, p.parseCallExpression)

	//滑动两次，以初始化curToken和peekToken
	p.nextToken()
	p.nextToken()

	return p
}
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p._curToken, Value: p.curTokenIs(token.TRUE)}
}
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}
	return program
}
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p._curToken,
		Operator: p._curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}
func (p *Parser) parseStatement() ast.Statement {
	switch p._curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

/*
以 let a = 10 ; 为例
进入该方法后，首先记录该语句的Token
对于一个Statement而言，需要记录的有以下三点
Token Name Value
Token 在调用该方法的时候即保存，对于该语句即为let
Name 在检测到变量时保存，对于该语句即为a
Value :TODO
*/
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p._curToken}

	if !p.expectedPeek(token.IDENT) {
		p.peekError(token.IDENT)
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p._curToken, Value: p._curToken.Literal}
	if !p.expectedPeek(token.ASSIGN) {
		p.peekError(token.ASSIGN)
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p._curToken}

	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)
	for !p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	p.nextToken()

	return stmt
}

/*
parseExpression方法用于解析前缀与中缀表达式，只通过观察代码理解较为抽象，因此举例说明
初始parseExpression传入的都是LOWEST
对于前缀表达式
5;	-15;	!a;
prefix会解析相应的前缀表达式并存入left中，在for循环中由于后续是分号，因此会返回
对于中缀表达式
5 + 10;
5首先会被识别为前缀表达式，并存入left
然后调用infix方法 调用的过程是 将 “+”号存入operator，然后递归调用parseExpression方法
10在递归调用中，还是会被识别为前缀表达式
因此整个Expreesion的结构为Left:5 Operator:"+" Right:10
对于复杂的中缀表达式
5 + 6 * 7；
同样，5首先会被识别为前缀表达式，并存入left
然后调用infix方法，将"+"号存入operator，递归调用parseExpression方法，此时的precedent为"+"号的
由于"+"的优先级低于"*"，因此第二轮for循环能够再次进行
此时Left:5 operator:"+" Right:(Left:6 Operator:"*" Right:7)
而如果是5 * 6 + 7;
由于"*"的优先级高于"+"，因此结构为
Left:(Left:5 Operator:"*" Right:6) Operator:"+" Right:7
*/
func (p *Parser) parseExpression(precedent int) ast.Expression {
	prefix := p._prefixParseFns[p._curToken.Type]

	if prefix == nil {
		err := "noPrefixParseFnErr! for:" + p._curToken.Literal
		p._errors = append(p._errors, err)
		return nil
	}
	left := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedent < p.peekPrecedence() {
		infix := p._infixParseFns[p._peekToken.Type]
		if infix == nil {
			return left
		}

		p.nextToken()

		left = infix(left)
	}
	return left

}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p._curToken}

	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p._curToken, Value: p._curToken.Literal}
}
func (p *Parser) parseIntergerLiberal() ast.Expression {

	value, err := strconv.ParseInt(p._curToken.Literal, 0, 64)
	if err != nil {
		p._errors = append(p._errors, err.Error())
	}
	return &ast.IntegerLiteral{Token: p._curToken, Value: value}
}
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p._curToken,
		Operator: p._curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p._curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p._peekToken.Type == t
}

func (p *Parser) expectedPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	err := fmt.Sprintf("expected next token to be %s,but got:%s instead", t, p._peekToken.Type)
	p._errors = append(p._errors, err)
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p._peekToken.Type]; ok {
		return p
	}
	return LOWEST
}
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p._curToken.Type]; ok {
		return p
	}
	return LOWEST
}
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectedPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p._curToken}

	if !p.expectedPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectedPeek(token.RPAREN) {
		return nil
	}

	if !p.expectedPeek(token.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectedPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p._curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}
	p.nextToken()

	ident := &ast.Identifier{Token: p._curToken, Value: p._curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p._curToken, Value: p._curToken.Literal}
		identifiers = append(identifiers, ident)
	}
	if !p.expectedPeek(token.RPAREN) {
		return nil
	}
	return identifiers
}
func (p *Parser) parseFunctionalLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p._curToken}

	if !p.expectedPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectedPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p._curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)

	return exp
}

/*
之前用来获取函数调用时的参数，目前被parseExpressionList方法取代
*/

// func (p *Parser) parseCallArguments() []ast.Expression {
// 	args := []ast.Expression{}

// 	if p.expectedPeek(token.RPAREN) {
// 		return args
// 	}

// 	p.nextToken()
// 	args = append(args, p.parseExpression(LOWEST))

// 	for p.peekTokenIs(token.COMMA) {
// 		p.nextToken()
// 		p.nextToken()
// 		args = append(args, p.parseExpression(LOWEST))
// 	}

//		if !p.expectedPeek(token.RPAREN) {
//			return nil
//		}
//		return args
//	}
func (p *Parser) Errors() []string {
	return p._errors
}
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p._curToken, Value: p._curToken.Literal}

}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p._curToken}

	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	var elements = []ast.Expression{}

	if p.expectedPeek(end) {
		return elements
	}
	p.nextToken()
	elements = append(elements, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		elements = append(elements, p.parseExpression(LOWEST))
	}
	if !p.expectedPeek(end) {
		return nil
	}
	return elements
}
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p._curToken, Left: left}
	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectedPeek(token.RBRACKET) {
		return nil
	}

	return exp
}
func (p *Parser) parseHashingLiteral() ast.Expression {
	hl := &ast.HashLiteral{Token: p._curToken, Pairs: make(map[ast.Expression]ast.Expression)}
	if p.expectedPeek(token.RBRACE) {
		return hl
	}
	p.nextToken()
	key := p.parseExpression(LOWEST)
	if !p.expectedPeek(token.COLON) {
		return nil
	}
	p.nextToken()

	value := p.parseExpression(LOWEST)

	hl.Pairs[key] = value
	for p.expectedPeek(token.COMMA) {
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectedPeek(token.COLON) {
			return nil
		}
		p.nextToken()
		value := p.parseExpression(LOWEST)
		hl.Pairs[key] = value
	}

	if !p.expectedPeek(token.RBRACE) {
		return nil
	}
	return hl
}
