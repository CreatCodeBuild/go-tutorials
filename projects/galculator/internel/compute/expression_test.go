package compute

import (
	"galculator/internel/lexer"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOperatorExpression(t *testing.T) {
	exp := OperatorExpression{Left: nil, Op: lexer.Sub, Right: NumberExpression{1}}
	result, err := exp.Value()
	require.NoError(t, err)
	require.Equal(t, int64(-1), result)

	exp = OperatorExpression{Left: exp, Op: lexer.Sub, Right: NumberExpression{1}}
	result, err = exp.Value()
	require.NoError(t, err)
	require.Equal(t, int64(-2), result)
}
