package eval

import (
	"fmt"

	"github.com/dxtym/maymun/object"
)

// NOTE: implicit infer of value type
var builtins = map[string]*object.Builtin{
	"uzunlik": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong argument number: %d", len(args))}
			}
			switch str := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: len(str.Value)}
			default:
				return &object.Error{Message: fmt.Sprintf("wrong argument: %s", str.Type())}
			}
		},
	},
}
