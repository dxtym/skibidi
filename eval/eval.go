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

func Eval(root ast.Node, env *object.Environment) object.Object {
	switch root := root.(type) {
	// statements
	case *ast.Program:
		return evalProgram(root.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(root.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatements(root.Statements, env)
	case *ast.ReturnStatement:
		val := Eval(root.Value, env)
		if checkError(val) { return val }
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(root.Value, env)
		if checkError(val) { return val }
		env.Set(root.Name.Value, val)
	// expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: root.Value}
	case *ast.Identifier:
		return evalIdentifer(root, env)
	case *ast.Boolean:
		return boolToBooleanObject(root.Value)
	case *ast.FunctionLiteral:
		args := root.Arguments
		body := root.Body
		return &object.Function{Arguments: args, Body: body, Env: env}
	case *ast.CallExpression:
		fn := Eval(root.Function, env)
		if checkError(fn) { return fn }
		args := evalExpressions(root.Arguments, env)
		if len(args) == 1 && checkError(args[0]) {
			return args[0]
		}
		// TODO: eval inside body creating inner env
	case *ast.PrefixExpression:
		right := Eval(root.Right, env)
		if checkError(right) { return right }
		return evalPrefixExpression(root.Operator, right)
	case *ast.InfixExpression:
		left := Eval(root.Left, env)
		if checkError(left) { return left }
		right := Eval(root.Right, env)
		if checkError(right) { return right }
		return evalInfixExpression(root.Operator, left, right)
	case *ast.IfElseExpression:
		return evalIfElseExpression(root, env)
	}

	return nil
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var res object.Object
	for _, node := range stmts {
		res = Eval(node, env)
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
func evalBlockStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var res object.Object
	for _, node := range stmts {
		res = Eval(node, env)
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
func evalIfElseExpression(exp *ast.IfElseExpression, env *object.Environment) object.Object {
	cond := Eval(exp.Predicate, env)
	if checkError(cond) { return cond }
	if checkTruthy(cond) {
		return Eval(exp.Consequence, env)
	} else if exp.Alternative != nil {
		return Eval(exp.Alternative, env)
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

func evalIdentifer(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("unbound indentifier: %s", node.Value)
	}
	return val
}

func evalExpressions(node []ast.Expression, env *object.Environment) []object.Object {
	res := []object.Object{}
	for _, arg := range node {
		evaled := Eval(arg, env)
		// check for error and return immediately
		if checkError(evaled) { return []object.Object{evaled} }
		res = append(res, evaled)
	}
	return res
}