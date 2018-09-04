package evaluator

import "github.com/yasaichi-sandbox/monkey/object"

var builtins = map[string]*object.Builtin{
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError(
					"argument to `first` must be ARRAY, got %s",
					args[0].Type(),
				)
			}

			array := args[0].(*object.Array)
			if len(array.Elements) == 0 {
				return NULL
			}

			return array.Elements[0]
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError(
					"argument to `last` must be ARRAY, got %s",
					args[0].Type(),
				)
			}

			array := args[0].(*object.Array)
			length := len(array.Elements)
			if length == 0 {
				return NULL
			}

			return array.Elements[length-1]
		},
	},
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			// NOTE: これ、nilが渡ってきたときにどうするんだ、、
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			}

			return newError(
				"argument to `len` not supported, got %s",
				args[0].Type(),
			)
		},
	},
	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError(
					"argument to `rest` must be ARRAY, got %s",
					args[0].Type(),
				)
			}

			array := args[0].(*object.Array)
			length := len(array.Elements)
			if length == 0 {
				return NULL
			}

			// 第2引数がlength, 第3引数がcapacity。capを省略するとlenと同じ値が設定される
			newElements := make([]object.Object, length-1)
			// NOTE: この場合ポインタの配列なので、`copy`でdeep copyする必要があるのか謎だった
			copy(newElements, array.Elements[1:])

			return &object.Array{Elements: newElements}
		},
	},
}
