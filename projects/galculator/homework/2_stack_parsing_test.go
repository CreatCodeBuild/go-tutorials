package homework

import "testing"
import "github.com/stretchr/testify/require"

func TestTwoStack(t *testing.T) {
	numbers, operators := TwoStack("1+1-2")
	require.Equal(t, []string{"1", "1", "2"}, numbers)
	require.Equal(t, []string{"+", "-"}, numbers)

	numbers, operators = TwoStack("-1+1-2")
	require.Equal(t, []string{"-1", "1", "2"}, numbers)
	require.Equal(t, []string{"+", "-"}, numbers)
}
