package config_test

import (
	"io"
	"log"
	"testing"

	"github.com/gnames/gnfinder/pkg/config"
	"github.com/gnames/gnfinder/pkg/ent/lang"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {

	log.SetOutput(io.Discard)

	t.Run("returns new Config object", func(t *testing.T) {
		cfg := config.New()
		assert.Equal(t, lang.English, cfg.Language)
		assert.Equal(t, "", cfg.LanguageDetected)
		assert.Equal(t, 0, cfg.TokensAround)
		assert.True(t, cfg.WithBayes)
		assert.False(t, cfg.WithPositionInBytes)
	})

	t.Run("takes language", func(t *testing.T) {
		cfg := config.New(config.OptLanguage(lang.English))
		assert.Equal(t, lang.English, cfg.Language)
		assert.Equal(t, "", cfg.LanguageDetected)
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
		assert.Equal(t, 4, cfg.TokensAround)
	})

	t.Run("sets find by annotation", func(t *testing.T) {
		cfg := config.New(config.OptWithFindByAnnotation(true))
		assert.Equal(t, true, cfg.WithFindByAnnotation)
	})

	t.Run("does not set 'bad' tokens number", func(t *testing.T) {
		cfg := config.New(config.OptTokensAround(-1))
		assert.Equal(t, 0, cfg.TokensAround)
		cfg = config.New(config.OptTokensAround(10))
		assert.Equal(t, 5, cfg.TokensAround)
	})

	t.Run("sets bayes' threshold", func(t *testing.T) {
		cfg := config.New(config.OptBayesOddsThreshold(200))
		assert.Equal(t, 200.0, cfg.BayesOddsThreshold)
	})

	t.Run("sets several options", func(t *testing.T) {
		opts := []config.Option{
			config.OptWithBayes(true),
			config.OptLanguage(lang.German),
		}
		cfg := config.New(opts...)
		assert.Equal(t, lang.German, cfg.Language)
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
			assert.Equal(t, v.hasErr, err != nil, v.msg)
			langOpt := config.OptLanguage(l)
			opts := []config.Option{langOpt}
			cfg := config.New(opts...)
			assert.Equal(t, v.langCfg, cfg.Language, v.msg)
		}
	})
}
