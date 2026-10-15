package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/gnames/gnfinder/internal/iofs"
	"github.com/gnames/gnfinder/pkg/config"
	"github.com/gnames/gnfinder/pkg/ent/lang"
	"github.com/gnames/gnfmt"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

// envVars maps config keys to environment variables. Names are given
// explicitly, because viper does not add a prefix to explicitly named
// variables. These match the fields included in config.ToOptions().
var envVars = map[string]string{
	"bayes_odds_threshold":    "GNF_BAYES_ODDS_THRESHOLD",
	"data_sources":            "GNF_DATA_SOURCES",
	"format":                  "GNF_FORMAT",
	"include_input_text":      "GNF_INCLUDE_INPUT_TEXT",
	"input_text_only":         "GNF_INPUT_TEXT_ONLY",
	"language":                "GNF_LANGUAGE",
	"tika_url":                "GNF_TIKA_URL",
	"tokens_around":           "GNF_TOKENS_AROUND",
	"verifier_url":            "GNF_VERIFIER_URL",
	"with_all_matches":        "GNF_WITH_ALL_MATCHES",
	"with_ambiguous_names":    "GNF_WITH_AMBIGUOUS_NAMES",
	"with_bayes":              "GNF_WITH_BAYES",
	"with_bayes_odds_details": "GNF_WITH_BAYES_ODDS_DETAILS",
	"with_odds_adjustment":    "GNF_WITH_ODDS_ADJUSTMENT",
	"with_plain_input":        "GNF_WITH_PLAIN_INPUT",
	"with_position_in_bytes":  "GNF_WITH_POSITION_IN_BYTES",
	"with_tika":               "GNF_WITH_TIKA",
	"with_tika_first":         "GNF_WITH_TIKA_FIRST",
	"with_unique_names":       "GNF_WITH_UNIQUE_NAMES",
	"with_verification":       "GNF_WITH_VERIFICATION",
	"log.format":              "GNF_LOG_FORMAT",
	"log.level":               "GNF_LOG_LEVEL",
	"log.destination":         "GNF_LOG_DESTINATION",
}

// envWithoutBayes is a legacy environment variable, it is the inverse
// of GNF_WITH_BAYES.
const envWithoutBayes = "GNF_WITHOUT_BAYES"

// initConfig reads config.yaml and environment variables and returns
// options for persistent configuration fields.
func initConfig(home string) ([]config.Option, error) {
	var err error
	cfgPath := config.ConfigFilePath(home)
	slog.Info("Initializing configuration from file", "config_path", cfgPath)

	v := viper.New()
	v.SetConfigFile(cfgPath)
	initEnvVars(v)

	if err = v.ReadInConfig(); err != nil {
		slog.Error("Failed to read config file", "error", err, "config_path", cfgPath)
		return nil, iofs.ReadFileError(cfgPath, err)
	}

	// Values missing from config.yaml and environment keep their defaults.
	res := config.New()
	if err = v.Unmarshal(&res, viper.DecodeHook(decodeHook())); err != nil {
		slog.Error("Failed to unmarshal config", "error", err, "config_path", cfgPath)
		return nil, configValueError(cfgPath, err)
	}

	opts := res.ToOptions()
	if b, ok := withoutBayesEnv(); ok {
		opts = append(opts, config.OptWithBayes(!b))
	}

	slog.Info("Configuration read successfully",
		"config_path", cfgPath,
		"options_count", len(opts),
	)
	return opts, nil
}

func initEnvVars(v *viper.Viper) {
	for k, env := range envVars {
		_ = v.BindEnv(k, env)
	}
	slog.Info("Environment variables bound", "count", len(envVars))
}

// withoutBayesEnv returns the value of the legacy GNF_WITHOUT_BAYES variable
// and true, if the variable is set to a valid boolean.
func withoutBayesEnv() (bool, bool) {
	s, ok := os.LookupEnv(envWithoutBayes)
	if !ok {
		return false, false
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		slog.Warn("Cannot parse environment variable", "name", envWithoutBayes, "value", s)
		return false, false
	}
	return b, true
}

// decodeHook converts strings from config.yaml or environment into
// gnfmt.Format, lang.Language and []int. It replaces viper's default hooks, so
// they are included as well.
func decodeHook() mapstructure.DecodeHookFunc {
	return mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		stringToTypesHook,
	)
}

func stringToTypesHook(from, to reflect.Type, data any) (any, error) {
	if from.Kind() != reflect.String {
		return data, nil
	}
	s := data.(string)

	switch to {
	case reflect.TypeFor[gnfmt.Format]():
		f, err := gnfmt.NewFormat(s)
		if err != nil {
			return nil, err
		}
		return f, nil
	case reflect.TypeFor[[]int]():
		// comma-separated list, for example GNF_DATA_SOURCES="1,11"
		var res []int
		for v := range strings.SplitSeq(s, ",") {
			if v = strings.TrimSpace(v); v == "" {
				continue
			}
			i, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("cannot convert '%s' to a list of integers", s)
			}
			res = append(res, i)
		}
		return res, nil
	case reflect.TypeFor[lang.Language]():
		l, err := lang.New(s)
		if err != nil {
			return nil, fmt.Errorf(
				"unknown language '%s', use one of %s or 'detect'",
				s, langStrings(),
			)
		}
		return l, nil
	}
	return data, nil
}
