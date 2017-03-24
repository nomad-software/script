package evaluator

import (
	"fmt"

	"github.com/nomad-software/script/ast"
	"github.com/nomad-software/script/object"
	"github.com/nomad-software/script/token"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval evaluates the AST.
func Eval(node ast.Node, env *object.Env) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.ReturnStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Env) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value

		case *object.Error:
			return result
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL
		// return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch node := right.(type) {
	case *object.Boolean:
		if node.Value {
			return FALSE
		}
		return TRUE

	case *object.Null:
		return TRUE

	case *object.Integer:
		if node.Value == 0 {
			return TRUE
		}
	}
	return FALSE
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if !right.IsType(object.INTEGER) {
		return newError("invalid operation: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {

	if left.IsType(object.INTEGER) && right.IsType(object.INTEGER) {
		return evalIntegerInfixExpression(operator, left, right)

	} else if operator == token.EQUAL {
		return nativeBoolToBooleanObject(evalTruth(left) == evalTruth(right))

	} else if operator == token.NOT_EQUAL {
		return nativeBoolToBooleanObject(evalTruth(left) != evalTruth(right))
	}

	return newError("invalid operation: %s %s %s", left.Type(), operator, right.Type())
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case token.PLUS:
		return &object.Integer{Value: leftVal + rightVal}
	case token.MINUS:
		return &object.Integer{Value: leftVal - rightVal}
	case token.ASTERISK:
		return &object.Integer{Value: leftVal * rightVal}
	case token.SLASH:
		return &object.Integer{Value: leftVal / rightVal}
	case token.LT:
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case token.GT:
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case token.EQUAL:
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case token.NOT_EQUAL:
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("invalid operation: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalTruth(obj object.Object) *object.Boolean {
	switch obj := obj.(type) {
	case *object.Null:
		return FALSE

	case *object.Boolean:
		return obj

	case *object.Integer:
		if obj.Value != 0 {
			return TRUE
		}
		return FALSE

	default:
		return TRUE
	}
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Env) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			if result.IsType(object.RETURN_VALUE) || result.IsType(object.ERROR) {
				return result
			}
		}
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression, env *object.Env) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if evalTruth(condition).Value {
		return Eval(ie.Consequence, env)

	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)

	} else {
		return NULL
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.IsType(object.ERROR)
	}
	return false
}

func evalIdentifier(node *ast.Identifier, env *object.Env) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	// if builtin, ok := builtins[node.Value]; ok {
	// 	return builtin
	// }

	return newError("identifier not found: " + node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Env) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {

	case *object.Function:
		env := object.NewChildEnv(fn.Env)

		for i, param := range fn.Parameters {
			env.Set(param.Value, args[i])
		}

		obj := Eval(fn.Body, env)

		if r, ok := obj.(*object.ReturnValue); ok {
			return r.Value
		}

		return obj

	// case *object.Builtin:
	// 	return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}
