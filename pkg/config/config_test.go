package config_test

import (
	"io"
	"log/slog"
	"testing"

	"github.com/gnames/gnfinder/pkg/config"
	"github.com/gnames/gnfinder/pkg/ent/lang"
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
