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

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkErrors(t, p)
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

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkErrors(t, p)
	assert.Len(t, program.Statements, 3)

	for _, s := range program.Statements {
		assert.IsType(t, &ast.ReturnStatement{}, s)
		assert.Equal(t, `return`, s.TokenLiteral())
	}
}

func checkErrors(t *testing.T, p *Parser) {
	assert.Emptyf(t, p.Errors(), "parser has %d errors:\n%s", len(p.Errors()), strings.Join(p.Errors(), "\n"))
}
