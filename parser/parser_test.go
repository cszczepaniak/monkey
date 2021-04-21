package parser

import (
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
		assert.Equal(t, tc.expectedIdentifier, letStmt.Name.Value)
		assert.Equal(t, tc.expectedIdentifier, letStmt.Name.TokenLiteral())
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

	for _, s := range program.Statements {
		assert.IsType(t, &ast.ReturnStatement{}, s)
		assert.Equal(t, `return`, s.TokenLiteral())
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
	assert.IsType(t, &ast.IntegerLiteral{}, stmt.Expression)
	expr := stmt.Expression.(*ast.IntegerLiteral)
	assert.Equal(t, int64(5), expr.Value)
	assert.Equal(t, `5`, expr.TokenLiteral())
}

func initializeParserTest(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkErrors(t, p)

	return program
}

func checkErrors(t *testing.T, p *Parser) {
	assert.Emptyf(t, p.Errors(), "parser has %d errors:\n%s", len(p.Errors()), strings.Join(p.Errors(), "\n"))
}
