package config_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/gnames/gnfinder/config"
	"github.com/gnames/gnfinder/ent/lang"
	"github.com/tj/assert"
)

func TestConfig(t *testing.T) {

	log.SetOutput(ioutil.Discard)

	t.Run("returns new Config object", func(t *testing.T) {
		cfg := config.New()
		assert.Equal(t, cfg.Language, lang.DefaultLanguage)
		assert.Equal(t, cfg.LanguageDetected, "")
		assert.Equal(t, cfg.TokensAround, 0)
		assert.True(t, cfg.WithBayes)
    assert.False(t, cfg.WithBytesOffset)
	})

	t.Run("takes language", func(t *testing.T) {
		cfg := config.New(config.OptLanguage(lang.English))
		assert.Equal(t, cfg.Language, lang.English)
		assert.False(t, cfg.WithLanguageDetection)
		assert.Equal(t, cfg.LanguageDetected, "")
	})

	t.Run("sets bayes", func(t *testing.T) {
		cfg := config.New(config.OptWithBayes(false))
		assert.False(t, cfg.WithBayes)
	})

	t.Run("sets offsets in bytes", func(t *testing.T) {
		cfg := config.New(config.OptWithBytesOffset(true))
		assert.True(t, cfg.WithBytesOffset)
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
		assert.False(t, cfg.WithLanguageDetection)
		assert.True(t, cfg.WithBayes)
	})
}
