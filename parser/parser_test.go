package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cszczepaniak/monkey/ast"
	"github.com/cszczepaniak/monkey/lexer"
	"github.com/stretchr/testify/assert"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expIdent string
		expValue interface{}
	}{{
		`let x = 5;`,
		`x`,
		5,
	}, {
		`let y = true;`,
		`y`,
		true,
	}, {
		`let foobar = y;`,
		`foobar`,
		`y`,
	}}

	for _, tc := range tests {
		program := assertProgram(t, tc.input, 1)
		assert.IsType(t, &ast.LetStatement{}, program.Statements[0])
		letStmt := program.Statements[0].(*ast.LetStatement)
		assertIdentifier(t, letStmt.Name, tc.expIdent)
		assertLiteralExpression(t, letStmt.Value, tc.expValue)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expValue interface{}
	}{{
		`return 5;`,
		5,
	}, {
		`return true;`,
		true,
	}, {
		`return y;`,
		`y`,
	}}

	for _, tc := range tests {
		program := assertProgram(t, tc.input, 1)
		assert.IsType(t, &ast.ReturnStatement{}, program.Statements[0])
		ret := program.Statements[0].(*ast.ReturnStatement)
		assertLiteralExpression(t, ret.ReturnValue, tc.expValue)
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`
	program := assertProgram(t, input, 1, &ast.ExpressionStatement{})
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	assert.IsType(t, &ast.Identifier{}, stmt.Expression)
	ident := stmt.Expression.(*ast.Identifier)
	assert.Equal(t, `foobar`, ident.Value)
	assert.Equal(t, `foobar`, ident.TokenLiteral())
}

func TestIntLiteralExpression(t *testing.T) {
	input := `5;`
	program := assertProgram(t, input, 1, &ast.ExpressionStatement{})
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	assertIntegerLiteral(t, stmt.Expression, 5)
}

func TestBoolLiteralExpression(t *testing.T) {
	tests := []struct {
		input  string
		expVal bool
	}{{
		`true;`,
		true,
	}, {
		`false;`,
		false,
	}}

	for _, tc := range tests {
		program := assertProgram(t, tc.input, 1, &ast.ExpressionStatement{})
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		assertLiteralExpression(t, stmt.Expression, tc.expVal)
	}
}

func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		expVal   interface{}
	}{
		{`!5;`, `!`, 5},
		{`-15;`, `-`, 15},
		{`!true;`, `!`, true},
		{`!false;`, `!`, false},
	}

	for _, tc := range tests {
		program := assertProgram(t, tc.input, 1, &ast.ExpressionStatement{})
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		assert.IsType(t, &ast.PrefixExpression{}, stmt.Expression)
		expr := stmt.Expression.(*ast.PrefixExpression)
		assert.Equal(t, tc.operator, expr.Operator)
		assertLiteralExpression(t, expr.Right, tc.expVal)
	}
}

func TestInfixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{`5 + 5`, 5, `+`, 5},
		{`5 - 5`, 5, `-`, 5},
		{`5 * 5`, 5, `*`, 5},
		{`5 / 5`, 5, `/`, 5},
		{`5 > 5`, 5, `>`, 5},
		{`5 < 5`, 5, `<`, 5},
		{`5 == 5`, 5, `==`, 5},
		{`5 != 5`, 5, `!=`, 5},
		{`true == true`, true, `==`, true},
		{`false != true`, false, `!=`, true},
		{`false == false`, false, `==`, false},
	}

	for _, tc := range tests {
		program := assertProgram(t, tc.input, 1, &ast.ExpressionStatement{})
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		assertInfixExpression(t, stmt.Expression, tc.leftValue, tc.operator, tc.rightValue)
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	program := assertProgram(t, input, 1, &ast.ExpressionStatement{})
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	assert.IsType(t, &ast.IfExpression{}, stmt.Expression)
	expr := stmt.Expression.(*ast.IfExpression)

	assertInfixExpression(t, expr.Condition, `x`, `<`, `y`)

	assert.Len(t, expr.Consequence.Statements, 1)
	consequence := expr.Consequence.Statements[0].(*ast.ExpressionStatement)
	assertIdentifier(t, consequence.Expression, `x`)
	assert.Nil(t, expr.Alternative)
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	program := assertProgram(t, input, 1, &ast.ExpressionStatement{})
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	assert.IsType(t, &ast.IfExpression{}, stmt.Expression)
	expr := stmt.Expression.(*ast.IfExpression)

	assertInfixExpression(t, expr.Condition, `x`, `<`, `y`)

	assert.Len(t, expr.Consequence.Statements, 1)
	assert.IsType(t, &ast.ExpressionStatement{}, expr.Consequence.Statements[0])
	consequence := expr.Consequence.Statements[0].(*ast.ExpressionStatement)
	assertIdentifier(t, consequence.Expression, `x`)

	assert.NotNil(t, expr.Alternative)
	assert.Len(t, expr.Alternative.Statements, 1)
	assert.IsType(t, &ast.ExpressionStatement{}, expr.Alternative.Statements[0])
	alternative := expr.Alternative.Statements[0].(*ast.ExpressionStatement)
	assertIdentifier(t, alternative.Expression, `y`)
}

func TestFuncLiteral(t *testing.T) {
	tests := []struct {
		input                string
		expArgs              []string
		expNumBodyStatements int
	}{{
		`fn() { x + y; };`,
		[]string{},
		1,
	}, {
		`fn(a) { x + y; };`,
		[]string{`a`},
		1,
	}, {
		`fn(a, b, c) { x + y; };`,
		[]string{`a`, `b`, `c`},
		1,
	}}

	for _, tc := range tests {
		program := assertProgram(t, tc.input, 1, &ast.ExpressionStatement{})

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		assert.IsType(t, &ast.FunctionLiteral{}, stmt.Expression)
		fn := stmt.Expression.(*ast.FunctionLiteral)

		assert.Len(t, fn.Args, len(tc.expArgs))
		for i, exp := range tc.expArgs {
			assertIdentifier(t, fn.Args[i], exp)
		}

		assert.Len(t, fn.Body.Statements, tc.expNumBodyStatements)
		assert.IsType(t, &ast.ExpressionStatement{}, fn.Body.Statements[0])
		bodyStmt := fn.Body.Statements[0].(*ast.ExpressionStatement)
		assertInfixExpression(t, bodyStmt.Expression, `x`, `+`, `y`)
	}
}

func TestCallExpression(t *testing.T) {
	tests := []struct {
		input       string
		expFunction string
		expArgs     []string
	}{{
		`add()`,
		`add`,
		[]string{},
	}, {
		`add(1, 2)`,
		`add`,
		[]string{`1`, `2`},
	}, {
		`add(1 + 1, 2 * 2)`,
		`add`,
		[]string{`(1 + 1)`, `(2 * 2)`},
	}, {
		`add(1 + 1, sub(2, 2 * 2))`,
		`add`,
		[]string{`(1 + 1)`, `sub(2, (2 * 2))`},
	}, {
		`fn(a, b) { a + b; }(x, y)`,
		`fn(a, b) { (a + b); }`,
		[]string{`x`, `y`},
	}}

	for _, tc := range tests {
		program := assertProgram(t, tc.input, 1, &ast.ExpressionStatement{})
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		assert.IsType(t, &ast.CallExpression{}, stmt.Expression)
		call := stmt.Expression.(*ast.CallExpression)
		assert.Equal(t, tc.expFunction, call.Function.String())
		assert.Len(t, call.Args, len(tc.expArgs))
		for i, arg := range call.Args {
			assert.Equal(t, tc.expArgs[i], arg.String())
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{{
		`a + add(b * c) + d`,
		`((a + add((b * c))) + d)`,
	}, {
		`add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))`,
		`add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))`,
	}, {
		`add(a + b + c * d / f + g)`,
		`add((((a + b) + ((c * d) / f)) + g))`,
	}, {
		`1 + (2 + 3) + 4`,
		`((1 + (2 + 3)) + 4)`,
	}, {
		`(2 + 2) * 5`,
		`((2 + 2) * 5)`,
	}, {
		`2 / (5 + 5)`,
		`(2 / (5 + 5))`,
	}, {
		`-(5 + 5)`,
		`(-(5 + 5))`,
	}, {
		`!(true == true)`,
		`(!(true == true))`,
	}, {
		`true`,
		`true`,
	}, {
		`false`,
		`false`,
	}, {
		`3 > 5 == false`,
		`((3 > 5) == false)`,
	}, {
		`3 < 5 == true`,
		`((3 < 5) == true)`,
	}, {
		`-a * b`,
		`((-a) * b)`,
	}, {
		`!-a`,
		`(!(-a))`,
	}, {
		`a + b + c`,
		`((a + b) + c)`,
	}, {
		`a + b - c`,
		`((a + b) - c)`,
	}, {
		`a * b * c`,
		`((a * b) * c)`,
	}, {
		`a * b / c`,
		`((a * b) / c)`,
	}, {
		`a + b / c`,
		`(a + (b / c))`,
	}, {
		`a + b * c + d / e - f`,
		`(((a + (b * c)) + (d / e)) - f)`,
	}, {
		`3 + 4; -5 * 5;`,
		`(3 + 4)((-5) * 5)`,
	}, {
		`5 > 4 == 3 < 4`,
		`((5 > 4) == (3 < 4))`,
	}, {
		`5 > 4 != 3 < 4`,
		`((5 > 4) != (3 < 4))`,
	}, {
		`3 + 4 * 5 == 3 * 1 + 4 * 5`,
		`((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))`,
	}}
	for _, tc := range tests {
		program := assertProgram(t, tc.input, 1)
		assert.Equal(t, tc.expected, program.String())
	}
}

func assertProgram(t *testing.T, input string, expNumStatements int, expStatementTypes ...interface{}) *ast.Program {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkErrors(t, p)

	if len(expStatementTypes) > 0 {
		assert.Len(t, program.Statements, expNumStatements)
		assert.Len(t, expStatementTypes, expNumStatements)

		for i, s := range program.Statements {
			assert.IsType(t, expStatementTypes[i], s)
		}
	}

	return program
}

func assertIntegerLiteral(t *testing.T, expr ast.Expression, expVal int64) {
	assert.IsType(t, &ast.IntegerLiteral{}, expr)
	il := expr.(*ast.IntegerLiteral)
	assert.Equal(t, expVal, il.Value)
	assert.Equal(t, fmt.Sprintf(`%d`, expVal), il.TokenLiteral())
}

func assertBooleanLiteral(t *testing.T, expr ast.Expression, expVal bool) {
	assert.IsType(t, &ast.BooleanLiteral{}, expr)
	bl := expr.(*ast.BooleanLiteral)
	assert.Equal(t, expVal, bl.Value)
	assert.Equal(t, fmt.Sprintf(`%v`, expVal), bl.TokenLiteral())
}

func assertIdentifier(t *testing.T, expr ast.Expression, expVal string) {
	assert.IsType(t, &ast.Identifier{}, expr)
	ident := expr.(*ast.Identifier)
	assert.Equal(t, expVal, ident.Value)
	assert.Equal(t, expVal, ident.TokenLiteral())
}

func assertLiteralExpression(t *testing.T, expr ast.Expression, expected interface{}) {
	switch v := expected.(type) {
	case bool:
		assertBooleanLiteral(t, expr, v)
	case int:
		assertIntegerLiteral(t, expr, int64(v))
	case int64:
		assertIntegerLiteral(t, expr, v)
	case string:
		assertIdentifier(t, expr, v)
	}
}

func assertInfixExpression(t *testing.T, expr ast.Expression, expLeft interface{}, expOp string, expRight interface{}) {
	assert.IsType(t, &ast.InfixExpression{}, expr)
	infixExpr := expr.(*ast.InfixExpression)
	assertLiteralExpression(t, infixExpr.Left, expLeft)
	assert.Equal(t, expOp, infixExpr.Operator)
	assertLiteralExpression(t, infixExpr.Right, expRight)
}

func checkErrors(t *testing.T, p *Parser) {
	assert.Emptyf(t, p.Errors(), "parser has %d errors:\n%s", len(p.Errors()), strings.Join(p.Errors(), "\n"))
}
