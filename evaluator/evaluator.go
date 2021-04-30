package evaluator

import (
	"fmt"

	"github.com/cszczepaniak/monkey/ast"
	"github.com/cszczepaniak/monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n, env)
	case *ast.BlockStatement:
		return evalBlockStatement(n, env)
	case *ast.ExpressionStatement:
		return Eval(n.Expression, env)
	case *ast.IfExpression:
		return evalIfExpression(n, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(n, env)
	case *ast.LetStatement:
		val := Eval(n.Value, env)
		if val.Type() == object.ERROR {
			return val
		}
		return env.Set(n.Name.Value, val)
	case *ast.Identifier:
		val, ok := env.Get(n.Value)
		if !ok {
			return newErrorf(`identifier not found: %s`, n.Value)
		}
		return val
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBoolObject(n.Value)
	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if right.Type() == object.ERROR {
			return right
		}
		return evalPrefixExpression(n.Operator, right)
	case *ast.InfixExpression:
		left := Eval(n.Left, env)
		if left.Type() == object.ERROR {
			return left
		}
		right := Eval(n.Right, env)
		if right.Type() == object.ERROR {
			return right
		}
		return evalInfixExpression(n.Operator, left, right)
	default:
		return nil
	}
}

func evalProgram(p *ast.Program, env *object.Environment) object.Object {
	var res object.Object
	for _, stmt := range p.Statements {
		res = Eval(stmt, env)

		switch r := res.(type) {
		case *object.ReturnValue:
			return r.Value
		case *object.Error:
			return r
		}
	}
	return res
}

func evalBlockStatement(bs *ast.BlockStatement, env *object.Environment) object.Object {
	var res object.Object
	for _, stmt := range bs.Statements {
		res = Eval(stmt, env)
		if res != nil && (res.Type() == object.RETURN || res.Type() == object.ERROR) {
			return res
		}
	}
	return res
}

func evalIfExpression(is *ast.IfExpression, env *object.Environment) object.Object {
	c := Eval(is.Condition, env)
	if c.Type() == object.ERROR {
		return c
	}
	if c == NULL || c == FALSE {
		if is.Alternative != nil {
			return evalBlockStatement(is.Alternative, env)
		}
		return NULL
	}
	return evalBlockStatement(is.Consequence, env)
}

func evalReturnStatement(rs *ast.ReturnStatement, env *object.Environment) object.Object {
	res := Eval(rs.ReturnValue, env)
	if res.Type() == object.ERROR {
		return res
	}
	return &object.ReturnValue{Value: res}
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case `!`:
		return evalBangPrefixExpression(right)
	case `-`:
		return evalMinusPrefixExpression(right)
	default:
		return newErrorf(`unknown operator: %s%s`, op, right.Type())
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
		return newErrorf(`unknown operator: -%s`, right.Type())
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
	case left.Type() != right.Type():
		return newErrorf(`type mismatch: %s %s %s`, left.Type(), op, right.Type())
	default:
		return newErrorf(`unknown operator: %s %s %s`, left.Type(), op, right.Type())
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
		return newErrorf(`unknown operator: %s %s %s`, left.Type(), op, right.Type())
	}
}

func nativeBoolToBoolObject(val bool) *object.Boolean {
	if val {
		return TRUE
	}
	return FALSE
}

func newErrorf(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
