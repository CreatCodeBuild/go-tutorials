package compute

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompute(t *testing.T) {
	require.Equal(t, "-3", Compute("1 + 1 - 3 * 9 - 1 / 3", nil))
	require.Equal(t, "-8", Compute("1 + 1 - 10", nil))
	require.Equal(t, "2", Compute("1+1", nil))
	require.Equal(t, "2", Compute("(1+1)", nil))
	require.Equal(t, "233", Compute("233)))", nil))
	require.Equal(t, "1", Compute("a = 1", nil))            // ?
	require.Equal(t, "4", Compute("a = (1 + 3)", nil))      // ?
	require.Equal(t, "4", Compute("a = 1 + 3", nil))        // ?
	require.Equal(t, "-2", Compute("(a = 1) - 4 + a", nil)) // ?
	require.Equal(t, "3", Compute("(a=1)+a+1", nil))
	require.Equal(t, "-8", Compute("(a=1)-10+a", nil))
}

func TestComputeError(t *testing.T) {
	require.Equal(t, "Parsing Error: Missing 1 ) parentheses", Compute("(1+2", nil))
}
