package cmd

import (
	"os"
	"testing"

	"github.com/gnames/gnfinder/internal/iofs"
	"github.com/gnames/gnfinder/pkg/config"
	"github.com/gnames/gnfinder/pkg/ent/lang"
	"github.com/gnames/gnfmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupHome creates a home directory with config.yaml containing content.
// Empty content uses the default embedded config.yaml.
func setupHome(t *testing.T, content string) string {
	t.Helper()
	for _, env := range envVars {
		t.Setenv(env, "")
		require.NoError(t, os.Unsetenv(env))
	}
	t.Setenv(envWithoutBayes, "")
	require.NoError(t, os.Unsetenv(envWithoutBayes))

	home := t.TempDir()
	require.NoError(t, iofs.EnsureDirs(home))
	if content == "" {
		content = iofs.ConfigYAML
	}
	err := os.WriteFile(config.ConfigFilePath(home), []byte(content), 0644)
	require.NoError(t, err)
	return home
}

func configFrom(t *testing.T, home string) config.Config {
	t.Helper()
	opts, err := initConfig(home)
	require.NoError(t, err)
	return config.New(opts...)
}

func TestInitConfigDefaults(t *testing.T) {
	home := setupHome(t, "")
	assert.Equal(t, config.New(), configFrom(t, home))
}

func TestInitConfigFile(t *testing.T) {
	home := setupHome(t, `
format: pretty
language: detect
bayes_odds_threshold: 100
tokens_around: 2
data_sources: [1, 11]
with_bayes: false
with_unique_names: true
log:
  level: debug
  format: tint
`)
	cfg := configFrom(t, home)
	assert.Equal(t, gnfmt.PrettyJSON, cfg.Format)
	assert.Equal(t, lang.None, cfg.Language)
	assert.Equal(t, 100.0, cfg.BayesOddsThreshold)
	assert.Equal(t, 2, cfg.TokensAround)
	assert.Equal(t, []int{1, 11}, cfg.DataSources)
	assert.True(t, cfg.WithVerification, "data-sources turn on verification")
	assert.False(t, cfg.WithBayes)
	assert.True(t, cfg.WithUniqueNames)
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "tint", cfg.Log.Format)
	assert.Equal(t, "file", cfg.Log.Destination, "missing keys keep defaults")
	assert.Equal(t, config.New().TikaURL, cfg.TikaURL)
}

func TestInitConfigEnv(t *testing.T) {
	home := setupHome(t, "format: pretty\n")
	t.Setenv("GNF_FORMAT", "tsv")
	t.Setenv("GNF_LANGUAGE", "deu")
	t.Setenv("GNF_DATA_SOURCES", "1,12")
	t.Setenv("GNF_WITH_AMBIGUOUS_NAMES", "true")
	t.Setenv("GNF_TIKA_URL", "http://localhost:9998")
	t.Setenv("GNF_LOG_DESTINATION", "stderr")

	cfg := configFrom(t, home)
	assert.Equal(t, gnfmt.TSV, cfg.Format, "env overrides config file")
	assert.Equal(t, lang.German, cfg.Language)
	assert.Equal(t, []int{1, 12}, cfg.DataSources)
	assert.True(t, cfg.WithAmbiguousNames)
	assert.Equal(t, "http://localhost:9998", cfg.TikaURL)
	assert.Equal(t, "stderr", cfg.Log.Destination)
}

func TestInitConfigWithoutBayesEnv(t *testing.T) {
	home := setupHome(t, "")
	t.Setenv(envWithoutBayes, "true")
	assert.False(t, configFrom(t, home).WithBayes)
}

func TestInitConfigInvalid(t *testing.T) {
	tests := []struct {
		msg, content string
	}{
		{"format", "format: xml\n"},
		{"language", "language: klingon\n"},
		{"number", "tokens_around: many\n"},
		{"yaml", "format: [csv\n"},
	}
	for _, v := range tests {
		home := setupHome(t, v.content)
		_, err := initConfig(home)
		assert.Error(t, err, v.msg)
	}
}

func TestInitConfigInvalidValueIgnored(t *testing.T) {
	home := setupHome(t, "tika_url: ''\nlog:\n  level: loud\n")
	cfg := configFrom(t, home)
	assert.Equal(t, config.New().TikaURL, cfg.TikaURL)
	assert.Equal(t, "info", cfg.Log.Level)
}

func newTestCmd(t *testing.T, args ...string) *cobra.Command {
	t.Helper()
	cmd := &cobra.Command{}
	cmd.Flags().AddFlagSet(rootCmd.Flags())
	// AddFlagSet shares flags with rootCmd, reset their state.
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
	require.NoError(t, cmd.Flags().Parse(args))
	return cmd
}

func TestFlagOpts(t *testing.T) {
	t.Run("no flags keep config values", func(t *testing.T) {
		cfg := config.New(
			config.OptFormat(gnfmt.PrettyJSON),
			config.OptLanguage(lang.German),
			config.OptWithBayes(false),
		)
		cfg.Update(flagOpts(newTestCmd(t))...)
		assert.Equal(t, gnfmt.PrettyJSON, cfg.Format)
		assert.Equal(t, lang.German, cfg.Language)
		assert.False(t, cfg.WithBayes)
	})

	t.Run("flags override config values", func(t *testing.T) {
		cfg := config.New(config.OptFormat(gnfmt.PrettyJSON))
		cmd := newTestCmd(t,
			"-f", "csv", "-l", "detect", "-n", "-u", "-w", "3",
			"-s", "1, 11", "-t", "http://localhost:9998",
		)
		cfg.Update(flagOpts(cmd)...)
		assert.Equal(t, gnfmt.CSV, cfg.Format)
		assert.Equal(t, lang.None, cfg.Language)
		assert.False(t, cfg.WithBayes)
		assert.True(t, cfg.WithUniqueNames)
		assert.Equal(t, 3, cfg.TokensAround)
		assert.Equal(t, []int{1, 11}, cfg.DataSources)
		assert.True(t, cfg.WithVerification)
		assert.Equal(t, "http://localhost:9998", cfg.TikaURL)
	})

	t.Run("all matches turn on verification", func(t *testing.T) {
		cfg := config.New()
		cfg.Update(flagOpts(newTestCmd(t, "-M"))...)
		assert.True(t, cfg.WithAllMatches)
		assert.True(t, cfg.WithVerification)
	})

	t.Run("bad values are ignored", func(t *testing.T) {
		cfg := config.New(config.OptFormat(gnfmt.TSV))
		cfg.Update(flagOpts(newTestCmd(t, "-f", "xml", "-l", "klingon"))...)
		assert.Equal(t, gnfmt.TSV, cfg.Format)
		assert.Equal(t, lang.English, cfg.Language)
	})
}

func TestParseDataSources(t *testing.T) {
	assert.Equal(t, []int{1, 11}, parseDataSources("1, 11"))
	assert.Equal(t, []int{3}, parseDataSources("0,a,3"))
	assert.Empty(t, parseDataSources("a"))
}

func TestInitConfigBadEnvList(t *testing.T) {
	home := setupHome(t, "")
	t.Setenv("GNF_DATA_SOURCES", "1,x")
	_, err := initConfig(home)
	assert.Error(t, err)
}
