package object

import "fmt"

type Type string

const (
	NULL_OBJ         = "NULL"
	ERROR_OBJ        = "ERROR"
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	FUNCTION_OBJ     = "FUNCTION"
)

type Object interface {
	Type() Type
	Inspect() string
	IsType(Type) bool
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type             { return INTEGER_OBJ }
func (i *Integer) Inspect() string        { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) IsType(other Type) bool { return i.Type() == other }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type             { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string        { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) IsType(other Type) bool { return b.Type() == other }

type Null struct {
}

func (n *Null) Type() Type             { return NULL_OBJ }
func (n *Null) Inspect() string        { return "null" }
func (n *Null) IsType(other Type) bool { return n.Type() == other }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() Type             { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string        { return rv.Value.Inspect() }
func (rv *ReturnValue) IsType(other Type) bool { return rv.Type() == other }

type Error struct {
	Message string
}

func (e *Error) Type() Type             { return ERROR_OBJ }
func (e *Error) Inspect() string        { return "ERROR: " + e.Message }
func (e *Error) IsType(other Type) bool { return e.Type() == other }
