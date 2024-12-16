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
				return &object.Error{Message: fmt.Sprintf("cant yap: %d", len(args))}
			}

			switch obj := args[0].(type) {
			case *object.String:
				return &object.String{Value: obj.Value}
			case *object.Integer:
				return &object.Integer{Value: obj.Value}
			case *object.Array:
				return &object.Array{Elements: obj.Elements}
			default:
				return &object.Error{Message: fmt.Sprintf("got mogged: %v", obj.Type())}
			}
		},
	},
	"aura": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("zero aura: %d", len(args))}
			}

			switch obj := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: len(obj.Value)}
			case *object.Array:
				return &object.Integer{Value: len(obj.Elements)}
			default:
				return &object.Error{Message: fmt.Sprintf("got mogged: %s", obj.Type())}
			}
		},
	},
}
