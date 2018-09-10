package object

import (
	"bytes"
	"fmt"
	"github.com/yasaichi-sandbox/monkey/ast"
	"hash/fnv"
	"strings"
)

type BuiltinFunction func(args ...Object) Object

type HashKey struct {
	Type  ObjectType
	Value uint64
}

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
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

type Hashable interface {
	HashKey() HashKey
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Array struct {
	Elements []Object
}

func (a *Array) Inspect() string {
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}
func (*Array) Type() ObjectType { return ARRAY_OBJ }

type Boolean struct {
	Value bool
}

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
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

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Inspect() string {
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, pair.Key.Inspect()+":"+pair.Value.Inspect())
	}

	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}
func (*Hash) Type() ObjectType { return HASH_OBJ }

type Integer struct {
	Value int64
}

func (i *Integer) HashKey() HashKey {
	// NOTE: `i.Value`が表現できる負の値（int64なので2**63-1まで）は`uint64()`とすると
	// int64では表現できない2**63 ~ 2**64-1にマッピングされる（＝int64で表現できる正の値と
	// 被ることはない）ので、ハッシュ値を計算する用途においてはこれでも大丈夫。なんだけども、、
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
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

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
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
