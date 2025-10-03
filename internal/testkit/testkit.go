package testkit

import (
	"net/url"
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

func NewURL(t *testing.T, s string) *url.URL {
	url, err := url.Parse(s)
	require.NoError(t, err)
	return url
}
