package config_test

import (
	"io"
	"log/slog"
	"testing"

	"github.com/gnames/gnfinder/pkg/config"
	"github.com/gnames/gnfinder/pkg/ent/lang"
	"github.com/gnames/gnfmt"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {

	t.Run("returns new Config object", func(t *testing.T) {
		slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
		cfg := config.New()
		assert.Equal(t, cfg.Language, lang.English)
		assert.Equal(t, cfg.LanguageDetected, "")
		assert.Equal(t, cfg.TokensAround, 0)
		assert.True(t, cfg.WithBayes)
		assert.False(t, cfg.WithPositionInBytes)
	})

	t.Run("takes language", func(t *testing.T) {
		cfg := config.New(config.OptLanguage(lang.English))
		assert.Equal(t, cfg.Language, lang.English)
		assert.Equal(t, cfg.LanguageDetected, "")
	})

	t.Run("sets bayes", func(t *testing.T) {
		cfg := config.New(config.OptWithBayes(false))
		assert.False(t, cfg.WithBayes)
	})

	t.Run("sets offsets in bytes", func(t *testing.T) {
		cfg := config.New(config.OptWithPositonInBytes(true))
		assert.True(t, cfg.WithPositionInBytes)
	})

	t.Run("sets tokens number", func(t *testing.T) {
		cfg := config.New(config.OptTokensAround(4))
		assert.Equal(t, cfg.TokensAround, 4)
	})

	t.Run("does not set 'bad' tokens number", func(t *testing.T) {
		cfg := config.New(config.OptTokensAround(-1))
		assert.Equal(t, cfg.TokensAround, 0)
		cfg = config.New(config.OptTokensAround(10))
		assert.Equal(t, cfg.TokensAround, 5)
	})

	t.Run("sets bayes' threshold", func(t *testing.T) {
		cfg := config.New(config.OptBayesOddsThreshold(200))
		assert.Equal(t, cfg.BayesOddsThreshold, 200.0)
	})

	t.Run("sets several options", func(t *testing.T) {
		opts := []config.Option{
			config.OptWithBayes(true),
			config.OptLanguage(lang.German),
		}
		cfg := config.New(opts...)
		assert.Equal(t, cfg.Language, lang.German)
		assert.True(t, cfg.WithBayes)
	})

	t.Run("sets language options", func(t *testing.T) {
		tests := []struct {
			msg, lang string
			langCfg   lang.Language
			hasErr    bool
		}{
			{"default", "", lang.English, false},
			{"eng", "eng", lang.English, false},
			{"deu", "deu", lang.German, false},
			{"unknown", "notlang", lang.English, true},
			{"detect", "detect", lang.None, false},
		}

		for _, v := range tests {
			l, err := lang.New(v.lang)
			assert.Equal(t, err != nil, v.hasErr, v.msg)
			langOpt := config.OptLanguage(l)
			opts := []config.Option{langOpt}
			cfg := config.New(opts...)
			assert.Equal(t, cfg.Language, v.langCfg, v.msg)
		}
	})
}

func TestLogDefaults(t *testing.T) {
	cfg := config.New()
	assert.Equal(t, "json", cfg.Log.Format)
	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, "file", cfg.Log.Destination)
	assert.Equal(t, "", cfg.HomeDir)
}

func TestInvalidOptions(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	def := config.New()
	cfg := config.New(
		config.OptBayesOddsThreshold(-1),
		config.OptFormat(gnfmt.FormatNone),
		config.OptTikaURL(""),
		config.OptVerifierURL("  "),
		config.OptLogFormat("xml"),
		config.OptLogLevel("loud"),
		config.OptLogDestination("printer"),
		config.OptHomeDir(""),
	)
	assert.Equal(t, def, cfg)
}

func TestUpdate(t *testing.T) {
	cfg := config.New()
	cfg.Update(
		config.OptFormat(gnfmt.PrettyJSON),
		config.OptLogLevel("debug"),
		config.OptHomeDir(" /home/user "),
	)
	assert.Equal(t, gnfmt.PrettyJSON, cfg.Format)
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "/home/user", cfg.HomeDir)

	assert.False(t, cfg.WithVerification)
	cfg.Update(config.OptDataSources([]int{1, 11}))
	assert.True(t, cfg.WithVerification, "data-sources turn on verification")
}

func TestToOptions(t *testing.T) {
	t.Run("round-trip of persistent fields", func(t *testing.T) {
		src := config.New(
			config.OptBayesOddsThreshold(100),
			config.OptFormat(gnfmt.TSV),
			config.OptIncludeInputText(true),
			config.OptLanguage(lang.None),
			config.OptDataSources([]int{1, 11}),
			config.OptTikaURL("http://localhost:9998"),
			config.OptTokensAround(3),
			config.OptVerifierURL("http://localhost:8888/api/v1/"),
			config.OptWithAllMatches(true),
			config.OptWithAmbiguousNames(true),
			config.OptWithBayes(false),
			config.OptWithBayesOddsDetails(true),
			config.OptWithOddsAdjustment(true),
			config.OptWithPlainInput(true),
			config.OptWithPositonInBytes(true),
			config.OptWithUniqueNames(true),
			config.OptLogFormat("tint"),
			config.OptLogLevel("warn"),
			config.OptLogDestination("stderr"),
		)
		cfg := config.New(src.ToOptions()...)
		assert.Equal(t, src, cfg)

		src.DataSources[0] = 3
		assert.Equal(t, []int{1, 11}, cfg.DataSources, "data-sources are copied")
	})

	t.Run("runtime fields are not included", func(t *testing.T) {
		src := config.New(config.OptHomeDir("/home/user"))
		src.LanguageDetected = "deu"
		src.APIDoc = "http://example.org"
		cfg := config.New(src.ToOptions()...)
		assert.Equal(t, "", cfg.HomeDir)
		assert.Equal(t, "", cfg.LanguageDetected)
		assert.Equal(t, config.New().APIDoc, cfg.APIDoc)
	})

	t.Run("empty values do not override defaults", func(t *testing.T) {
		src := config.Config{WithBayes: true, Language: lang.English}
		cfg := config.New(src.ToOptions()...)
		assert.Equal(t, config.New(), cfg)
	})
}
