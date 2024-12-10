package evaluator

import (
	"interpreter/lexer"
	"interpreter/object"
	"interpreter/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"12", 12},
		{"-5", -5},
		{"10", 10},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
	}
	for _, tt := range tests {
		if !testIntergerObject(t, testEval(tt.input), tt.expected) {
			return
		}
	}
}
func TestBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
	}

	for _, tt := range tests {
		testBoolean(t, testEval(tt.input), tt.expected)
	}
}
func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!4", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		eval := testEval(tt.input)
		testBoolean(t, eval, tt.expected)
	}
}

func TestIfElseStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		res := testEval(tt.input)

		switch res.(type) {
		case *object.Interger:
			i64 := int64(tt.expected.(int))
			testIntergerObject(t, res, i64)
		case *object.NULL:
			if tt.expected != nil {
				t.Errorf("res is not nil.got=%s", res.Type())
			}
		}
	}
}
func TestMinusOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"-5", -5},
		{"5", 5},
		{"--5", 5},
		{"---5", -5},
	}

	for _, tt := range tests {
		testIntergerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`if (10 > 1){
			if (10 > 1){
			return 10;
			}
			return 1;
			}
		`, 10},
	}
	for _, tt := range tests {
		obj := testEval(tt.input)
		switch obj.(type) {
		case *object.Interger:
			i64 := int64(tt.expected.(int))
			testIntergerObject(t, obj, i64)
		}
	}
}
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"true + false + true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return true + false;
  }

  return 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.ErrorType)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}
func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntergerObject(t, testEval(tt.input), tt.expected)
	}
}
func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}
func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {

		testIntergerObject(t, testEval(tt.input), tt.expected)
	}
}
func TestClosures(t *testing.T) {
	input := `
	let newAdder = fn(x){
	fn(y){x+y;};
	};
	
	let addTwo = newAdder(2);
	addTwo(2);
	`

	testIntergerObject(t, testEval(input), 4)
}
func TestStringLiteral(t *testing.T) {
	input := `"Hello World!";`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not string.got=%T", evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("str.Value !=Hello World!.got=%s", str.Value)
	}
}
func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " "+ "World!"`
	eval := testEval(input)
	str, ok := eval.(*object.String)
	if !ok {
		t.Fatalf("object is not string.got=%T", eval)
	}
	if str.Value != "Hello World!" {
		t.Errorf("str.Value !=Hello World!.got=%s", str.Value)
	}
}
func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one","two")`, `wrong number of arguments. got=2, want=1`},
	}

	for _, tt := range tests {
		eval := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntergerObject(t, eval, int64(expected))
		case string:
			errObj, ok := eval.(*object.ErrorType)
			if !ok {
				t.Errorf("object is not Error.got=%T", eval)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message,expected=%q,got=%q", expected, errObj.Message)
			}
		}
	}
}
func TestArrayLiteral(t *testing.T) {
	input := `[1,2*2,3+3]`
	eval := testEval(input)
	result, ok := eval.(*object.Array)
	if !ok {
		t.Fatalf("eval is not (*object.Array).got=%T", eval)
	}
	testIntergerObject(t, result.Elements[0], 1)
	testIntergerObject(t, result.Elements[1], 4)
	testIntergerObject(t, result.Elements[2], 6)
}
func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntergerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}
func TestMap(t *testing.T) {
	input := `
	let map = fn(arr, f) {
	let iter = fn(arr, accumulated) {
	if (len(arr) == 0) {
	accumulated
	} else {
	iter(rest(arr), push(accumulated, f(first(arr))));
			}
	};
	iter(arr, []);
	};
	let a = [1, 2, 3, 4, 5];
	let double = fn(x) { x * 2 };
	map(a, double);
	`
	eval := testEval(input)
	arr, ok := eval.(*object.Array)
	if ok {
		for i, ele := range arr.Elements {
			testIntergerObject(t, ele, int64(2*i+2))
		}
	}
}
func TestReduce(t *testing.T) {
	input := `
	let reduce = fn(arr, initial, f) {
	let iter = fn(arr, result) {
		if (len(arr) == 0) {
		result
		} else {
	iter(rest(arr), f(result, first(arr)));
	}
	};
	iter(arr, initial);
	};
	let sum = fn(arr) {
	reduce(arr, 0, fn(initial, el) { initial + el });
	};
	sum([1, 2, 3, 4, 5]);
	`
	eval := testEval(input)
	in, ok := eval.(*object.Interger)
	if ok {
		testIntergerObject(t, in, int64(15))
	}
}
func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Interger{Value: 4}).HashKey():     4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntergerObject(t, pair.Value, expectedValue)
	}
}
func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntergerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}
func testEval(input string) object.Object {
	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()
	env := object.NewEnvironment(nil)
	return Eval(program, env)
}
func testNullObject(t *testing.T, obj object.Object) bool {
	if obj.Type() != object.NULL_OBJ {
		t.Fatalf("%s != NULL.", obj.Inspect())
		return false
	}
	return true
}
func testIntergerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Interger)

	if !ok {
		t.Errorf("obj is not *object.Interger.got=%T", obj)
		t.Errorf(obj.(*object.ErrorType).Message)
		return false
	}

	if result.Value != expected {
		t.Errorf("result.Value != %d,got=%d", expected, result.Value)
		return false
	}
	return true
}
func testBoolean(t *testing.T, obj object.Object, expected bool) bool {
	bl, ok := obj.(*object.BooleanType)
	if !ok {
		t.Errorf("obj is not (*object.BooleanType).got=%T", obj)
		return false
	}

	if bl.Value != expected {
		t.Errorf("want %t,got=%t", expected, bl.Value)
		return false
	}
	return true
}
