package compute

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompute(t *testing.T) {
	require.Equal(t, "-3", Compute("1 + 1 - 3 * 9 - 1 / 3"))
	require.Equal(t, "-8", Compute("1 + 1 - 10"))
	require.Equal(t, "2", Compute("1+1"))
	require.Equal(t, "2", Compute("(1+1)"))
	require.Equal(t, "233", Compute("233)))"))
	require.Equal(t, "1", Compute("a = 1"))            // ?
	require.Equal(t, "4", Compute("a = (1 + 3)"))      // ?
	require.Equal(t, "4", Compute("a = 1 + 3"))        // ?
	require.Equal(t, "-2", Compute("(a = 1) - 4 + a")) // ?
}

func TestComputeError(t *testing.T) {
	require.Equal(t, "Parsing Error: Missing 1 ) parentheses", Compute("(1+2"))
}
