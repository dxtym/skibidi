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
		{"ijobiy", true},
		{"salbiy", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"ijobiy == ijobiy", true},
		{"salbiy == salbiy", true},
		{"ijobiy == salbiy", false},
		{"ijobiy != salbiy", true},
		{"salbiy != ijobiy", true},
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
		{"!ijobiy", false},
		{"!salbiy", true},
		{"!!ijobiy", true},
		{"!!salbiy", false},
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
		{"agar (1) {2};", 2},
		{"agar (ijobiy) {1};", 1},
		{"agar (2 > 1) {2};", 2},
		{"agar (salbiy) {1};", nil},
		{"agar (1 > 2) {2} yana {1};", 1},
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
		{"qaytar 1; 2;", 1},
		{"1 * 2; qaytar 2; 1;", 2},
		{"qaytar 2; 2 * 1;", 2},
		{"qaytar 1; qaytar 2;", 1},
		{"amal(x, y) { x + y; }(1, 2)", 3},
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
		{"1 + ijobiy;", "type mismatch: INTEGER + BOOLEAN"},
		{"-ijobiy;", "unknown operator: -BOOLEAN"},
		{"ijobiy + salbiy;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"1 - ijobiy; 1;", "type mismatch: INTEGER - BOOLEAN"},
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
		{"deylik a = 1; a;", 1},
		{"deylik a = 1; deylik b = a; b", 1},
		{"deylik a = 1 + 2; deylik b = a + 1; b;", 4},
	}

	for _, tt := range tests {
		evaled := testEval(tt.got)
		testIntegerObject(t, evaled, tt.want)
	}
}

func TestFunction(t *testing.T) {
	got := "amal(x) { x + 2; }"
	evaled := testEval(got)

	fn, ok := evaled.(*object.Function)
	if !ok {
		t.Fatalf("evaled not *object.Function: got=%T", evaled)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("fn.Parameters must be 1 statement: got=%d", len(fn.Parameters))
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("fn.Parameters[0].String not equal to x: got=%s", fn.Parameters[0].String())
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
		{"deylik a = amal(x) { x + 1; }; a(1);", 2},
		{"deylik a = amal(x) { x + 1; }(1); a;", 2},
		{"deylik a = amal(x, y) { qaytar x + y; }; a(1, 2);", 3},
		{"deylik a = amal(x) { amal(y) { x + y; } }; a(1)(2);", 3},
	}

	for _, tt := range tests {
		evaled := testEval(tt.got)
		testIntegerObject(t, evaled, tt.want)
	}
}