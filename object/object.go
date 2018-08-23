package object

import (
	"bytes"
	"fmt"
	"github.com/yasaichi-sandbox/monkey/ast"
	"strings"
)

type BuiltinFunction func(args ...Object) Object
type ObjectType string

const (
	BOOLEAN_OBJ      = "BOOLEAN"
	INTEGER_OBJ      = "INTEGER"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
func (*Boolean) Type() ObjectType  { return BOOLEAN_OBJ }

type Builtin struct {
	Fn BuiltinFunction
}

func (*Builtin) Inspect() string  { return "builtin function" }
func (*Builtin) Type() ObjectType { return BUILTIN_OBJ }

type Error struct {
	Message string
}

func (e *Error) Inspect() string { return "ERROR: " + e.Message }
func (*Error) Type() ObjectType  { return ERROR_OBJ }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Inspect() string {
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	buf := &bytes.Buffer{}
	fmt.Fprintf(
		buf,
		"fn(%s) {\n%s\n}",
		strings.Join(params, ", "),
		f.Body.String(),
	)

	return buf.String()
}
func (*Function) Type() ObjectType { return FUNCTION_OBJ }

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (*Integer) Type() ObjectType  { return INTEGER_OBJ }

type Null struct {
}

func (*Null) Inspect() string  { return "null" }
func (*Null) Type() ObjectType { return NULL_OBJ }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }
func (*ReturnValue) Type() ObjectType   { return RETURN_VALUE_OBJ }

type String struct {
	Value string
}

func (s *String) Inspect() string { return s.Value }
func (*String) Type() ObjectType  { return STRING_OBJ }

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

func NewEnvironment() *Environment {
	return &Environment{store: map[string]Object{}}
}

func (e *Environment) Get(name string) (Object, bool) {
	if obj, ok := e.store[name]; ok || e.outer == nil {
		return obj, ok
	}

	// NOTE: Find value recursively
	obj, ok := e.outer.Get(name)
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
