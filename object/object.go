package object

import "fmt"

const (
	INTEGER = "INTEGER"
	BOOLEAN = "BOOLEAN"
	NULL    = "NULL"
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
