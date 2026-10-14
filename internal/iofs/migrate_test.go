package iofs_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gnames/gnfinder/internal/iofs"
	"github.com/gnames/gnfinder/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// setupLegacy creates a home with XDG directories and a legacy config file.
func setupLegacy(t *testing.T, content string) (string, string) {
	t.Helper()
	home := t.TempDir()
	require.NoError(t, iofs.EnsureDirs(home))
	legacy := filepath.Join(home, ".config", iofs.LegacyConfigFile)
	require.NoError(t, os.WriteFile(legacy, []byte(content), 0644))
	return home, legacy
}

func readMigrated(t *testing.T, home string) (string, map[string]any) {
	t.Helper()
	bs, err := os.ReadFile(config.ConfigFilePath(home))
	require.NoError(t, err)
	var res map[string]any
	require.NoError(t, yaml.Unmarshal(bs, &res))
	return string(bs), res
}

func TestMigrateLegacyConfig(t *testing.T) {
	home, legacy := setupLegacy(t, `
# comment
BayesOddsThreshold: 100
DataSources: [1, 11]
Format: pretty
Language: detect
TikaURL: http://localhost:9998
TokensAround: 2
WithUniqueNames: true
WithVerification: false
WithoutBayes: true
`)
	keys, err := iofs.MigrateLegacyConfig(legacy, home)
	require.NoError(t, err)
	assert.Equal(t, []string{
		"bayes_odds_threshold", "data_sources", "format", "language",
		"tika_url", "tokens_around", "with_bayes", "with_unique_names",
	}, keys)

	txt, res := readMigrated(t, home)
	assert.True(t, strings.HasPrefix(txt, iofs.ConfigYAML))
	assert.Contains(t, txt, "# Settings migrated from "+legacy)
	assert.Equal(t, 100, res["bayes_odds_threshold"])
	assert.Equal(t, []any{1, 11}, res["data_sources"])
	assert.Equal(t, "pretty", res["format"])
	assert.Equal(t, "detect", res["language"])
	assert.Equal(t, "http://localhost:9998", res["tika_url"])
	assert.Equal(t, 2, res["tokens_around"])
	assert.Equal(t, false, res["with_bayes"])
	assert.Equal(t, true, res["with_unique_names"])
	assert.NotContains(t, res, "with_verification", "false values are skipped")

	bs, err := os.ReadFile(legacy)
	require.NoError(t, err)
	assert.Contains(t, string(bs), "WithoutBayes: true", "legacy file untouched")
}

func TestMigrateLegacyConfigEmpty(t *testing.T) {
	// legacy file created by gnfinder v1.1 has only comments
	home, legacy := setupLegacy(t, "# Format: csv\n# WithoutBayes: false\n")
	keys, err := iofs.MigrateLegacyConfig(legacy, home)
	require.NoError(t, err)
	assert.Empty(t, keys)
	_, err = os.Stat(config.ConfigFilePath(home))
	assert.True(t, os.IsNotExist(err), "config.yaml is not created")
}

func TestMigrateLegacyConfigInvalid(t *testing.T) {
	home, legacy := setupLegacy(t, `
format: xml
language: klingon
tokensaround: -1
datasources: [a]
withverification: maybe
unknownkey: 1
withuniquenames: "true"
`)
	keys, err := iofs.MigrateLegacyConfig(legacy, home)
	require.NoError(t, err)
	assert.Equal(t, []string{"with_unique_names"}, keys)
}

func TestMigrateLegacyConfigTokensAround(t *testing.T) {
	home, legacy := setupLegacy(t, "TokensAround: 7\n")
	_, err := iofs.MigrateLegacyConfig(legacy, home)
	require.NoError(t, err)
	_, res := readMigrated(t, home)
	assert.Equal(t, 5, res["tokens_around"])
}

func TestMigrateLegacyConfigAliases(t *testing.T) {
	home, legacy := setupLegacy(t, "PreferredSources: \"4, 12\"\n")
	_, err := iofs.MigrateLegacyConfig(legacy, home)
	require.NoError(t, err)
	_, res := readMigrated(t, home)
	assert.Equal(t, []any{4, 12}, res["data_sources"])

	home, legacy = setupLegacy(t, "PreferredSources: [4]\nDataSources: [1]\n")
	_, err = iofs.MigrateLegacyConfig(legacy, home)
	require.NoError(t, err)
	_, res = readMigrated(t, home)
	assert.Equal(t, []any{1}, res["data_sources"], "DataSources wins")
}

func TestMigrateLegacyConfigBadYAML(t *testing.T) {
	home, legacy := setupLegacy(t, "Format: [csv\n")
	_, err := iofs.MigrateLegacyConfig(legacy, home)
	assert.Error(t, err)
}
