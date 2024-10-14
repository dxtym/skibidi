package eval

import (
	"fmt"

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
	case *ast.BlockStatement:
		return evalBlockStatements(root.Statements)
	case *ast.ReturnStatement:
		val := Eval(root.Value)
		if checkError(val) { return val }
		return &object.ReturnValue{Value: val}
	// values
	case *ast.IntegerLiteral:
		return &object.Integer{Value: root.Value}
	case *ast.Boolean:
		return boolToBooleanObject(root.Value)
	case *ast.PrefixExpression:
		right := Eval(root.Right)
		if checkError(right) { return right }
		return evalPrefixExpression(root.Operator, right)
	case *ast.InfixExpression:
		left := Eval(root.Left)
		if checkError(left) { return left }
		right := Eval(root.Right)
		if checkError(right) { return right }
		return evalInfixExpression(root.Operator, left, right)
	case *ast.IfElseExpression:
		return evalIfElseExpression(root)
	}

	return nil
}

func evalProgram(stmts []ast.Statement) object.Object {
	var res object.Object
	for _, node := range stmts {
		res = Eval(node)
		switch res := res.(type) {
		case *object.ReturnValue:
			return res.Value
		case *object.Error:
			return res
		}
	}
	return res
}

// NOTE: 
// to avoid return outermost value in nested
// statements. not unwrap value but check its type
func evalBlockStatements(stmts []ast.Statement) object.Object {
	var res object.Object
	for _, node := range stmts {
		res = Eval(node)
		if res != nil {
			rt := res.Type()
			if rt == object.ERROR_OBJECT || rt == object.RETURN_VAL_OBJECT {
				return res
			}
		}
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
	default:
		return newError("unknown operator: %s%s", op, right.Type())
	}
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
		return newError("unknown operator: -%s", right.Type())
	}

	val := right.(*object.Integer).Value
	return &object.Integer{Value: -val}
}

func evalInfixExpression(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJECT && right.Type() == object.INTEGER_OBJECT:
		return evalIntegerInfixExpression(op, left, right)
	// pointer comparison (works because of TRUE and FALSE)
	case op == "!=":
		return boolToBooleanObject(left != right)
	case op == "==":
		return boolToBooleanObject(left == right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), op, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
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
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

// TODO: revise complicated logic
func evalIfElseExpression(exp *ast.IfElseExpression) object.Object {
	cond := Eval(exp.Predicate)
	if checkError(cond) { return cond }
	if checkTruthy(cond) {
		return Eval(exp.Consequence)
	} else if exp.Alternative != nil {
		return Eval(exp.Alternative)
	} else {
		return NULL
	}
}

func checkTruthy(cond object.Object) bool {
	switch cond {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL:
		return false
	default:
		return true
	}
}

func newError(format string, a ...any) object.Object {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// to avoid errors being passed around 
func checkError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJECT
	}
	return false
}