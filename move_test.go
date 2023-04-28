package file_mover

import "testing"
import "github.com/stretchr/testify/require"

func TestFileExist(t *testing.T) {
	e := fileExists("mover.go")
	require.Equal(t, true, e)

	e = fileExists("foo.bar")
	require.Equal(t, false, e)
}
