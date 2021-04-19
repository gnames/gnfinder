// Copyright © 2019 Dmitry Mozzherin <dmozzherin@gmail.com>
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

	"github.com/gnames/bayes"
	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/ent/lang"
	"github.com/gnames/gnfinder/ent/nlp"
	"github.com/gnames/gnfinder/ent/verifier"
	"github.com/gnames/gnfinder/io/dict"
	"github.com/spf13/cobra"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find text_with_names.txt",
	Short: "Finds scientific names in UTF-8-encoded plain texts",
	Long: `
 Name finding happens in two stages. First we apply heuristic rules, and
 then, unless opted out, Bayesian algorithms to find scientific names.
 Optionally, gnfinder verifies found names against gnindex database located
 at https://index.globalnames.org. Found names and metadata are returned in
 JSON format to the standard output.

 Optional verification process returns 'the best' result for the match. If
 specific datasets are important for verification, they can be set with '-s'
 '--sources' flag using IDs from https://index.globalnames.org/datasource.
 The default sources are 'Catalogue of life' (ID 1), GBIF (ID 11), and Open
 Tree of Life (ID 179).`,
	Run: func(cmd *cobra.Command, args []string) {
		var data []byte
		sources := sources(cmd)

		lang, err := cmd.Flags().GetString("lang")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		tokensNum, err := cmd.Flags().GetInt("tokens-around")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		noBayes, err := cmd.Flags().GetBool("no-bayes")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		oddsDetails, err := cmd.Flags().GetBool("odds-details")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		verify, err :=
			cmd.Flags().GetBool("check-names")
		if err != nil {
			log.Println(err)
			os.Exit(1)
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

		findNames(data, lang, noBayes, verify, sources, tokensNum, oddsDetails)
	},
}

func init() {
	rootCmd.AddCommand(findCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// findCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	findCmd.Flags().BoolP("no-bayes", "n", false, "do not run Bayes algorithms.")
	findCmd.Flags().BoolP("check-names", "c", false, "verify found name-strings.")
	findCmd.Flags().BoolP("odds-details", "o", false, "show details of odds calculation.")
	findCmd.Flags().StringP("lang", "l", "", "text's language or 'detect' for automatic detection.")
	findCmd.Flags().IntSliceP("sources", "s", []int{},
		"IDs of data sources to display for matches, for example '1,11,179'")
	findCmd.Flags().IntP("tokens-around", "t", 0, "number of tokens kept around name-strings")
}

func findNames(
	data []byte,
	langString string,
	noBayes bool,
	verify bool,
	sources []int,
	tokensNum int,
	oddsDetails bool,
) {
	dict := dict.LoadDictionary()
	var opts []gnfinder.Option
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

func sources(cmd *cobra.Command) []int {
	res, err := cmd.Flags().GetIntSlice("sources")
	if err != nil {
		log.Panic(err)
	}
	return res
}
