package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/nomad-software/script/ast"
)

type Type string

const (
	NULL         = "NULL"
	ERROR        = "ERROR"
	INTEGER      = "INTEGER"
	BOOLEAN      = "BOOLEAN"
	RETURN_VALUE = "RETURN_VALUE"
	FUNCTION     = "FUNCTION"
)

type Object interface {
	Type() Type
	Inspect() string
	IsType(Type) bool
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type             { return INTEGER }
func (i *Integer) Inspect() string        { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) IsType(other Type) bool { return i.Type() == other }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type             { return BOOLEAN }
func (b *Boolean) Inspect() string        { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) IsType(other Type) bool { return b.Type() == other }

type Null struct {
}

func (n *Null) Type() Type             { return NULL }
func (n *Null) Inspect() string        { return "null" }
func (n *Null) IsType(other Type) bool { return n.Type() == other }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() Type             { return RETURN_VALUE }
func (rv *ReturnValue) Inspect() string        { return rv.Value.Inspect() }
func (rv *ReturnValue) IsType(other Type) bool { return rv.Type() == other }

type Error struct {
	Message string
}

func (e *Error) Type() Type             { return ERROR }
func (e *Error) Inspect() string        { return "ERROR: " + e.Message }
func (e *Error) IsType(other Type) bool { return e.Type() == other }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Env
}

func (f *Function) Type() Type             { return FUNCTION }
func (f *Function) IsType(other Type) bool { return f.Type() == other }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
