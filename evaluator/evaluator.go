package evaluator

import (
	"github.com/cszczepaniak/monkey/ast"
	"github.com/cszczepaniak/monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n)
	case *ast.ExpressionStatement:
		return Eval(n.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBoolObject(n.Value)
	case *ast.PrefixExpression:
		right := Eval(n.Right)
		return evalPrefixExpression(n.Operator, right)
	case *ast.InfixExpression:
		left := Eval(n.Left)
		right := Eval(n.Right)
		return evalInfixExpression(n.Operator, left, right)
	default:
		return nil
	}
}

func evalProgram(p *ast.Program) object.Object {
	var res object.Object
	for _, stmt := range p.Statements {
		res = Eval(stmt)
	}
	return res
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case `!`:
		return evalBangPrefixExpression(right)
	case `-`:
		return evalMinusPrefixExpression(right)
	default:
		return NULL
	}
}

func evalBangPrefixExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return NULL
	}
	integer := right.(*object.Integer)
	return &object.Integer{Value: -integer.Value}
}

func evalInfixExpression(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(op, left, right)
	case op == `==`:
		return nativeBoolToBoolObject(left == right)
	case op == `!=`:
		return nativeBoolToBoolObject(left != right)

	default:
		return NULL
	}
}

func evalIntegerInfixExpression(op string, left, right object.Object) object.Object {
	l, r := left.(*object.Integer).Value, right.(*object.Integer).Value
	switch op {
	case `+`:
		return &object.Integer{Value: l + r}
	case `-`:
		return &object.Integer{Value: l - r}
	case `*`:
		return &object.Integer{Value: l * r}
	case `/`:
		return &object.Integer{Value: l / r}
	case `==`:
		return nativeBoolToBoolObject(l == r)
	case `!=`:
		return nativeBoolToBoolObject(l != r)
	case `>`:
		return nativeBoolToBoolObject(l > r)
	case `<`:
		return nativeBoolToBoolObject(l < r)
	default:
		return NULL
	}
}

func nativeBoolToBoolObject(val bool) *object.Boolean {
	if val {
		return TRUE
	}
	return FALSE
}
