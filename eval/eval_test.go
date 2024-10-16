package eval_test

import (
	"testing"

	"github.com/dxtym/maymun/eval"
	"github.com/dxtym/maymun/lexer"
	"github.com/dxtym/maymun/object"
	"github.com/dxtym/maymun/parser"
)

var (
	env = object.NewEnvironment()

	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func testEval(in string) object.Object {
	l := lexer.NewLexer(in)
	p := parser.NewParser(l)
	program := p.Parse()
	env := object.NewEnvironment()

	return eval.Eval(program, env)
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		got  string
		want int
	}{
		{"1", 1},
		{"2", 2},
		{"-1", -1},
		{"-2", -2},
		{"1 + 2", 3},
		{"2 - 1", 1},
		{"2 * 2", 4},
		{"4 / 2", 2},
		{"(1 + 2) * 3", 9},
		{"6 / (1 - 3)", -3},
	}

	for _, tt := range tests {
		eval := testEval(tt.got)
		testIntegerObject(t, eval, tt.want)
	}
}

func testIntegerObject(t *testing.T, eval object.Object, out int) bool {
	res, ok := eval.(*object.Integer)
	if !ok {
		t.Fatalf("eval not *object.Integer: got=%T", eval)
		return false
	}
	if res.Value != out {
		t.Fatalf("res.Value not equal to %d: got=%d", out, res.Value)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		got  string
		want bool
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
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
	}

	for _, tt := range tests {
		eval := testEval(tt.got)
		testBooleanObject(t, eval, tt.want)
	}
}

func testBooleanObject(t *testing.T, eval object.Object, out bool) bool {
	res, ok := eval.(*object.Boolean)
	if !ok {
		t.Fatalf("eval not *object.Boolean: got=%T", eval)
		return false
	}
	if res.Value != out {
		t.Fatalf("res.Value not equal to %t: got=%t", out, res.Value)
		return false
	}
	return true
}

func TestNotOperator(t *testing.T) {
	tests := []struct {
		got  string
		want bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!1", false},
	}

	for _, tt := range tests {
		evaled := testEval(tt.got)
		testBooleanObject(t, evaled, tt.want)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		got string
		want any
	}{
		{"if (1) {2};", 2},
		{"if (true) {1};", 1},
		{"if (2 > 1) {2};", 2},
		{"if (false) {1};", nil},
		{"if (1 > 2) {2} else {1};", 1},
	}

	for _, tt := range tests {
		evaled := testEval(tt.got)
		num, ok := tt.want.(int)
		if ok {
			testIntegerObject(t, evaled, num)
		} else {
			testNullObject(t, evaled)
		}
	}
}

func testNullObject(t *testing.T, eval object.Object) bool {
	if eval != NULL {
		t.Fatalf("object.Object not equal to NULL: got=%T", eval)
		return false
	}
	return true
}

func TestReturnValue(t *testing.T) {
	tests := []struct {
		got string
		want int
	}{
		{"return 1; 2;", 1},
		{"1 * 2; return 2; 1;", 2},
		{"return 2; 2 * 1;", 2},
		{"return 1; return 2;", 1},
	}

	for _, tt := range tests {
		evaled := testEval(tt.got)
		testIntegerObject(t, evaled, tt.want)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		got string
		want string
	}{
		{"1 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true;", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"1 - true; 1;", "type mismatch: INTEGER - BOOLEAN"},
		{"foobar;", "unbound indentifier: foobar"},
	}

	for _, tt := range tests {
		evaled := testEval(tt.got)
		err, ok := evaled.(*object.Error)
		if !ok {
			t.Fatalf("evaled not *object.Error: got=%T", evaled)
		}

		if err.Message != tt.want {
			t.Fatalf("err.Message not equal to %s: got=%s", tt.want, err.Message)
		}
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		got string
		want int
	}{
		{"let a = 1; a;", 1},
		{"let a = 1; let b = a; b", 1},
		{"let a = 1 + 2; let b = a + 1; b;", 4},
	}

	for _, tt := range tests {
		evaled := testEval(tt.got)
		testIntegerObject(t, evaled, tt.want)
	}
}

func TestFunction(t *testing.T) {
	got := "func(x) { x + 2; }"
	evaled := testEval(got)

	fn, ok := evaled.(*object.Function)
	if !ok {
		t.Fatalf("evaled not *object.Function: got=%T", evaled)
	}

	if len(fn.Arguments) != 1 {
		t.Fatalf("fn.Arguments must be 1 statement: got=%d", len(fn.Arguments))
	}
	if fn.Arguments[0].String() != "x" {
		t.Fatalf("fn.Arguments[0].String not equal to x: got=%s", fn.Arguments[0].String())
	}

	if fn.Body.String() != "(x + 2)" {
		t.Fatalf("fn.Body.String not equal to (x + 2): got=%s", fn.Body.String())
	}
}

func TestCallExpression(t *testing.T) {
	tests := []struct {
		got string
		want int
	}{
		{"let a = fn(x) { x + 1; }; a(1);", 2},
		{"let a = fn(x) { x + 1; }(1); a;", 2},
		{"let a = fn(x, y) { return x + y; }; a(1, 2);", 3},
		{"let a = fn(x) { fn(y) { x + y; } }; a(1)(2);", 3},
	}

	for _, tt := range tests {
		evaled := testEval(tt.got)
		testIntegerObject(t, evaled, tt.want)
	}
}