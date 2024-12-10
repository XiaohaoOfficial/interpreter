package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"interpreter/ast"
	"strings"
)

const (
	INTEGER_OBJ  = `INTEGER`
	BOOLEAN_OBJ  = `BOOLEAN`
	NULL_OBJ     = `NULL`
	RETURN_OBJ   = `RETURN`
	ERROR_OBJ    = `ERROR`
	FUNCTION_OBJ = `FUNCTION`
	STRING_OBJ   = `STRING`
	BUILTIN_OBJ  = `BUILTIN`
	ARRAY_OBJ    = `ARRAY`
	HASH_OBJ     = "HASH"
)

type ObjectType string
type Object interface {
	Type() ObjectType
	Inspect() string
}
type Environment struct {
	_store map[string]Object
	_args  map[string]Object
	_outer *Environment
}

func NewEnvironment(outer *Environment) *Environment {
	res := &Environment{_store: make(map[string]Object), _outer: outer, _args: make(map[string]Object)}

	return res
}

func (e *Environment) Get(key string) (Object, bool) {
	obj, ok := e._store[key]
	if !ok && len(e._args) != 0 {
		obj, ok = e._args[key]
	}
	if !ok && e._outer != nil {
		obj, ok = e._outer.Get(key)
	}

	return obj, ok
}
func (e *Environment) Set(key string, value Object) Object {
	e._store[key] = value
	return value
}
func (e *Environment) SetArgs(key string, value Object) Object {
	e._args[key] = value
	return value
}
func (e *Environment) IsNull() bool {
	return len(e._args) == 0
}

type Interger struct {
	Value int64
}

func (i *Interger) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Interger) Type() ObjectType { return INTEGER_OBJ }

type BooleanType struct {
	Value bool
}

func (bl *BooleanType) Inspect() string  { return fmt.Sprintf("%t", bl.Value) }
func (bl *BooleanType) Type() ObjectType { return BOOLEAN_OBJ }

type NULL struct{}

func (n *NULL) Inspect() string  { return "null" }
func (n *NULL) Type() ObjectType { return NULL_OBJ }

type ReturnType struct {
	Value Object
}

func (rt *ReturnType) Inspect() string  { return rt.Value.Inspect() }
func (rt *ReturnType) Type() ObjectType { return RETURN_OBJ }

type ErrorType struct {
	Message string
}

func (et *ErrorType) Inspect() string  { return "ERROR: " + et.Message }
func (et *ErrorType) Type() ObjectType { return ERROR_OBJ }

type Function struct {
	Parameters  []*ast.Identifier
	Body        *ast.BlockStatement
	Environment *Environment
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString("){\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}
func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

type String struct {
	Value string
}

func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return STRING_OBJ }

type BuiltinFunction func(args ...Object) Object
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (ar *Array) Type() ObjectType { return ARRAY_OBJ }
func (ar *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ar.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// 由于在repl中需要打印key value，如果使用map[HashKey]Object，就会打印出HashKey,而Hash值毫无意义
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *BooleanType) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}
func (i *Interger) HashKey() HashKey {
	value := uint64(i.Value)
	return HashKey{Type: i.Type(), Value: value}
}
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}
