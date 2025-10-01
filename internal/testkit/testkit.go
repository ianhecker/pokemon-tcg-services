package testkit

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func ReadTestdata(t *testing.T, filename string) []byte {
	_, callerFile, _, ok := runtime.Caller(1)
	require.True(t, ok)

	dir := filepath.Dir(callerFile)
	path := filepath.Join(dir, "testdata", filename)

	bytes, err := os.ReadFile(path)
	require.NoError(t, err)

	return bytes
}
