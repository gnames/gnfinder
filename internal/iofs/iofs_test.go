package iofs_test

import (
	"os"
	"testing"

	"github.com/gnames/gnfinder/internal/iofs"
	"github.com/gnames/gnfinder/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureDirs(t *testing.T) {
	home := t.TempDir()
	require.NoError(t, iofs.EnsureDirs(home))
	for _, dir := range []string{
		config.ConfigDir(home),
		config.CacheDir(home),
		config.LogDir(home),
	} {
		info, err := os.Stat(dir)
		require.NoError(t, err)
		assert.True(t, info.IsDir())
	}
	// second call is a no-op
	assert.NoError(t, iofs.EnsureDirs(home))
}

func TestEnsureConfigFile(t *testing.T) {
	home := t.TempDir()
	require.NoError(t, iofs.EnsureDirs(home))
	require.NoError(t, iofs.EnsureConfigFile(home))

	path := config.ConfigFilePath(home)
	bs, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, iofs.ConfigYAML, string(bs))

	// existing config file is not overwritten
	require.NoError(t, os.WriteFile(path, []byte("format: tsv\n"), 0644))
	require.NoError(t, iofs.EnsureConfigFile(home))
	bs, err = os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "format: tsv\n", string(bs))
}
