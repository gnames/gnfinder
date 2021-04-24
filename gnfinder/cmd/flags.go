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

func sourcesFlag(cmd *cobra.Command) config.Option {
	sources, _ := cmd.Flags().GetString("sources")
	return config.OptPreferredSources(parseDataSources(sources))
}

func formatFlag(cmd *cobra.Command) config.Option {
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
	return config.OptFormat(format)
}

func langFlag(cmd *cobra.Command) config.Option {
	langString, _ := cmd.Flags().GetString("lang")

	if langString == "" {
		return config.OptWithLanguageDetection(false)
	}

	if langString == "detect" {
		return config.OptWithLanguageDetection(true)
	}

	l, err := lang.NewLanguage(langString)
	if err != nil {
		l = lang.DefaultLanguage
		log.Print(err)
		log.Printf("Supported language codes: %s.", langsToString())
		log.Printf("To detect language automatically use '-l detect'.")
	}
	return config.OptLanguage(l)
}

func wordsFlag(cmd *cobra.Command) config.Option {
	wordsNum, _ := cmd.Flags().GetInt("words-around")
	return config.OptTokensAround(wordsNum)
}

func bayesFlag(cmd *cobra.Command) config.Option {
	noBayes, _ := cmd.Flags().GetBool("no-bayes")
	return config.OptWithBayes(!noBayes)
}

func adjustOddsFlag(cmd *cobra.Command) config.Option {
	adj, _ := cmd.Flags().GetBool("adjust-odds")
	return config.OptWithOddsAdjustment(adj)
}

func oddsDetailsFlag(cmd *cobra.Command) config.Option {
	oddsDetails, _ := cmd.Flags().GetBool("details-odds")
	return config.OptWithBayesOddsDetails(oddsDetails)
}

func verifFlag(cmd *cobra.Command) config.Option {
	verif, _ := cmd.Flags().GetBool("verify")
	return config.OptWithVerification(verif)
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
