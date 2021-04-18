package lexer

import (
	"testing"

	"github.com/cszczepaniak/monkey/token"
	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for _, tc := range tests {
		tok := l.NextToken()

		assert.Equal(t, tc.expectedType, tok.Type)
		assert.Equal(t, tc.expectedLiteral, tok.Literal)
	}
}
