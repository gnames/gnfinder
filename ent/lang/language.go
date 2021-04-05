package lang

import (
	"fmt"

	"github.com/abadojack/whatlanggo"
)

// Language represents the language of a text.
type Language int

// Set of languages. Last one is an indicator of the 'edge', as well as
// a default value for GnFinder.Language.
const (
	DefaultLanguage Language = iota
	English
	German
	NotSet
)

func (l Language) String() string {
	languages := [...]string{"eng", "eng", "deu", ""}
	return languages[l]
}

// NewLanguage takes a string and returns matching language, or error, if
// a matching language cannot be found.
func NewLanguage(lang string) (Language, error) {
	for _, l := range SupportedLanguages() {
		if l.String() == lang {
			return l, nil
		}
	}
	var l Language
	return l, fmt.Errorf("unknown language %s", lang)
}

// SupportedLanguages returns a slice of supported by gnfinder languages.
func SupportedLanguages() []Language {
	var res []Language
	for i := 1; i < int(NotSet); i++ {
		l := Language(i)
		res = append(res, l)
	}
	return res
}

// LanguagesSet returns a 'set' of languages for more effective
// lookup of a language.
func LanguagesSet() map[Language]struct{} {
	var empty struct{}
	ls := make(map[Language]struct{})
	for i := 0; i < int(NotSet); i++ {
		ls[Language(i)] = empty
	}
	return ls
}

// DetectLanguage finds the most probable language for a text.
func DetectLanguage(text []rune) (Language, string) {
	sampleLength := len(text)
	if sampleLength > 40000 {
		sampleLength = 40000
	}
	info := whatlanggo.Detect(string(text[0:sampleLength]))
	code := whatlanggo.LangToString(info.Lang)
	switch code {
	case "eng":
		return English, code
	case "deu":
		return German, code
	default:
		return DefaultLanguage, code
	}
}
