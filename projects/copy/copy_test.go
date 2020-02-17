package copy

import "testing"
import "github.com/stretchr/testify/require"

func TestSerial(t *testing.T) {
	err := Serial("./testdata/go", "testdata/go_copy")
	require.NoError(t, err)
}

func TestConcurrent(t *testing.T) {
	err := Concurrent("./testdata/go", "testdata/go_copy_concurrent")
	require.NoError(t, err)
}