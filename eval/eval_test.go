package eval_test

import (
	"testing"

	"github.com/dxtym/skibidi/eval"
	"github.com/dxtym/skibidi/lexer"
	"github.com/dxtym/skibidi/object"
	"github.com/dxtym/skibidi/parser"
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
		t.Errorf("eval not *object.Integer: got=%T", eval)
		return false
	}
	if res.Value != out {
		t.Errorf("res.Value not equal to %d: got=%d", out, res.Value)
		return false
	}
	return true
}

func TestStringExpression(t *testing.T) {
	input := `"hello, world"`
	evaled := testEval(input)
	str, ok := evaled.(*object.String)
	if !ok {
		t.Errorf("evaled not *object.String: got=%T", evaled)
	}

	if str.Value != "hello, world" {
		t.Errorf("str.Value not equal to %s: got=%s", "hello, world", str.Value)
	}
}

func TestStringConcatExpression(t *testing.T) {
	input := `"hello," + " " + "world"`
	evaled := testEval(input)
	str, ok := evaled.(*object.String)
	if !ok {
		t.Errorf("evaled not *object.String: got=%T", evaled)
	}

	if str.Value != "hello, world" {
		t.Errorf("str.Value not equal to %s: got=%s", "hello, world", str.Value)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		got  string
		want bool
	}{
		{"kino", true},
		{"slop", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"kino == kino", true},
		{"slop == slop", true},
		{"kino == slop", false},
		{"kino != slop", true},
		{"slop != kino", true},
	}

	for _, tt := range tests {
		eval := testEval(tt.got)
		testBooleanObject(t, eval, tt.want)
	}
}

func testBooleanObject(t *testing.T, eval object.Object, out bool) bool {
	res, ok := eval.(*object.Boolean)
	if !ok {
		t.Errorf("eval not *object.Boolean: got=%T", eval)
		return false
	}
	if res.Value != out {
		t.Errorf("res.Value not equal to %t: got=%t", out, res.Value)
		return false
	}
	return true
}

func TestNotOperator(t *testing.T) {
	tests := []struct {
		got  string
		want bool
	}{
		{"!kino", false},
		{"!slop", true},
		{"!!kino", true},
		{"!!slop", false},
		{"!1", false},
	}

	for _, tt := range tests {
		evaled := testEval(tt.got)
		testBooleanObject(t, evaled, tt.want)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		got  string
		want any
	}{
		{"hawk (1) {2};", 2},
		{"hawk (kino) {1};", 1},
		{"hawk (2 > 1) {2};", 2},
		{"hawk (slop) {1};", nil},
		{"hawk (1 > 2) {2} tuah {1};", 1},
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
		t.Errorf("object.Object not equal to NULL: got=%T", eval)
		return false
	}
	return true
}

func TestReturnValue(t *testing.T) {
	tests := []struct {
		got  string
		want int
	}{
		{"rizz 1; 2;", 1},
		{"1 * 2; rizz 2; 1;", 2},
		{"rizz 2; 2 * 1;", 2},
		{"rizz 1; rizz 2;", 1},
		{"brainrot(x, y) { x + y; }(1, 2)", 3},
	}

	for _, tt := range tests {
		evaled := testEval(tt.got)
		testIntegerObject(t, evaled, tt.want)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		got  string
		want string
	}{
		{"1 + kino;", "type mismatch: INTEGER + BOOLEAN"},
		{"-kino;", "unknown operator: -BOOLEAN"},
		{"kino + slop;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"1 - kino; 1;", "type mismatch: INTEGER - BOOLEAN"},
		{"foobar;", "unbound identifier: foobar"},
		{`"foobar" - "barfoo";`, "unknown operator: STRING - STRING"},
	}

	for _, tt := range tests {
		evaled := testEval(tt.got)
		err, ok := evaled.(*object.Error)
		if !ok {
			t.Errorf("evaled not *object.Error: got=%T", evaled)
		}

		if err.Message != tt.want {
			t.Errorf("err.Message not equal to %s: got=%s", tt.want, err.Message)
		}
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		got  string
		want int
	}{
		{"amogus a = 1; a;", 1},
		{"amogus a = 1; amogus b = a; b", 1},
		{"amogus a = 1 + 2; amogus b = a + 1; b;", 4},
	}

	for _, tt := range tests {
		evaled := testEval(tt.got)
		testIntegerObject(t, evaled, tt.want)
	}
}

func TestFunction(t *testing.T) {
	got := "brainrot(x) { x + 2; }"
	evaled := testEval(got)

	fn, ok := evaled.(*object.Function)
	if !ok {
		t.Errorf("evaled not *object.Function: got=%T", evaled)
	}

	if len(fn.Parameters) != 1 {
		t.Errorf("fn.Parameters must be 1 statement: got=%d", len(fn.Parameters))
	}
	if fn.Parameters[0].String() != "x" {
		t.Errorf("fn.Parameters[0].String not equal to x: got=%s", fn.Parameters[0].String())
	}

	if fn.Body.String() != "(x + 2)" {
		t.Errorf("fn.Body.String not equal to (x + 2): got=%s", fn.Body.String())
	}
}

func TestCallExpression(t *testing.T) {
	tests := []struct {
		got  string
		want int
	}{
		{"amogus a = brainrot(x) { x + 1; }; a(1);", 2},
		{"amogus a = brainrot(x) { x + 1; }(1); a;", 2},
		{"amogus a = brainrot(x, y) { rizz x + y; }; a(1, 2);", 3},
		{"amogus a = brainrot(x) { brainrot(y) { x + y; } }; a(1)(2);", 3},
	}

	for _, tt := range tests {
		evaled := testEval(tt.got)
		testIntegerObject(t, evaled, tt.want)
	}
}

func TestLenBuiltin(t *testing.T) {
	tests := []struct {
		got  string
		want any
	}{
		{`aura("")`, 0},
		{`aura("hello");`, 5},
		{`aura("hello world")`, 11},
		{`aura(1)`, "wrong argument: INTEGER"},
		{`aura("one", "two")`, "wrong argument number: 2"},
	}

	for _, tt := range tests {
		evaled := testEval(tt.got)
		switch want := tt.want.(type) {
		case int:
			testIntegerObject(t, evaled, want)
		case string:
			obj, ok := evaled.(*object.Error)
			if !ok {
				t.Errorf("evaled not *object.Error: got=%T", evaled)
			}
			if obj.Message != tt.want {
				t.Errorf("obj.Message not equal to %s: got=%s", tt.want, obj.Message)
			}
		}
	}
}
