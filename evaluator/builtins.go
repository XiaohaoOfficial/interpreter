package evaluator

import (
	"fmt"
	"interpreter/object"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.ErrorType{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Interger{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Interger{Value: int64(len(arg.Elements))}
			default:
				return &object.ErrorType{Message: fmt.Sprintf("argument to `len` not supported, got %s", arg.Type())}
			}
		},
	},
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.ErrorType{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[0]
				}
				return NULL
			default:
				return &object.ErrorType{Message: fmt.Sprintf("argument to `len` not supported, got %s", arg.Type())}
			}
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.ErrorType{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[len(arg.Elements)-1]
				}
				return NULL
			default:
				return &object.ErrorType{Message: fmt.Sprintf("argument to `len` not supported, got %s", arg.Type())}
			}
		},
	},
	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.ErrorType{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					arr := make([]object.Object, len(arg.Elements)-1)
					copy(arr, arg.Elements[1:len(arg.Elements)])
					return &object.Array{Elements: arr}
				}
				return NULL
			default:
				return &object.ErrorType{Message: fmt.Sprintf("argument to `len` not supported, got %s", arg.Type())}
			}
		},
	},
	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.ErrorType{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(args))}
			}
			switch arg := args[0].(type) {
			case *object.Array:

				arr := make([]object.Object, len(arg.Elements)+1)
				copy(arr, arg.Elements[0:len(arg.Elements)])
				arr[len(arg.Elements)] = args[1]
				return &object.Array{Elements: arr}

			default:
				return &object.ErrorType{Message: fmt.Sprintf("argument to `len` not supported, got %s", arg.Type())}
			}
		},
	},
	"puts": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
}
