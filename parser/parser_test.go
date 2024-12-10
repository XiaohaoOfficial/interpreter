package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 8888;
	`

	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	chenckParserErrors(t, parser)
	if program == nil {
		t.Fatalf("ParseProgram return nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program statements does not contain 3 Statements.got =%d", len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
		expectedValue      string
	}{
		{"x", "5"},
		{"y", "10"},
		{"foobar", "8888"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier, tt.expectedValue) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string, value string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral isn't let.got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement.got=%T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'.got='%s' ", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'.got='%s' ", name, letStmt.Name.TokenLiteral())
		return false
	}
	if letStmt.Value.String() != value {
		t.Errorf("letStmt.Value.String() !=  value.got=%s", letStmt.Value.String())
	}
	return true
}

func chenckParserErrors(t *testing.T, p *Parser) {
	errors := p._errors
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error :%q", msg)
	}
	t.FailNow()
}

func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 99933;
	`

	lexer := lexer.New(input)
	p := New(lexer)

	program := p.ParseProgram()
	chenckParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program statements does not contain 3 Statements.got =%d", len(program.Statements))
	}
	expectedValue := []string{"5", "10", "99933"}
	for i, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement.got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return' got:%q", returnStmt.TokenLiteral())
		}
		if returnStmt.ReturnValue.String() != expectedValue[i] {
			t.Errorf("returnStmt.ReturnValue.String() !=%s,got=%s", expectedValue[i], returnStmt.ReturnValue.String())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	chenckParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements.got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement,got =%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.Identifier,got =%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value is not foobar,got=%s", ident.Value)
	}

	if ident.Token.Literal != "foobar" {
		t.Errorf("ident.Token.Literal is not foobar,got=%s", ident.Token.Literal)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`
	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	chenckParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements.got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement,got =%T", program.Statements[0])
	}

	liberal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IntegerLiteral,got =%T", stmt.Expression)
	}

	if liberal.Value != 5 {
		t.Errorf("ident.Value is not foobar,got=%d", liberal.Value)
	}

	if liberal.Token.Literal != "5" {
		t.Errorf("ident.Token.Literal is not foobar,got=%s", liberal.Token.Literal)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	inputs := []struct {
		Input         string
		Operator      string
		IntergerValue interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true", "!", true},
	}

	for _, tt := range inputs {
		lexer := lexer.New(tt.Input)
		parse := New(lexer)
		program := parse.ParseProgram()
		chenckParserErrors(t, parse)

		if len(program.Statements) != 1 {
			t.Errorf("len(program.Statements) != 1,got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Errorf("program.Statements[0] is not *ast.ExpressionStatement,got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Errorf("stmt.Expression is not *ast.PrefixExpression,got=%T", stmt.Expression)
		}

		testLiteralExpression(t, exp.Right, tt.IntergerValue)

	}
}

func testIntergerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	intergerLiteral, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il is not *ast.IntegerLiteral,got=%T", il)
		return false
	}

	if intergerLiteral.Value != value {
		t.Errorf("il.value is not %d, got=%d", value, intergerLiteral.Value)
		return false
	}
	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5+5;", 5, "+", 5},
		{"5-5;", 5, "-", 5},
		{"5*5;", 5, "*", 5},
		{"5/5;", 5, "/", 5},
		{"5>5;", 5, ">", 5},
		{"5<5;", 5, "<", 5},
		{"5==5;", 5, "==", 5},
		{"5!=5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		lexer := lexer.New(tt.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		chenckParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("len(program.Statements)!=1.got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not (*ast.Expression).got=%T", program.Statements[0])
		}

		infixExpress, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt is not (*ast.InfixExpression).got=%T", stmt)
		}

		if !testIntergerLiteral(t, infixExpress.Left, tt.leftValue) {
			return
		}

		if infixExpress.Operator != tt.operator {
			t.Fatalf("infixExpress.Operator != tt.operator,got=%T", infixExpress.Operator)
		}

		if !testIntergerLiteral(t, infixExpress.Right, tt.rightValue) {
			return
		}

	}
}
func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a *[1,2,3,4][b*c]*d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		lexer := lexer.New(tt.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		chenckParserErrors(t, parser)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y){
	x ;
	}else{
	x+y;
	}
	`

	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	chenckParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) != 1,got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not (*ast.ExpressionStatement) got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not (*ast.IfExpression).got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("len(exp.Consequence.Statements)!=1,got=%d", len(exp.Consequence.Statements))
	}

	consequese, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not (*ast.ExpressionStatement).got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequese.Expression, "x") {
		return
	}

	stmt = exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	ep := stmt.Expression.(*ast.InfixExpression)
	if !testInfixExpression(t, ep, "x", "+", "y") {
		return
	}
}

func TestFunctionalLiteralParsing(t *testing.T) {
	input := `fn(x,y){x + y;}`

	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	chenckParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) != 1.got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not (*ast.ExpressionStatement).got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not (*ast.FunctionLiteral).got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("len(function.Parameters)!=2.got=%d", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("len(function.Body.Statements)!=1.got=%d", len(function.Body.Statements))
	}
	testInfixExpression(t, function.Body.Statements[0].(*ast.ExpressionStatement).Expression, "x", "+", "y")
}

func TestFunctionParamterParsing(t *testing.T) {
	inputs := []struct {
		input         string
		expectedParas []string
	}{
		{input: "fn(){};", expectedParas: []string{}},
		{input: "fn(x){};", expectedParas: []string{"x"}},
		{input: "fn(x,y,z){};", expectedParas: []string{"x", "y", "z"}},
	}

	for _, tt := range inputs {
		lexer := lexer.New(tt.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		chenckParserErrors(t, parser)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		lib := stmt.Expression.(*ast.FunctionLiteral)

		if len(lib.Parameters) != len(tt.expectedParas) {
			t.Fatalf("len(lib.Parameters) != len(tt.expectedParas),expected=%d,got=%d", len(lib.Parameters), len(tt.expectedParas))
		}

		for i, ident := range tt.expectedParas {
			testLiteralExpression(t, lib.Parameters[i], ident)
		}

	}
}
func TestCallExpressionParsing(t *testing.T) {
	input := `add(1,2*3,4+5);`

	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	chenckParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("len(program.Statements) !=1.got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not(*ast.ExpressionStatement).got=%T", program.Statements[0])
	}

	ce, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not (*ast.CallExpression).got=%T", stmt.Expression)
	}

	if !testIdentifier(t, ce.Function, "add") {
		return
	}

	if len(ce.Arguments) != 3 {
		t.Fatalf("len(ce.Arguments)!=3,got=%d", len(ce.Arguments))
	}

	testLiteralExpression(t, ce.Arguments[0], 1)
	testInfixExpression(t, ce.Arguments[1], 2, "*", 3)
	testInfixExpression(t, ce.Arguments[2], 4, "+", 5)
}
func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	chenckParserErrors(t, parser)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	sl := stmt.Expression.(*ast.StringLiteral)

	if sl.Value != "hello world" {
		t.Errorf("sl.Value != hello world.got=%s", sl.Value)
	}
}
func TestParsingArrayLiterals(t *testing.T) {
	input := `[1,2*2,3+3];`

	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	chenckParserErrors(t, parser)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0].(*ast.ExpressionStatement),got=%T", program.Statements[0])
	}
	al, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not (*ast.ArrayLiteral).got=%T", stmt.Expression)
	}

	if len(al.Elements) != 3 {
		t.Fatalf("len(al.Elements)!=3.got=%d", len(al.Elements))
	}
	testIntergerLiteral(t, al.Elements[0], 1)
	testInfixExpression(t, al.Elements[1], 2, "*", 2)
	testInfixExpression(t, al.Elements[2], 3, "+", 3)
}
func TestParsingIndexExpressions(t *testing.T) {
	input := `myArray[1+1]`
	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	chenckParserErrors(t, parser)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0].(*ast.ExpressionStatement),got=%T", program.Statements[0])
	}
	ie, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not (*ast.IndexExpression).got=%T", stmt.Expression)
	}

	if !testIdentifier(t, ie.Left, "myArray") {
		return
	}
	if !testInfixExpression(t, ie.Index, 1, "+", 1) {
		return
	}

}
func TestParsingHashLiteral(t *testing.T) {
	input := `{"one" : 1, "two" : 2, "three" : 3};`
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	chenckParserErrors(t, parser)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not (*ast.ExpressionStatement). got=%T", program.Statements[0])
	}
	hl, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not(*ast.HashLiteral).got=%T", stmt.Expression)
	}
	i := 0
	for key, value := range hl.Pairs {

		testIntergerLiteral(t, value, expected[key.String()])
		i++
	}
}
func TestParsingHashLiteralWithExpressions(t *testing.T) {
	input := `{"one" : 0 + 1,"two" : 10 - 8 ,"three" : 15 /5 }`
	expected := map[string]struct {
		left     int
		operator string
		right    int
	}{
		"one":   {0, "+", 1},
		"two":   {10, "-", 8},
		"three": {15, "/", 5},
	}
	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	chenckParserErrors(t, parser)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not (*ast.ExpressionStatement). got=%T", program.Statements[0])
	}
	hl, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not(*ast.HashLiteral).got=%T", stmt.Expression)
	}
	for key, value := range hl.Pairs {
		testInfixExpression(t, value, expected[key.String()].left, expected[key.String()].operator, expected[key.String()].right)
	}
}
func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntergerLiteral(t, exp, int64(v))
	case int64:
		return testIntergerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolLiteral(t, exp, v)
	}
	t.Errorf("type of exp does not fit.got=%T", exp)
	return false
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier.got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value != %s,got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() != %s,got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not *ast.InfixExpression,got=%T", exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("opExp.Operator != %s,got=%s", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}
func testBoolLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bl, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp is not *ast.Boolean,got=%T", exp)
		return false
	}

	if bl.Value != value {
		t.Errorf("bo.Value != %t,got=%t", value, bl.Value)
		return false
	}

	if bl.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bl.TokenLiteral() !=%s.got=%s", fmt.Sprintf("%t", value), bl.TokenLiteral())
		return false
	}

	return true
}
