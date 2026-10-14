package iologger_test

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/gnfinder/internal/iologger"
	"github.com/gnames/gnfinder/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitFile(t *testing.T) {
	defer slog.SetDefault(slog.Default())
	dir := t.TempDir()
	path := filepath.Join(dir, iologger.LogFile)
	cfg := config.LogConfig{Format: "json", Level: "info", Destination: "file"}

	require.NoError(t, iologger.Init(dir, cfg, false))
	slog.Info("first")
	slog.Debug("hidden")

	require.NoError(t, iologger.Init(dir, cfg, true))
	slog.Info("second")

	bs, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(bs), `"msg":"first"`)
	assert.Contains(t, string(bs), `"msg":"second"`)
	assert.NotContains(t, string(bs), "hidden")

	// without append the log file is truncated
	require.NoError(t, iologger.Init(dir, cfg, false))
	slog.Info("third")
	bs, err = os.ReadFile(path)
	require.NoError(t, err)
	assert.NotContains(t, string(bs), "first")
	assert.Contains(t, string(bs), "third")
}

func TestInitTintFile(t *testing.T) {
	defer slog.SetDefault(slog.Default())
	dir := t.TempDir()
	cfg := config.LogConfig{Format: "tint", Level: "debug", Destination: "file"}

	require.NoError(t, iologger.Init(dir, cfg, false))
	slog.Debug("tinted")

	bs, err := os.ReadFile(filepath.Join(dir, iologger.LogFile))
	require.NoError(t, err)
	assert.Contains(t, string(bs), "tinted")
	assert.NotContains(t, string(bs), "\x1b[", "no color codes in log file")
}

func TestInitBadDir(t *testing.T) {
	cfg := config.LogConfig{Format: "json", Level: "info", Destination: "file"}
	err := iologger.Init(filepath.Join(t.TempDir(), "missing"), cfg, false)
	assert.Error(t, err)
}
