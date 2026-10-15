package cmd

import (
	"strconv"
	"strings"

	"github.com/gnames/gn"
	"github.com/gnames/gnfinder/pkg/config"
	"github.com/gnames/gnfinder/pkg/ent/lang"
	"github.com/gnames/gnfmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// flagOpts converts CLI flags into config options. Only flags given by
// the user are converted, so they override config.yaml and environment
// variables, but do not reset them to the flags' default values.
func flagOpts(cmd *cobra.Command) []config.Option {
	var res []config.Option
	fs := cmd.Flags()

	boolFlags := []struct {
		name string
		opt  func(bool) config.Option
	}{
		{"ambiguous-uninomials", config.OptWithAmbiguousNames},
		{"adjust-odds", config.OptWithOddsAdjustment},
		{"bytes-offset", config.OptWithPositonInBytes},
		{"details-odds", config.OptWithBayesOddsDetails},
		{"input", config.OptIncludeInputText},
		{"input-only", config.OptInputTextOnly},
		{"utf8-input", config.OptWithPlainInput},
		{"unique-names", config.OptWithUniqueNames},
		{"verify", config.OptWithVerification},
		{"no-bayes", func(b bool) config.Option {
			return config.OptWithBayes(!b)
		}},
		{"no-tika", func(b bool) config.Option {
			return config.OptWithTika(!b)
		}},
		{"tika-first", config.OptWithTikaFirst},
	}
	for _, v := range boolFlags {
		if b, ok := changedBool(fs, v.name); ok {
			res = append(res, v.opt(b))
		}
	}

	if b, ok := changedBool(fs, "all-matches"); ok {
		res = append(res, config.OptWithAllMatches(b))
		if b {
			res = append(res, config.OptWithVerification(true))
		}
	}

	if s, ok := changedString(fs, "format"); ok {
		if f, err := gnfmt.NewFormat(s); err != nil {
			gn.Warn("Unknown format <em>%s</em>, use csv, tsv, compact or pretty", s)
		} else {
			res = append(res, config.OptFormat(f))
		}
	}

	if s, ok := changedString(fs, "lang"); ok {
		// can take "detect" as a string
		if l, err := lang.New(s); err != nil {
			gn.Warn(
				"Unknown language <em>%s</em>, use one of %s or 'detect'",
				s, langStrings(),
			)
		} else {
			res = append(res, config.OptLanguage(l))
		}
	}

	if s, ok := changedString(fs, "sources"); ok {
		if sources := parseDataSources(s); len(sources) > 0 {
			res = append(res, config.OptDataSources(sources))
		}
		res = append(res, config.OptWithVerification(true))
	}

	if s, ok := changedString(fs, "tika-url"); ok {
		res = append(res, config.OptTikaURL(s))
	}

	if s, ok := changedString(fs, "verifier-url"); ok {
		res = append(res, config.OptVerifierURL(s))
	}

	if fs.Changed("words-around") {
		i, _ := fs.GetInt("words-around")
		res = append(res, config.OptTokensAround(i))
	}

	return res
}

func changedBool(fs *pflag.FlagSet, name string) (bool, bool) {
	if !fs.Changed(name) {
		return false, false
	}
	b, _ := fs.GetBool(name)
	return b, true
}

func changedString(fs *pflag.FlagSet, name string) (string, bool) {
	if !fs.Changed(name) {
		return "", false
	}
	s, _ := fs.GetString(name)
	return s, true
}

func parseDataSources(s string) []int {
	dss := strings.Split(s, ",")
	res := make([]int, 0, len(dss))
	for _, v := range dss {
		v = strings.TrimSpace(v)
		ds, err := strconv.Atoi(v)
		if err != nil {
			gn.Warn("Cannot convert data-source <em>%s</em> to ID, skipping", v)
			continue
		}
		if ds < 1 {
			gn.Warn("Data source ID %d is less than one, skipping", ds)
			continue
		}
		res = append(res, ds)
	}
	return res
}

func langStrings() string {
	return strings.Join(lang.LangStrings(), ", ")
}
