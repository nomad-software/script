package evaluator

import (
	"github.com/nomad-software/script/ast"
	"github.com/nomad-software/script/object"
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
		if node.Value {
			return TRUE
		}
		return FALSE

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		// if isError(right) {
		// 	return right
		// }
		return evalPrefixExpression(node.Operator, right)

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
