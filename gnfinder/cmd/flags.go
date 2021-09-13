package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/config"
	"github.com/gnames/gnfinder/ent/lang"
	"github.com/gnames/gnfmt"
	"github.com/spf13/cobra"
)

func adjustOddsFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("adjust-odds")
	if b {
		opts = append(opts, config.OptWithOddsAdjustment(b))
	}
}

func bayesFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("no-bayes")
	if b {
		opts = append(opts, config.OptWithBayes(false))
	}
}

func bytesOffsetFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("bytes-offset")
	if b {
		opts = append(opts, config.OptWithBytesOffset(true))
	}
}

func formatFlag(cmd *cobra.Command) {
	format := gnfmt.CSV
	s, _ := cmd.Flags().GetString("format")

	if s == "" {
		return
	}
	if s != "csv" {
		format, _ = gnfmt.NewFormat(s)
		if format == gnfmt.FormatNone {
			log.Printf(
				"Cannot set format from '%s', setting format to csv",
				s,
			)
		}
	}
	opts = append(opts, config.OptFormat(format))
}

func inputFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("input")
	if b {
		opts = append(opts, config.OptIncludeInputText(b))
	}
}

func inputOnlyFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("input-only")
	if b {
		opts = append(opts, config.OptInputTextOnly(b))
	}
}

func langFlag(cmd *cobra.Command) {
	s, _ := cmd.Flags().GetString("lang")

	if s == "" {
		return
	}

	if s == "detect" {
		opts = append(opts, config.OptWithLanguageDetection(true))
		return
	}

	l, err := lang.New(s)
	if err != nil {
		log.Print(err)
		log.Printf("Supported language codes: %s.", langsToString())
		log.Printf("To detect language automatically use '-l detect'.")
	}
	opts = append(opts, config.OptLanguage(l))
}

func oddsDetailsFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("details-odds")
	if b {
		opts = append(opts, config.OptWithBayesOddsDetails(b))
	}
}

func plainInputFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("utf8-input")
	if b {
		opts = append(opts, config.OptWithPlainInput(true))
	}
}

func portFlag(cmd *cobra.Command) int {
	port, _ := cmd.Flags().GetInt("port")
	return port
}

func sourcesFlag(cmd *cobra.Command) {
	s, _ := cmd.Flags().GetString("sources")
	if s != "" {
		opts = append(opts, config.OptPreferredSources(parseDataSources(s)))
	}
}

func tikaURLFlag(cmd *cobra.Command) {
	s, _ := cmd.Flags().GetString("tika-url")
	if s != "" {
		opts = append(opts, config.OptTikaURL(s))
	}
}

func uniqueFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("unique-names")
	if b {
		opts = append(opts, config.OptWithUniqueNames(b))
	}
}

func verifFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("verify")
	if b {
		opts = append(opts, config.OptWithVerification(b))
	}
}

func verifURLFlag(cmd *cobra.Command) {
	s, _ := cmd.Flags().GetString("verifier-url")
	if s != "" {
		opts = append(opts, config.OptVerifierURL(s))
	}
}

func versionFlag(cmd *cobra.Command) bool {
	b, _ := cmd.Flags().GetBool("version")
	if b {
		fmt.Printf("\nversion: %s\nbuild: %s\n\n", gnfinder.Version, gnfinder.Build)
		return true
	}
	return false
}

func wordsFlag(cmd *cobra.Command) {
	i, _ := cmd.Flags().GetInt("words-around")
	if i > 0 {
		opts = append(opts, config.OptTokensAround(i))
	}
}

func parseDataSources(s string) []int {
	dss := strings.Split(s, ",")
	res := make([]int, 0, len(dss))
	for _, v := range dss {
		v = strings.Trim(v, " ")
		ds, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("Cannot convert data-source '%s' to list, skipping", v)
			return nil
		}
		if ds < 1 {
			log.Printf("Data source ID %d is less than one, skipping", ds)
		} else {
			res = append(res, int(ds))
		}
	}
	if len(res) > 0 {
		return res
	}
	return nil
}
