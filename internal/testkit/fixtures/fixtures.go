package fixtures

import (
	"embed"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed testdata/*
var files embed.FS

func Read(t *testing.T, name string) []byte {
	bytes, err := fs.ReadFile(files, "testdata/"+name)
	require.NoError(t, err)

	return bytes
}
