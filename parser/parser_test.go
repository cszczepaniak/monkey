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
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
		`

	program := initializeParserTest(t, input)

	assert.NotNil(t, program)
	assert.Len(t, program.Statements, 3)

	tests := []struct {
		expectedIdentifier string
	}{
		{`x`},
		{`y`},
		{`foobar`},
	}
	for i, tc := range tests {
		assert.Equal(t, `let`, program.Statements[i].TokenLiteral())
		assert.IsType(t, &ast.LetStatement{}, program.Statements[i])
		letStmt := program.Statements[i].(*ast.LetStatement)
		assertIdentifier(t, letStmt.Name, tc.expectedIdentifier)
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 239332;
		`

	program := initializeParserTest(t, input)

	assert.Len(t, program.Statements, 3)

	tests := []struct {
		expVal int64
	}{{
		5,
	}, {
		10,
	}, {
		239332,
	}}

	for i, _ := range tests {
		s := program.Statements[i]
		fmt.Print(s)
		assert.IsType(t, &ast.ReturnStatement{}, program.Statements[i])
		ret := program.Statements[i].(*ast.ReturnStatement)
		assert.Equal(t, `return`, ret.TokenLiteral())
		// assertIntegerLiteral(t, ret.ReturnValue, tc.expVal)
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`
	program := initializeParserTest(t, input)

	assert.Len(t, program.Statements, 1)
	assert.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	assert.IsType(t, &ast.Identifier{}, stmt.Expression)
	ident := stmt.Expression.(*ast.Identifier)
	assert.Equal(t, `foobar`, ident.Value)
	assert.Equal(t, `foobar`, ident.TokenLiteral())
}

func TestIntLiteralExpression(t *testing.T) {
	input := `5;`
	program := initializeParserTest(t, input)

	assert.Len(t, program.Statements, 1)
	assert.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])
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
		program := initializeParserTest(t, tc.input)
		assert.Len(t, program.Statements, 1)
		assert.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])
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
		program := initializeParserTest(t, tc.input)
		assert.Len(t, program.Statements, 1)
		assert.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])
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
		program := initializeParserTest(t, tc.input)
		assert.Len(t, program.Statements, 1)
		assert.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		assertInfixExpression(t, stmt.Expression, tc.leftValue, tc.operator, tc.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{{
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
		program := initializeParserTest(t, tc.input)
		assert.Equal(t, tc.expected, program.String())
	}
}

func initializeParserTest(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkErrors(t, p)

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
