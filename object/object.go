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
	STRING       = "STRING"
	BUILTIN      = "BUILTIN"
	ARRAY        = "ARRAY"
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

type String struct {
	Value string
}

func (s *String) Type() Type             { return STRING }
func (s *String) Inspect() string        { return s.Value }
func (s *String) IsType(other Type) bool { return s.Type() == other }

// func (s *String) HashKey() HashKey {
// 	h := fnv.New64a()
// 	h.Write([]byte(s.Value))

// 	return HashKey{Type: s.Type(), Value: h.Sum64()}
// }

type Builtin struct {
	Fn func(args ...Object) Object
}

func (b *Builtin) Type() Type             { return BUILTIN }
func (b *Builtin) Inspect() string        { return "builtin function" }
func (b *Builtin) IsType(other Type) bool { return b.Type() == other }

type Array struct {
	Elements []Object
}

func (ao *Array) Type() Type             { return ARRAY }
func (ao *Array) IsType(other Type) bool { return ao.Type() == other }
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// type HashPair struct {
// 	Key   Object
// 	Value Object
// }

// type Hash struct {
// 	Pairs map[HashKey]HashPair
// }

// func (h *Hash) Type() Type { return HASH }
// func (h *Hash) Inspect() string {
// 	var out bytes.Buffer

// 	pairs := []string{}
// 	for _, pair := range h.Pairs {
// 		pairs = append(pairs, fmt.Sprintf("%s: %s",
// 			pair.Key.Inspect(), pair.Value.Inspect()))
// 	}

// 	out.WriteString("{")
// 	out.WriteString(strings.Join(pairs, ", "))
// 	out.WriteString("}")

// 	return out.String()
// }
