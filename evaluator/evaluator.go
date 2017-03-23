package evaluator

import (
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
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		// if isError(right) {
		// 	return right
		// }
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		// if isError(left) {
		// 	return left
		// }

		right := Eval(node.Right)
		// if isError(right) {
		// 	return right
		// }

		return evalInfixExpression(node.Operator, left, right)
	}

	return nil
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement)

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
	if !right.IsType(object.INTEGER_OBJ) {
		return NULL
		// return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {

	if left.IsType(object.INTEGER_OBJ) && right.IsType(object.INTEGER_OBJ) {
		return evalIntegerInfixExpression(operator, left, right)

	} else if operator == token.EQUAL {
		return nativeBoolToBooleanObject(evalTruth(left) == evalTruth(right))

	} else if operator == token.NOT_EQUAL {
		return nativeBoolToBooleanObject(evalTruth(left) != evalTruth(right))
	}

	return NULL
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
		return NULL
		// 	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalTruth(obj object.Object) *object.Boolean {
	switch obj.Type() {

	case object.NULL_OBJ:
		return FALSE

	case object.BOOLEAN_OBJ:
		return obj.(*object.Boolean)

	case object.INTEGER_OBJ:
		if obj.(*object.Integer).Value != 0 {
			return TRUE
		}
		return FALSE

	default:
		return FALSE
	}
}
