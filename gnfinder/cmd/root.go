// Copyright © 2018 Dmitry Mozzherin <dmozzherin@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/ent/lang"
	"github.com/gnames/gnfinder/ent/nlp"
	"github.com/gnames/gnfinder/ent/verifier"
	"github.com/gnames/gnfinder/io/dict"
	"github.com/gnames/gnfmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gnfinder [flags] text-file.txt",
	Short: "Finds scientific names in UTF-8-encoded plain texts",
	Long: `
Name finding happens in two stages. First we apply heuristic rules, and then,
unless opted out, Bayesian algorithms to find scientific names.  Optionally,
gnfinder verifies found names against gnindex database located at
https://index.globalnames.org. Found names and metadata are returned in JSON
format to the standard output.

Optional verification process returns 'the best' result for the match. If
specific datasets are important for verification, they can be set with '-s'
'--sources' flag using IDs from https://verifier.globalnames.org/datasource.
`,

	// Uncomment the following line if your bare application has an action
	// associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		version, err := cmd.Flags().GetBool("version")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		if version {
			fmt.Printf("\nversion: %s\nbuild: %s\n\n", gnfinder.Version, gnfinder.Build)
			os.Exit(0)
		}

		var data []byte
		var sources []int
		data_sources, _ := cmd.Flags().GetString("sources")
		if data_sources != "" {
			sources = parseDataSources(data_sources)
		}
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
		lang, err := cmd.Flags().GetString("lang")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		wordsNum, err := cmd.Flags().GetInt("words-around")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		noBayes, err := cmd.Flags().GetBool("no-bayes")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		oddsDetails, err := cmd.Flags().GetBool("details-odds")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		verify, err :=
			cmd.Flags().GetBool("verify")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		if len(sources) > 0 {
			verify = true
		}

		switch len(args) {
		case 0:
			if !checkStdin() {
				_ = cmd.Help()
				os.Exit(0)
			}
			data, err = io.ReadAll(os.Stdin)
			if err != nil {
				log.Println(err)
			}
		case 1:
			data, err = os.ReadFile(args[0])
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
		default:
			_ = cmd.Help()
			os.Exit(0)
		}

		findNames(data, format, lang, noBayes,
			verify, sources, wordsNum, oddsDetails)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ver string) {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolP("details-odds", "d", false, "show details of odds calculation.")
	rootCmd.Flags().StringP("lang", "l", "", "text's language or 'detect' for automatic detection.")
	rootCmd.Flags().StringP("format", "f", "csv", `Format of the output: "compact", "pretty", "csv".
  compact: compact JSON,
  pretty: pretty JSON,
  csv: CSV (DEFAULT)`)
	rootCmd.Flags().BoolP("no-bayes", "n", false, "do not run Bayes algorithms.")
	rootCmd.Flags().StringP("sources", "s", "", `IDs of important data-sources to verify against (ex "1,11").
If sources are set and there are matches to their data,
such matches are returned in "preferred_result" results.
To find IDs refer to "https://resolver.globalnames.org/data_sources".
1 - Catalogue of Life
3 - ITIS
4 - NCBI
9 - WoRMS
11 - GBIF
12 - Encyclopedia of Life
167 - IPNI
170 - Arctos
172 - PaleoBioDB
181 - IRMNG`)
	rootCmd.Flags().IntP("words-around",
		"w", 0, "show this many words surrounding name-strings")
	rootCmd.Flags().BoolP("version", "V", false, "show version.")
	rootCmd.Flags().BoolP("verify", "v", false, "verify found name-strings.")
	log.SetFlags(0)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfgFile := "gnfinder"
	// Find home directory.
	home, err := os.UserConfigDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Search config in home directory with name ".gnfinder" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func findNames(
	data []byte,
	format gnfmt.Format,
	langString string,
	noBayes bool,
	verify bool,
	sources []int,
	tokensNum int,
	oddsDetails bool,
) {
	dict := dict.LoadDictionary()
	opts := []gnfinder.Option{gnfinder.OptFormat(format)}
	var weights map[lang.Language]*bayes.NaiveBayes

	if langString == "detect" {
		opts = append(opts, gnfinder.OptWithLanguageDetection(true))
	} else if langString != "" {
		l, err := lang.NewLanguage(langString)
		if err != nil {
			log.Printf("Error: %s\n", err)
			log.Printf("Supported language codes: %s.\n", langsToString())
			log.Printf("To detect language automatically use '-l detect'.")
			os.Exit(1)
		}
		opts = append(opts, gnfinder.OptLanguage(l))
	}

	if tokensNum > 0 {
		opts = append(opts, gnfinder.OptTokensAround(tokensNum))
	}

	if oddsDetails {
		opts = append(opts, gnfinder.OptWithBayesOddsDetails(oddsDetails))
	}

	opts = append(opts, gnfinder.OptPreferredSources(sources))

	if !noBayes {
		weights = nlp.BayesWeights()
	}
	opts = append(opts, gnfinder.OptWithBayes(!noBayes))
	opts = append(opts, gnfinder.OptWithVerification(verify))

	cfg := gnfinder.NewConfig(opts...)
	gnf := gnfinder.New(cfg, dict, weights)
	res := gnf.Find(data)

	if gnf.GetConfig().WithVerification {
		verif := verifier.New(gnf.GetConfig().PreferredSources)
		verifiedNames := verif.Verify(res.UniqueNameStrings())
		res.MergeVerification(verifiedNames)
	}
	fmt.Println(string(res.ToJSON()))
}

func langsToString() string {
	langs := lang.SupportedLanguages()
	res := make([]string, len(langs))
	for i, v := range langs {
		res[i] = v.String()
	}
	return strings.Join(res, ", ")
}

func checkStdin() bool {
	stdInFile := os.Stdin
	stat, err := stdInFile.Stat()
	if err != nil {
		log.Panic(err)
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
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
