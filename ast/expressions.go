package ast

import (
	"bytes"

	"github.com/cszczepaniak/monkey/token"
)

type Expression interface {
	Node
	expressionNode()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
	return `(` + pe.Operator + pe.Right.String() + `)`
}

type InfixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
	Left     Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *InfixExpression) String() string {
	return `(` + ie.Left.String() + ` ` + ie.Operator + ` ` + ie.Right.String() + `)`
}

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) TokenLiteral() string {
	return bl.Token.Literal
}
func (bl *BooleanLiteral) String() string {
	return bl.Token.Literal
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString(`if`)
	out.WriteString(ie.Condition.String())
	out.WriteString(` `)
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString(`else `)
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token token.Token
	Args  []*Identifier
	Body  *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	out.WriteString(`fn(`)
	for i, arg := range fl.Args {
		out.WriteString(arg.String())
		if i < len(fl.Args)-1 {
			out.WriteString(`, `)
		}
	}
	out.WriteString(`) `)
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token    token.Token
	Function Expression
	Args     []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ce.Function.String())
	out.WriteString(`(`)
	for i, arg := range ce.Args {
		out.WriteString(arg.String())
		if i < len(ce.Args)-1 {
			out.WriteString(`, `)
		}
	}
	out.WriteString(`)`)

	return out.String()
}
