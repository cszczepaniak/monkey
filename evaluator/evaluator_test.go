package evaluator

import (
	"testing"

	"github.com/cszczepaniak/monkey/lexer"
	"github.com/cszczepaniak/monkey/object"
	"github.com/cszczepaniak/monkey/parser"
	"github.com/stretchr/testify/assert"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{{
		`5`, 5,
	}, {
		`10`, 10,
	}, {
		`-5`, -5,
	}, {
		`-10`, -10,
	}, {
		`5 + 5 + 5 + 5 - 10`, 10,
	}, {
		`2 * 2 * 2 * 2 * 2`, 32,
	}, {
		`-50 + 100 + -50`, 0,
	}, {
		`5 * 2 + 10`, 20,
	}, {
		`5 + 2 * 10`, 25,
	}, {
		`20 + 2 * -10`, 0,
	}, {
		`50 / 2 * 2 + 10`, 60,
	}, {
		`2 * (5 + 10)`, 30,
	}, {
		`3 * 3 * 3 + 10`, 37,
	}, {
		`3 * (3 * 3) + 10`, 37,
	}, {
		`(5 + 10 * 2 + 15 / 3) * 2 + -10`, 50,
	}}

	for _, tc := range tests {
		result := evalInput(tc.input)
		assertIntegerObject(t, result, tc.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{{
		`true`, true,
	}, {
		`false`, false,
	}, {
		`1 == 1`, true,
	}, {
		`1 != 1`, false,
	}, {
		`1 == 2`, false,
	}, {
		`1 != 2`, true,
	}, {
		`1 < 2`, true,
	}, {
		`1 > 2`, false,
	}, {
		`1 > 1`, false,
	}, {
		`1 < 1`, false,
	}, {
		`true == true`, true,
	}, {
		`false == false`, true,
	}, {
		`false == true`, false,
	}, {
		`false != true`, true,
	}, {
		`true != false`, true,
	}, {
		`(1 < 2) != false`, true,
	}, {
		`(1 < 2) == false`, false,
	}, {
		`true != (1 < 2)`, false,
	}, {
		`true == (1 < 2)`, true,
	}}

	for _, tc := range tests {
		result := evalInput(tc.input)
		assertBooleanObject(t, result, tc.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{{
		`!true`, false,
	}, {
		`!false`, true,
	}, {
		`!5`, false,
	}, {
		`!0`, false,
	}, {
		`!!true`, true,
	}, {
		`!!false`, false,
	}, {
		`!!5`, true,
	}}

	for _, tc := range tests {
		result := evalInput(tc.input)
		assertBooleanObject(t, result, tc.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{{
		`if (true) { 10 }`, 10,
	}, {
		`if (false) { 10 }`, nil,
	}, {
		`if (1) { 10 }`, 10,
	}, {
		`if (1 < 2) { 10 }`, 10,
	}, {
		`if (1 > 2) { 10 }`, nil,
	}, {
		`if (1 > 2) { 10 } else { 20 }`, 20,
	}, {
		`if (1 < 2) { 10 } else { 20 }`, 10,
	}}

	for _, tc := range tests {
		result := evalInput(tc.input)
		exp, ok := tc.expected.(int)
		if ok {
			assertIntegerObject(t, result, int64(exp))
		} else {
			assertNullObject(t, result)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{{
		`return 10;`, 10,
	}, {
		`return 10; 9;`, 10,
	}, {
		`return 2 * 5; 9;`, 10,
	}, {
		`9; return 2 * 5; 9;`, 10,
	}, {
		`
		if (10 > 1) {
			if (10 > 1) {
				return 10;
			}
			return 1;
		}`,
		10,
	}}

	for _, tc := range tests {
		result := evalInput(tc.input)
		assertIntegerObject(t, result, tc.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input  string
		expMsg string
	}{{
		`5 + true;`, `type mismatch: INTEGER + BOOLEAN`,
	}, {
		`5 + true; 5;`, `type mismatch: INTEGER + BOOLEAN`,
	}, {
		`-true`, `unknown operator: -BOOLEAN`,
	}, {
		`true + false`, `unknown operator: BOOLEAN + BOOLEAN`,
	}, {
		`5; true + false; 5;`, `unknown operator: BOOLEAN + BOOLEAN`,
	}, {
		`if (1 < 2) { true + false; }`, `unknown operator: BOOLEAN + BOOLEAN`,
	}, {
		`
		if (10 > 1) {
			if (10 > 1) {
				return true + false;
			}
			return 1;
		}`,
		`unknown operator: BOOLEAN + BOOLEAN`,
	}, {
		`foobar;`, `identifier not found: foobar`,
	}}

	for _, tc := range tests {
		result := evalInput(tc.input)
		assert.IsType(t, &object.Error{}, result)
		errObj := result.(*object.Error)
		assert.Equal(t, tc.expMsg, errObj.Message)
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{{
		`let a = 5; a;`, 5,
	}, {
		`let a = 5 * 5; a;`, 25,
	}, {
		`let a = 5; let b = a; b;`, 5,
	}, {
		`let a = 5; let b = 6 * a; b;`, 30,
	}, {
		`let a = 5; let b = 6; let c = a * b; c;`, 30,
	}}

	for _, tc := range tests {
		result := evalInput(tc.input)
		assertIntegerObject(t, result, tc.expected)
	}
}

func evalInput(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func assertIntegerObject(t *testing.T, obj object.Object, exp int64) {
	assert.IsType(t, &object.Integer{}, obj)
	integer := obj.(*object.Integer)
	assert.Equal(t, exp, integer.Value)
}

func assertBooleanObject(t *testing.T, obj object.Object, exp bool) {
	assert.IsType(t, &object.Boolean{}, obj)
	boolean := obj.(*object.Boolean)
	assert.Equal(t, exp, boolean.Value)
}

func assertNullObject(t *testing.T, obj object.Object) {
	assert.Equal(t, NULL, obj)
}
