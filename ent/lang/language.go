package lang

import (
	"fmt"
	"sort"

	"github.com/abadojack/whatlanggo"
)

// Language represents the language of a text.
type Language int

// Set of languages. Last one is an indicator of the 'edge', as well as
// a default value for GnFinder.Language.
const (
	None Language = iota
	English
	German
)

var langMap = map[Language]string{
	None:    "",
	English: "eng",
	German:  "deu",
}

var langStrMap = func() map[string]Language {
	res := make(map[string]Language)
	for k, v := range langMap {
		if v != "" {
			res[v] = k
		}
	}
	return res
}()

// LanguagesSet contains supported languages and their string representation.
var LanguagesSet = func() map[Language]string {
	res := make(map[Language]string)
	for k, v := range langMap {
		if k != None {
			res[k] = v
		}
	}
	return res
}()

func (l Language) String() string {
	return langMap[l]
}

// New takes a string and returns a matching language. If string is "detect",
// returns lang.None. If language could not be found by the string, the
// function returns lang.English and an error.
func New(s string) (Language, error) {
	switch s {
	case "":
		return English, nil
	case "detect":
		return None, nil
	}

	if l, ok := langStrMap[s]; ok {
		return l, nil
	}
	return English, fmt.Errorf("unknown language %s", s)
}

func LangStrings() []string {
	res := make([]string, 0, 10)
	for k := range langStrMap {
		if k != "" {
			res = append(res, k)
		}
	}
	sort.Strings(res)
	return res
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
		return English, code
	}
}
