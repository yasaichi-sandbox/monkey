package evaluator

import (
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
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}

	return nil
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

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
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

	return NULL // NOTE: エラー処理を実装したらNULLを返す以外の実装になるかもしれない
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	// NOTE: 副作用を与えた後に最後の評価結果を返す、がやりたいのならもっと良い書き方がある気がする
	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}
