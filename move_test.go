package file_mover

import (
	"fmt"
	"os"
	"testing"
	"time"
)
import "github.com/stretchr/testify/require"

func TestFileExist(t *testing.T) {
	e := fileExists("mover.go")
	require.Equal(t, true, e)

	e = fileExists("foo.bar")
	require.Equal(t, false, e)
}

func TestAll(t *testing.T) {

	err := os.RemoveAll("/tmp/src")
	require.NoError(t, err)
	err = os.RemoveAll("/tmp/dest")
	require.NoError(t, err)

	folders := map[SourceFolder]*DestinationFolder{}
	folders["/tmp/src"] = NewDestinationFolder("/tmp/dest", 1024*1024*100)

	mover := NewMover(folders)

	go func() {
		err := mover.Move()
		require.NoError(t, err)
	}()

	time.Sleep(1 * time.Second)
	for i := 0; i < 101; i++ {
		err := generateFiles(t, fmt.Sprintf("/tmp/src/%04d.jpg", i+1), 1024*1024)
		require.NoError(t, err)
	}

	time.Sleep(5 * time.Second)

}

func TestPreExistingFiles(t *testing.T) {

	err := os.RemoveAll("/tmp/src")
	require.NoError(t, err)
	err = os.RemoveAll("/tmp/dest")
	require.NoError(t, err)

	err = os.MkdirAll("/tmp/src", 0755)
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		err := generateFiles(t, fmt.Sprintf("/tmp/src/%04d.jpg", i+1), 1024)
		require.NoError(t, err)
	}

	folders := map[SourceFolder]*DestinationFolder{}
	folders["/tmp/src"] = NewDestinationFolder("/tmp/dest", 1024*11)

	mover := NewMover(folders)

	go func() {
		err := mover.Move()
		require.NoError(t, err)
	}()

	time.Sleep(1 * time.Second)

	files, err := os.ReadDir("/tmp/dest")
	require.NoError(t, err)

	require.Equal(t, 10, len(files))
}

func generateFiles(t *testing.T, file string, size int) error {
	t.Helper()
	err := os.WriteFile(file, make([]byte, size), 0644)
	if err != nil {
		return fmt.Errorf("generating file: %w", err)
	}
	return nil
}
