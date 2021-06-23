package lang_test

import (
	"testing"

	"github.com/gnames/gnfinder/ent/lang"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert.Equal(t, lang.English.String(), "eng")
	assert.Equal(t, lang.German.String(), "deu")
	assert.Equal(t, lang.None.String(), "")
}

func TestLanguageSet(t *testing.T) {
	ls := lang.LanguagesSet
	_, ok := ls[lang.English]
	assert.True(t, ok)
	_, ok = ls[lang.German]
	assert.True(t, ok)
}

func TestDetectLang(t *testing.T) {
	text := `
          should be permitted to remain ; and this should be trained up, with a
          single stem, to the utmost height of its growth, and never stop'd or
          cut back. The horizontal branches or head will then be found to form
          itself, by pushing out shoots immediately around the point of the
          year's per- pendicular shoot or stem ; and as this will be long or
          short, according to the soil and situation, the horizontal tiers of
          branches will be at pro- portional and proper distances ; and thus
          the tree will assume the shape and growth of the fir or the wild
          cherry-tree. If any irregular shoots should push out on the sides of
          the stem, or too many horizontals, they may be removed. And if the
          perpendicular stem or leading shoot should be destroyed, one of the
          horizontals may be fixed
					`
	l, code := lang.DetectLanguage([]rune(text))
	assert.Equal(t, l, lang.English)
	assert.Equal(t, code, "eng")

	// unknown language detection
	text = "Однажды в студеную, зимнюю пору я из лесу вышел"
	l, code = lang.DetectLanguage([]rune(text))
	assert.Equal(t, l, lang.English)
	assert.Equal(t, code, "rus")
}
