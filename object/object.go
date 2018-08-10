package object

import (
	"fmt"
)

type ObjectType string

const (
	BOOLEAN_OBJ      = "BOOLEAN"
	INTEGER_OBJ      = "INTEGER"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
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
