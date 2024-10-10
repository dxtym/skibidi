package eval

import (
	"github.com/dxtym/maymun/ast"
	"github.com/dxtym/maymun/object"
)

// define for all time usage
var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(root ast.Node) object.Object {
	switch root := root.(type) {
	// statements
	case *ast.Program:
		return evalProgram(root.Statements)
	case *ast.ExpressionStatement:
		return Eval(root.Expression)
	// values
	case *ast.IntegerLiteral:
		return &object.Integer{Value: root.Value}
	case *ast.Boolean:
		return boolToBooleanObject(root.Value)
	case *ast.PrefixExpression:
		right := Eval(root.Right)
		return evalPrefixExpression(root.Operator, right)
	case *ast.InfixExpression:
		left := Eval(root.Left)
		right := Eval(root.Right)
		return evalInfixExression(root.Operator, left, right)
	}

	return nil
}

func evalProgram(stmts []ast.Statement) object.Object {
	var res object.Object
	for _, node := range stmts {
		res = Eval(node)
	}
	return res
}

func boolToBooleanObject(in bool) object.Object {
	if in {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalNotOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	}

	return NULL // TODO: better user exp
}

func evalNotOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJECT {
		return NULL
	}

	val := right.(*object.Integer).Value
	return &object.Integer{Value: -val}
}

func evalInfixExression(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJECT && right.Type() == object.INTEGER_OBJECT:
		return evalIntegerInfixExpression(op, left, right)
	// pointer comparison (works because of TRUE and FALSE)
	case op == "!=":
		return boolToBooleanObject(left != right)
	case op == "==":
		return boolToBooleanObject(left == right)
	default:
		return NULL
	}
}

func evalIntegerInfixExpression(op string, left, right object.Object) object.Object {
	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value
	switch op {
	case "+":
		return &object.Integer{Value: l + r}
	case "-":
		return &object.Integer{Value: l - r}
	case "*":
		return &object.Integer{Value: l * r}
	case "/":
		return &object.Integer{Value: l / r}
	case "<":
		return boolToBooleanObject(l < r)
	case ">":
		return boolToBooleanObject(l > r)
	case "!=":
		return boolToBooleanObject(l != r)
	case "==":
		return boolToBooleanObject(l == r)
	default:
		return NULL
	}
}
