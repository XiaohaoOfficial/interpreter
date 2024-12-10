package evaluator

import (
	"fmt"
	"interpreter/ast"
	"interpreter/object"
)

var (
	NULL  = &object.NULL{}
	TRUE  = &object.BooleanType{Value: true}
	FALSE = &object.BooleanType{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	//fmt.Appendln([]byte(node.String()))
	switch node := node.(type) {
	case *ast.HashLiteral:
		mp := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}
		for key, value := range node.Pairs {
			k := Eval(key, env)
			v := Eval(value, env)
			hp := object.HashPair{Key: k, Value: v}
			var hk object.HashKey
			switch k := k.(type) {
			case *object.Interger:
				hk = k.HashKey()
			case *object.String:
				hk = k.HashKey()
			case *object.BooleanType:
				hk = k.HashKey()
			}
			mp.Pairs[hk] = hp
		}
		return mp
	case *ast.IndexExpression:
		switch left := Eval(node.Left, env).(type) {
		case *object.Array:
			el := Eval(node.Index, env)
			idx, ok := el.(*object.Interger)
			if !ok {
				return &object.ErrorType{Message: fmt.Sprintf("index:%s is not INTEGER", node.Index.TokenLiteral())}
			}
			if int(idx.Value) >= len(left.Elements) || idx.Value < 0 {
				return NULL
			}
			return left.Elements[idx.Value]
		case *object.Hash:
			el := Eval(node.Index, env)
			var key object.HashKey
			switch el := el.(type) {
			case *object.BooleanType:
				key = el.HashKey()
			case *object.String:
				key = el.HashKey()
			case *object.Interger:
				key = el.HashKey()
			}
			value, ok := left.Pairs[key]
			if !ok {
				return NULL
			}
			return value.Value
		}

	case *ast.ArrayLiteral:
		var elements []object.Object
		for _, ele := range node.Elements {
			elements = append(elements, Eval(ele, env))
		}
		return &object.Array{Elements: elements}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.CallExpression:
		switch Callfn := Eval(node.Function, env).(type) {
		case *object.Function:
			if node.Function.String() != node.Function.TokenLiteral() {
				Newenv := object.NewEnvironment(env)
				fn := Eval(node.Function, Newenv)
				fun, _ := fn.(*object.Function)
				if len(fun.Parameters) != len(node.Arguments) {
					return &object.ErrorType{Message: fmt.Sprintf("want %d Arguments get=%d", len(fun.Parameters), len(node.Arguments))}
				}
				for i, arg := range node.Arguments {
					env.Set(fun.Parameters[i].Value, Eval(arg, env))
				}
				res := Eval(fun.Body, Newenv)
				if rt, ok := res.(*object.ReturnType); ok {
					return rt.Value
				}
				return res
			}
			fn, _ := env.Get(node.Function.TokenLiteral())
			fun := fn.(*object.Function)
			if len(fun.Parameters) != len(node.Arguments) {
				return &object.ErrorType{Message: fmt.Sprintf("want %d Arguments get=%d", len(fun.Parameters), len(node.Arguments))}
			}
			var fnEnv *object.Environment
			if !fun.Environment.IsNull() {
				fnEnv = fun.Environment
			} else {
				fnEnv = object.NewEnvironment(env)
			}

			for i, arg := range node.Arguments {
				fnEnv.Set(fun.Parameters[i].Value, Eval(arg, env))
			}
			res := Eval(fun.Body, fnEnv)
			if rt, ok := res.(*object.ReturnType); ok {
				return rt.Value
			}
			return res
		case *object.Builtin:

			var args = []object.Object{}
			for i := range node.Arguments {
				arg := Eval(node.Arguments[i], env)
				args = append(args, arg)
			}
			return Callfn.Fn(args...)
		}
	case *ast.FunctionLiteral:
		res := &object.Function{Environment: env}
		res.Parameters = node.Parameters
		res.Body = node.Body
		return res
	case *ast.LetStatement:
		name := node.Name.Value
		value := Eval(node.Value, env)
		if value.Type() == object.FUNCTION_OBJ {
			value := value.(*object.Function)
			if para, ok := node.Value.(*ast.CallExpression); ok {
				gt, _ := env.Get(para.Function.TokenLiteral())
				fn := gt.(*object.Function)
				for i, pr := range para.Arguments {
					value.Environment.SetArgs(fn.Parameters[i].Value, Eval(pr, env))
				}
			}
		}
		env.Set(name, value)
		return NULL
	case *ast.ReturnStatement:
		value := Eval(node.ReturnValue, env)
		if value.Type() == object.ERROR_OBJ {
			return value
		}
		return &object.ReturnType{Value: value}
	case *ast.Program:
		return evalStatements(node.Statements, env)
	case *ast.IntegerLiteral:
		return &object.Interger{Value: node.Value}
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.Boolean:
		return nativeBooleanObject(node.Value)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node.Operator, Eval(node.Right, env))
	case *ast.Identifier:
		value, ok := env.Get(node.Value)
		if ok {
			return value
		}
		if builtin, ok := builtins[node.Value]; ok {
			return builtin
		}
		return &object.ErrorType{Message: fmt.Sprintf("identifier not found: %s", node.Value)}

	case *ast.InfixExpression:
		if node.Operator == "=" {
			if _, ok := node.Left.(*ast.Identifier); !ok {
				return &object.ErrorType{Message: fmt.Sprintf("unknown Assign for %s", node.Left.TokenLiteral())}
			}
			if _, ok := env.Get(node.Left.TokenLiteral()); !ok {
				return &object.ErrorType{Message: fmt.Sprintf("identifier not found: %s", node.Left.TokenLiteral())}
			}
			env.Set(node.Left.TokenLiteral(), Eval(node.Right, env))
			return Eval(node.Left, env)
		}
		left := Eval(node.Left, env)

		if left.Type() == object.ERROR_OBJ {
			return left
		}
		right := Eval(node.Right, env)
		if right.Type() == object.ERROR_OBJ {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalStatements(node.Statements, env)
	case *ast.IfExpression:
		cond := Eval(node.Condition, env)
		if cond != NULL && cond != FALSE {
			return Eval(node.Consequence, env)
		}
		if node.Alternative == nil {
			return NULL
		}
		return Eval(node.Alternative, env)
	}
	return NULL
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		if right.Type() != object.INTEGER_OBJ {
			return &object.ErrorType{Message: fmt.Sprintf("unknown operator: %s%s", operator, right.Type())}
		}
		return evalMinusOperatorExpression(right)
	default:
		return NULL
	}
}
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}
func evalMinusOperatorExpression(right object.Object) object.Object {
	rt, ok := right.(*object.Interger)
	if !ok {
		return NULL
	}
	return &object.Interger{Value: (-rt.Value)}
}
func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt, env)
		if result.Type() == object.RETURN_OBJ || result.Type() == object.ERROR_OBJ {
			return result
		}
	}

	return result
}

func nativeBooleanObject(input bool) *object.BooleanType {
	if input {
		return TRUE
	}
	return FALSE
}
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		if operator == "==" {
			if left == right {
				return TRUE
			}
			return FALSE
		}
		if operator == "!=" {
			if left != right {
				return TRUE
			}
			return FALSE
		}
		return &object.ErrorType{Message: fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, right.Type())}
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		if operator == "+" {
			return evalStringInfixExpression(operator, left, right)
		}
		return &object.ErrorType{Message: fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, right.Type())}
	default:
		return &object.ErrorType{Message: fmt.Sprintf("type mismatch: %s %s %s", left.Type(), operator, right.Type())}
	}
}
func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	l := left.(*object.String)
	r := right.(*object.String)
	switch operator {
	case "+":
		return &object.String{Value: l.Value + r.Value}
	}
	return &object.ErrorType{Message: fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, right.Type())}

}
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	l := left.(*object.Interger)
	r := right.(*object.Interger)
	switch operator {
	case "+":
		return &object.Interger{Value: l.Value + r.Value}
	case "-":
		return &object.Interger{Value: l.Value - r.Value}
	case "*":
		return &object.Interger{Value: l.Value * r.Value}
	case "/":
		return &object.Interger{Value: l.Value / r.Value}
	case ">":
		return returnBool(l.Value > r.Value)
	case "<":
		return returnBool(l.Value < r.Value)
	case "==":
		return returnBool(l.Value == r.Value)
	case "!=":
		return returnBool(l.Value != r.Value)
	case "<=":
		return returnBool(l.Value <= r.Value)
	case ">=":
		return returnBool(l.Value >= r.Value)
	default:
		return NULL
	}
}
func returnBool(value bool) *object.BooleanType {
	if value {
		return TRUE
	}
	return FALSE
}
