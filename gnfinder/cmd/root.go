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
	"strings"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/ent/lang"
	"github.com/gnames/gnfinder/ent/nlp"
	"github.com/gnames/gnfinder/ent/verifier"
	"github.com/gnames/gnfinder/io/dict"
	"github.com/gnames/gnfinder/io/rest"
	"github.com/spf13/cobra"
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
		var err error
		if versionFlag(cmd) {
			os.Exit(0)
		}

		if port := portFlag(cmd); port > 0 {
			dict := dict.LoadDictionary()
			weights := nlp.BayesWeights()

			cfg := gnfinder.NewConfig()
			gnf := gnfinder.New(cfg, dict, weights)
			rest.Run(gnf, port)

			os.Exit(0)
		}

		opts := []gnfinder.Option{
			formatFlag(cmd),
			langFlag(cmd),
			wordsFlag(cmd),
			bayesFlag(cmd),
			oddsDetailsFlag(cmd),
			verifFlag(cmd),
			sourcesFlag(cmd),
		}

		var data []byte
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

		findNames(data, opts)
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
	rootCmd.Flags().IntP("port",
		"p", 0, "port to run the gnfinder's RESTful API service.")
	rootCmd.Flags().IntP("words-around",
		"w", 0, "show this many words surrounding name-strings.")
	rootCmd.Flags().BoolP("version", "V", false, "show version.")
	rootCmd.Flags().BoolP("verify", "v", false, "verify found name-strings.")
	log.SetFlags(0)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}

func findNames(data []byte, opts []gnfinder.Option) {
	dict := dict.LoadDictionary()
	weights := nlp.BayesWeights()

	cfg := gnfinder.NewConfig(opts...)
	gnf := gnfinder.New(cfg, dict, weights)
	res := gnf.Find(data)

	if gnf.GetConfig().WithVerification {
		verif := verifier.New(gnf.GetConfig().PreferredSources)
		verifiedNames := verif.Verify(res.UniqueNameStrings())
		res.MergeVerification(verifiedNames)
	}
	fmt.Println(res.Format(cfg.Format))
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
