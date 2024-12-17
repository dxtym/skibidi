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

func testEval(got string) object.Object {
	l := lexer.NewLexer(got)
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
		{"fax", true},
		{"cap", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"fax == fax", true},
		{"cap == cap", true},
		{"fax == cap", false},
		{"fax != cap", true},
		{"cap != fax", true},
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
		{"!fax", false},
		{"!cap", true},
		{"!!fax", true},
		{"!!cap", false},
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
		{"hawk (fax) {1};", 1},
		{"hawk (2 > 1) {2};", 2},
		{"hawk (cap) {1};", nil},
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
		{"cook(x, y) { x + y; }(1, 2)", 3},
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
		{"1 + fax;", "touch grass: INTEGER + BOOLEAN"},
		{"-fax;", "baka: -BOOLEAN"},
		{"fax + cap;", "delulu: BOOLEAN + BOOLEAN"},
		{"1 - fax; 1;", "touch grass: INTEGER - BOOLEAN"},
		{"foobar;", "delulu: foobar"},
		{`"foobar" - "barfoo";`, "touch grass: STRING - STRING"},
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
	got := "cook(x) { x + 2; }"
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
		{"amogus a = cook(x) { x + 1; }; a(1);", 2},
		{"amogus a = cook(x) { x + 1; }(1); a;", 2},
		{"amogus a = cook(x, y) { rizz x + y; }; a(1, 2);", 3},
		{"amogus a = cook(x) { cook(y) { x + y; } }; a(1)(2);", 3},
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
		{`yap("hello")`, "hello"},
		{`yap(123)`, 123},
		{`yap([1, 2, 3])`, []int{1, 2, 3}},
		{`aura("")`, 0},
		{`aura("hello world")`, 11},
	}

	for _, tt := range tests {
		evaled := testEval(tt.got)
		switch want := tt.want.(type) {
		case int:
			testIntegerObject(t, evaled, want)
		case string:
			obj, ok := evaled.(*object.String)
			if !ok {
				t.Errorf("evaled not *object.Error: got=%T", evaled)
			}
			if obj.Value != tt.want {
				t.Errorf("obj.Value not equal to %s: got=%s", tt.want, obj.Value)
			}
		}
	}
}

func TestArrayLiteral(t *testing.T) {
	got := "[1, 2, 3, 4];"
	evaled := testEval(got)
	obj, ok := evaled.(*object.Array)
	if !ok {
		t.Errorf("evaled not *object.Array: got=%T", evaled)
	}

	testIntegerObject(t, obj.Elements[0], 1)
	testIntegerObject(t, obj.Elements[1], 2)
	testIntegerObject(t, obj.Elements[2], 3)
	testIntegerObject(t, obj.Elements[3], 4)
}

func TestArrayIndexExpression(t *testing.T) {
	tests := []struct {
		got  string
		want any
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
	}

	for _, tt := range tests {
		evaled := testEval(tt.got)
		val, ok := tt.want.(int)
		if ok {
			testIntegerObject(t, evaled, val)
		} else {
			testNullObject(t, evaled)
		}
	}
}

func TestMapLiteral(t *testing.T) {
	got := `{"foo": 1, 2: 2, fax: 3}`
	evaled := testEval(got)
	obj, ok := evaled.(*object.Map)
	if !ok {
		t.Errorf("evaled not *object.Map: got=%T", evaled)
	}

	want := map[object.Hash]int{
		(&object.String{Value: "foo"}).Hash(): 1,
		(&object.Integer{Value: 2}).Hash():    2,
		TRUE.Hash():                           3,
	}

	if len(obj.Pairs) != len(want) {
		t.Errorf("obj.Pairs must be 3 statments: got=%T", len(obj.Pairs))
	}

	for key, val := range want {
		pair, ok := obj.Pairs[key]
		if !ok {
			t.Errorf("no pair for %v", key)
		}

		testIntegerObject(t, pair.Value, val)
	}
}

func TestMapIndexExpression(t *testing.T) {
	tests := []struct {
		got  string
		want any
	}{
		{`{"foo": 1}["foo"]`, 1},
		{`{"foo": 1}["bar"]`, nil},
		{`{}["foo"]`, nil},
		{`{1: 2}[1]`, 2},
		{`{fax: 1}[fax]`, 1},
		{`{cap: 1}[cap]`, 1},
	}

	for _, tt := range tests {
		evaled := testEval(tt.got)
		integer, ok := tt.want.(int)
		if ok {
			testIntegerObject(t, evaled, integer)
		} else {
			testNullObject(t, evaled)
		}
	}
}
