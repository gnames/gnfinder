package iofs

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gnames/gn"
	"github.com/gnames/gnfinder/pkg/config"
	"github.com/gnames/gnfinder/pkg/ent/lang"
	"github.com/gnames/gnfmt"
	"gopkg.in/yaml.v3"
)

// LegacyConfigFile is the name of the config file used by gnfinder v1.1
// and earlier. It was located directly in the user's config directory.
const LegacyConfigFile = "gnfinder.yml"

// legacyKeys maps lowercased keys of the legacy config file to keys of
// config.yaml.
var legacyKeys = map[string]string{
	"bayesoddsthreshold":   "bayes_odds_threshold",
	"datasources":          "data_sources",
	"preferredsources":     "data_sources",
	"format":               "format",
	"includeinputtext":     "include_input_text",
	"inputtextonly":        "input_text_only",
	"language":             "language",
	"tikaurl":              "tika_url",
	"tokensaround":         "tokens_around",
	"verifierurl":          "verifier_url",
	"withallmatches":       "with_all_matches",
	"withambiguousnames":   "with_ambiguous_names",
	"withbayesoddsdetails": "with_bayes_odds_details",
	"withoddsadjustment":   "with_odds_adjustment",
	"withplaininput":       "with_plain_input",
	"withpositioninbytes":  "with_position_in_bytes",
	"withuniquenames":      "with_unique_names",
	"withverification":     "with_verification",
	"withoutbayes":         "with_bayes",
}

// MigrateLegacyConfig creates config.yaml with settings found in the legacy
// gnfinder.yml file. It returns the keys of migrated settings. If the legacy
// file has no settings, config.yaml is not created and no keys are returned.
//
// Like the legacy CLI, it only migrates values that differ from zero values:
// true booleans, non-empty strings and lists, positive numbers. Invalid
// values are skipped with a warning. The legacy file is not modified.
func MigrateLegacyConfig(legacyPath, homeDir string) ([]string, error) {
	bs, err := os.ReadFile(legacyPath)
	if err != nil {
		return nil, ReadFileError(legacyPath, err)
	}

	var legacy map[string]any
	if err = yaml.Unmarshal(bs, &legacy); err != nil {
		return nil, ReadFileError(legacyPath, err)
	}

	settings := convertLegacy(legacy)
	if len(settings) == 0 {
		return nil, nil
	}

	out, err := yaml.Marshal(settings)
	if err != nil {
		return nil, err
	}

	content := ConfigYAML + fmt.Sprintf(
		"\n# Settings migrated from %s on %s\n",
		legacyPath, time.Now().Format(time.DateOnly),
	) + string(out)

	configPath := config.ConfigFilePath(homeDir)
	if err = os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return nil, CopyFileError(configPath, err)
	}

	keys := make([]string, 0, len(settings))
	for k := range settings {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys, nil
}

func convertLegacy(legacy map[string]any) map[string]any {
	res := make(map[string]any)
	for k, v := range legacy {
		lk := strings.ToLower(k)
		key, ok := legacyKeys[lk]
		if !ok {
			gn.Warn("Unknown setting <em>%s</em> in old config file, skipping", k)
			continue
		}
		if v == nil {
			continue
		}

		val, ok := convertValue(lk, v)
		if !ok {
			gn.Warn("Invalid value <em>%v</em> of <em>%s</em> in old config file, skipping", v, k)
			continue
		}
		if val == nil {
			continue
		}
		// DataSources has a priority over its old alias PreferredSources.
		if key == "data_sources" && lk == "preferredsources" {
			if _, ok := res[key]; ok {
				continue
			}
		}
		res[key] = val
	}
	return res
}

// convertValue returns the converted value and true. If the value has the
// zero value, it returns nil and true, if it is invalid, it returns false.
func convertValue(key string, v any) (any, bool) {
	switch key {
	case "format":
		s, ok := toString(v)
		if !ok {
			return nil, false
		}
		if s == "" {
			return nil, true
		}
		if _, err := gnfmt.NewFormat(s); err != nil {
			return nil, false
		}
		return s, true
	case "language":
		s, ok := toString(v)
		if !ok {
			return nil, false
		}
		if s == "" {
			return nil, true
		}
		if _, err := lang.New(s); err != nil {
			return nil, false
		}
		return s, true
	case "tikaurl", "verifierurl":
		s, ok := toString(v)
		if !ok {
			return nil, false
		}
		if s == "" {
			return nil, true
		}
		return s, true
	case "bayesoddsthreshold":
		f, ok := toFloat(v)
		if !ok || f < 0 {
			return nil, false
		}
		if f == 0 {
			return nil, true
		}
		return f, true
	case "tokensaround":
		f, ok := toFloat(v)
		if !ok || f < 0 || f != float64(int(f)) {
			return nil, false
		}
		if f == 0 {
			return nil, true
		}
		if f > 5 {
			// legacy CLI silently used 5 for larger numbers
			gn.Warn("TokensAround %v in old config file is more than 5, using 5", f)
			f = 5
		}
		return int(f), true
	case "datasources", "preferredsources":
		return toIDs(v)
	case "withoutbayes":
		b, ok := toBool(v)
		if !ok {
			return nil, false
		}
		if !b {
			return nil, true
		}
		return false, true
	default:
		b, ok := toBool(v)
		if !ok {
			return nil, false
		}
		if !b {
			return nil, true
		}
		return true, true
	}
}

func toString(v any) (string, bool) {
	s, ok := v.(string)
	return strings.TrimSpace(s), ok
}

func toBool(v any) (bool, bool) {
	switch t := v.(type) {
	case bool:
		return t, true
	case string:
		b, err := strconv.ParseBool(strings.TrimSpace(t))
		return b, err == nil
	}
	return false, false
}

func toFloat(v any) (float64, bool) {
	switch t := v.(type) {
	case int:
		return float64(t), true
	case float64:
		return t, true
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(t), 64)
		return f, err == nil
	}
	return 0, false
}

// toIDs converts a list of data-source IDs. Like the legacy CLI, it also
// accepts a comma-separated string.
func toIDs(v any) (any, bool) {
	var items []any
	switch t := v.(type) {
	case []any:
		items = t
	case string:
		for s := range strings.SplitSeq(t, ",") {
			if s = strings.TrimSpace(s); s != "" {
				items = append(items, s)
			}
		}
	case int:
		items = []any{t}
	default:
		return nil, false
	}

	var res []int
	for _, item := range items {
		f, ok := toFloat(item)
		if !ok || f < 1 || f != float64(int(f)) {
			return nil, false
		}
		res = append(res, int(f))
	}
	if len(res) == 0 {
		return nil, true
	}
	return res, true
}
