package object

import "fmt"

const (
	INTEGER = "INTEGER"
	BOOLEAN = "BOOLEAN"
	NULL    = "NULL"
	RETURN  = "RETURN"
	ERROR   = "ERROR"
)

type Type string

type Object interface {
	Type() Type
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf(`%d`, i.Value)
}
func (i *Integer) Type() Type {
	return INTEGER
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf(`%t`, b.Value)
}
func (b *Boolean) Type() Type {
	return BOOLEAN
}

type Null struct{}

func (n *Null) Inspect() string {
	return `null`
}
func (n *Null) Type() Type {
	return NULL
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}
func (rv *ReturnValue) Type() Type {
	return RETURN
}

type Error struct {
	Message string
}

func (e *Error) Inspect() string {
	return `ERROR: ` + e.Message
}
func (e *Error) Type() Type {
	return ERROR
}
