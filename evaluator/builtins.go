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
}
