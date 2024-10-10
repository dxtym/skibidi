package eval_test

import (
	"testing"

	"github.com/dxtym/maymun/eval"
	"github.com/dxtym/maymun/lexer"
	"github.com/dxtym/maymun/object"
	"github.com/dxtym/maymun/parser"
)

func testEval(in string) object.Object {
	l := lexer.NewLexer(in)
	p := parser.NewParser(l)
	program := p.Parse()
	return eval.Eval(program)
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		in  string
		out int
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
		eval := testEval(tt.in)
		testIntegerObject(t, eval, tt.out)
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
		in  string
		out bool
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
		eval := testEval(tt.in)
		testBooleanObject(t, eval, tt.out)
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
		in  string
		out bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!1", false},
	}

	for _, tt := range tests {
		evaled := testEval(tt.in)
		testBooleanObject(t, evaled, tt.out)
	}
}
