package evaluator

import (
	"fmt"
	"github.com/yasaichi-sandbox/monkey/ast"
	"github.com/yasaichi-sandbox/monkey/object"
)

// NOTE: Goの定数では、構造体を除く値型しか定義できないので`var`を使っている、はず。
// たぶんコンパイル時に評価して値をスタック領域に詰めないからな気がする、たぶん。
var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.LetStatement:
		val := Eval(node.Value)
		if isError(val) {
			return val
		}

		// TODO: いい感じにする
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}

		right := Eval(node.Right)
		if isError(right) {
			return left
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}

	return nil
}

func isTruthy(obj object.Object) bool {
	return obj != FALSE && obj != NULL
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	}

	// NOTE: truthy values
	return FALSE
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)
		if result == nil { // NOTE: この`nil`チェックは近いうちに消される気がする
			continue
		}

		switch result.Type() {
		case object.ERROR_OBJ, object.RETURN_VALUE_OBJ:
			return result
		}
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative == nil {
		return NULL
	} else {
		return Eval(ie.Alternative)
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	// NOTE: ここから下に入るパターンは次の通り
	// * leftとrightのどちらかがInteger
	// => ==/!=以外の演算を異なるType同士で試みているので、type mismatch
	// * leftとrightのどちらもIntegerでない = monkeyの場合、BooleanかNull
	//   * BooleanとNull  => （同上の理由で）type mismatch
	//   * どちらもBoolean => ==/!=以外は未定義なので、unknown operator
	//   * どちらもNull    => ==/!=以外は未定義なので、unknown operator
	case left.Type() != right.Type():
		return newError(
			"type mismatch: %s %s %s",
			left.Type(), operator, right.Type(),
		)
	}

	return newError(
		"unknown operator: %s %s %s",
		left.Type(), operator, right.Type(),
	)
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	}

	return newError(
		"unknown operator: %s %s %s",
		left.Type(), operator, right.Type(),
	)
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	}

	return newError("unknown operator: %s%s", operator, right.Type())
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement)

		switch result := result.(type) {
		case *object.Error:
			return result
		case *object.ReturnValue:
			return result.Value
		}
	}

	return result
}

func isError(obj object.Object) bool {
	// NOTE: interfaceに対するnilチェックをせずにメソッド呼び出しを書きがちなので気をつけたい
	if obj == nil {
		return false
	}

	return obj.Type() == object.ERROR_OBJ
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
