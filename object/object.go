package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/dxtym/skibidi/ast"
)

type ObjectType string
type BuiltinFunction func(args ...Object) Object

const (
	INTEGER_OBJECT    = "INTEGER"
	STRING_OBJECT     = "STRING"
	BOOLEAN_OBJECT    = "BOOLEAN"
	NULL_OBJECT       = "NULL"
	RETURN_VAL_OBJECT = "RETURN_VAL"
	ERROR_OBJECT      = "ERROR"
	FUNCTION_OBJECT   = "FUNCTION"
	BUILTIN_OBJECT    = "BUILTIN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJECT }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJECT }
func (s *String) Inspect() string  { return s.Value }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJECT }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJECT }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VAL_OBJECT }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// TODO: add stack trace from extra fields of token
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJECT }
func (e *Error) Inspect() string  { return fmt.Sprintf("%s", e.Message) }

type Environment struct {
	store map[string]Object
	other *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
		other: nil,
	}
}

func NewEnclosedEnvironment(other *Environment) *Environment {
	env := NewEnvironment()
	env.other = other
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	val, ok := e.store[name]
	if !ok && e.other != nil {
		val, ok = e.other.Get(name)
	}
	return val, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJECT }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	args := []string{}
	for _, arg := range f.Parameters {
		args = append(args, arg.String())
	}

	out.WriteString("func(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJECT }
func (b *Builtin) Inspect() string  { return "builtin function" }
