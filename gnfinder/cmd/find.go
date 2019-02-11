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
	"io/ioutil"
	"log"
	"os"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/util"
	"github.com/gnames/gnfinder/verifier"
	"github.com/spf13/cobra"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find text_with_names.txt",
	Short: "Finds scientific names in UTF-8-encoded plain texts",
	Long: `
 Name finding happens in two stages. First we apply heuristic rules, and
 then, if it is possible, Bayesian algorithms to find scientific names.
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
		bayes, err := cmd.Flags().GetBool("bayes")
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
				cmd.Help()
				os.Exit(0)
			}
			data, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Println(err)
			}
		case 1:
			data, err = ioutil.ReadFile(args[0])
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
		default:
			cmd.Help()
			os.Exit(0)
		}

		findNames(data, lang, bayes, verify, sources)
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
	// findCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	findCmd.Flags().BoolP("bayes", "b", false, "always run Bayes algorithms.")
	findCmd.Flags().BoolP("check-names", "c", false, "verify found name-strings.")
	findCmd.Flags().StringP("lang", "l", "", "text's language.")
	findCmd.Flags().IntSliceP("sources", "s", []int{1, 11, 179},
		"IDs of data sources used in verification.")
}

func findNames(data []byte, langString string, bayes bool,
	verify bool, sources []int) {

	dictionary := dict.LoadDictionary()
	var opts []util.Opt
	if langString != "" {
		l, err := lang.NewLanguage(langString)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		opts = append(opts, util.WithLanguage(l))
	}

	opts = append(opts,
		util.WithBayes(bayes),
		util.WithVerification(verify),
		util.WithSources(sources),
	)
	m := util.NewModel(opts...)
	output := gnfinder.FindNames([]rune(string(data)), &dictionary, m)

	if m.Verifier.Verify {
		names := gnfinder.UniqueNameStrings(output.Names)
		namesResolved := verifier.Verify(names, m)
		for i, n := range output.Names {
			if v, ok := namesResolved[n.Name]; ok {
				output.Names[i].Verification = v
			}
		}
	}
	fmt.Println(string(output.ToJSON()))
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
