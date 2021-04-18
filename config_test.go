package gnfinder_test

import (
	"testing"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/ent/lang"
	"github.com/tj/assert"
)

func TestConfig(t *testing.T) {

	t.Run("returns new Config object", func(t *testing.T) {
		cfg := gnfinder.NewConfig()
		assert.Equal(t, cfg.Language, lang.DefaultLanguage)
		assert.Equal(t, cfg.LanguageDetected, "")
		assert.Equal(t, cfg.TokensAround, 0)
		assert.True(t, cfg.WithBayes)
	})

	t.Run("takes language", func(t *testing.T) {
		cfg := gnfinder.NewConfig(gnfinder.OptLanguage(lang.English))
		assert.Equal(t, cfg.Language, lang.English)
		assert.False(t, cfg.WithLanguageDetection)
		assert.Equal(t, cfg.LanguageDetected, "")
	})

	t.Run("sets bayes", func(t *testing.T) {
		cfg := gnfinder.NewConfig(gnfinder.OptWithBayes(false))
		assert.False(t, cfg.WithBayes)
	})

	t.Run("sets tokens number", func(t *testing.T) {
		cfg := gnfinder.NewConfig(gnfinder.OptTokensAround(4))
		assert.Equal(t, cfg.TokensAround, 4)
	})

	t.Run("does not set 'bad' tokens number", func(t *testing.T) {
		cfg := gnfinder.NewConfig(gnfinder.OptTokensAround(-1))
		assert.Equal(t, cfg.TokensAround, 0)
		cfg = gnfinder.NewConfig(gnfinder.OptTokensAround(10))
		assert.Equal(t, cfg.TokensAround, 5)
	})

	t.Run("sets bayes' threshold", func(t *testing.T) {
		cfg := gnfinder.NewConfig(gnfinder.OptBayesThreshold(200))
		assert.Equal(t, cfg.BayesOddsThreshold, 200.0)
	})

	t.Run("sets several options", func(t *testing.T) {
		opts := []gnfinder.Option{
			gnfinder.OptWithBayes(true),
			gnfinder.OptLanguage(lang.German),
		}
		cfg := gnfinder.NewConfig(opts...)
		assert.Equal(t, cfg.Language, lang.German)
		assert.False(t, cfg.WithLanguageDetection)
		assert.True(t, cfg.WithBayes)
	})
}
