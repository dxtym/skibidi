package eval

import (
	"fmt"

	"github.com/dxtym/skibidi/object"
)

// NOTE: implicit infer of value type
var builtins = map[string]*object.Builtin{
	"yap": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("bruh: %d", len(args))}
			}

			switch obj := args[0].(type) {
			case *object.String:
				return &object.String{Value: obj.Value}
			case *object.Integer:
				return &object.Integer{Value: obj.Value}
			case *object.Array:
				return &object.Array{Elements: obj.Elements}
			default:
				return &object.Error{Message: fmt.Sprintf("bruh: %s", obj.Type())}
			}
		},
	},
	"aura": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("bruh: %d", len(args))}
			}

			switch obj := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: len(obj.Value)}
			case *object.Array:
				return &object.Integer{Value: len(obj.Elements)}
			default:
				return &object.Error{Message: fmt.Sprintf("bruh: %s", obj.Type())}
			}
		},
	},
	"chad": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("bruh: %d", len(args))}
			}
			if args[0].Type() != object.ARRAY_OBJECT {
				return &object.Error{Message: fmt.Sprintf("bruh: %s", args[0])}
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"skuf": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("bruh: %d", len(args))}
			}
			if args[0].Type() != object.ARRAY_OBJECT {
				return &object.Error{Message: fmt.Sprintf("bruh: %s", args[0])}
			}

			arr := args[0].(*object.Array)
			len := len(arr.Elements)
			if len > 0 {
				return arr.Elements[len-1]
			}

			return NULL
		},
	},
	"fam": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("bruh: %d", len(args))}
			}
			if args[0].Type() != object.ARRAY_OBJECT {
				return &object.Error{Message: fmt.Sprintf("bruh: %s", args[0])}
			}

			arr := args[0].(*object.Array)
			len := len(arr.Elements)
			if len > 0 {
				arr2 := make([]object.Object, len-1)
				copy(arr2, arr.Elements[1:])
				return &object.Array{Elements: arr2}
			}

			return NULL
		},
	},
	"yeet": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("bruh: %d", len(args))}
			}
			if args[0].Type() != object.ARRAY_OBJECT {
				return &object.Error{Message: fmt.Sprintf("bruh: %s", args[0])}
			}

			arr := args[0].(*object.Array)
			len := len(arr.Elements)
			arr2 := make([]object.Object, len+1)
			copy(arr2, arr.Elements)
			arr2[len] = args[1]

			return &object.Array{Elements: arr2}
		},
	},
}
