package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/ent/lang"
	"github.com/gnames/gnfmt"
	"github.com/spf13/cobra"
)

func versionFlag(cmd *cobra.Command) bool {
	version, _ := cmd.Flags().GetBool("version")
	if version {
		fmt.Printf("\nversion: %s\nbuild: %s\n\n", gnfinder.Version, gnfinder.Build)
		return true
	}
	return false
}

func portFlag(cmd *cobra.Command) int {
	port, _ := cmd.Flags().GetInt("port")
	return port
}

func sourcesFlag(cmd *cobra.Command) gnfinder.Option {
	sources, _ := cmd.Flags().GetString("sources")
	return gnfinder.OptPreferredSources(parseDataSources(sources))
}

func formatFlag(cmd *cobra.Command) gnfinder.Option {
	format := gnfmt.CSV
	formatString, _ := cmd.Flags().GetString("format")
	if formatString != "csv" {
		format, _ = gnfmt.NewFormat(formatString)
		if format == gnfmt.FormatNone {
			log.Printf(
				"Cannot set format from '%s', setting format to csv",
				formatString,
			)
		}
	}
	return gnfinder.OptFormat(format)
}

func langFlag(cmd *cobra.Command) gnfinder.Option {
	langString, _ := cmd.Flags().GetString("lang")

	if langString == "" {
		return gnfinder.OptWithLanguageDetection(false)
	}

	if langString == "detect" {
		return gnfinder.OptWithLanguageDetection(true)
	}

	l, err := lang.NewLanguage(langString)
	if err != nil {
		l = lang.DefaultLanguage
		log.Print(err)
		log.Printf("Supported language codes: %s.", langsToString())
		log.Printf("To detect language automatically use '-l detect'.")
	}
	return gnfinder.OptLanguage(l)
}

func wordsFlag(cmd *cobra.Command) gnfinder.Option {
	wordsNum, _ := cmd.Flags().GetInt("words-around")
	return gnfinder.OptTokensAround(wordsNum)
}

func bayesFlag(cmd *cobra.Command) gnfinder.Option {
	noBayes, _ := cmd.Flags().GetBool("no-bayes")
	return gnfinder.OptWithBayes(!noBayes)
}

func oddsDetailsFlag(cmd *cobra.Command) gnfinder.Option {
	oddsDetails, _ := cmd.Flags().GetBool("details-odds")
	return gnfinder.OptWithBayesOddsDetails(oddsDetails)
}

func verifFlag(cmd *cobra.Command) gnfinder.Option {
	verif, _ := cmd.Flags().GetBool("verify")
	return gnfinder.OptWithVerification(verif)
}

func parseDataSources(s string) []int {
	if s == "" {
		return nil
	}
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
