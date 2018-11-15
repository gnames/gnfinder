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
	"io/ioutil"
	"log"
	"os"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/dict"
	"github.com/gnames/gnfinder/lang"
	"github.com/gnames/gnfinder/util"
	"github.com/gnames/gnfinder/verifier"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	buildVersion string
	buildDate    string
	cfgFile      string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gnfinder text_with_names.txt",
	Short: "Finds scientific names in plain texts",
	Long: `Finds scientific names in plain UTF-8 encoded texts

 Name finding happens in two stages. First we apply heuristic rules, and
 then, if it is possible, Bayesian algorithms to find scientific names.
 Optionally, gnfinder verifies found names against gnindex database located
 at https://index.globalnames.org. Found names and metadata are returned in
 JSON format to the standard output.

 Verification returns 'the best' result for the match. If specific datasets
 are important for verification, they can be set with '-s' '--sources' flag
 using IDs from https://index.globalnames.org/datasource. The default sources
 are 'Catalogue of life' (ID 1), GBIF (ID 11), and Open Tree of Life (ID
 179).`,

	// Uncomment the following line if your bare application has an action
	// associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var data []byte

		version, err := cmd.Flags().GetBool("version")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		if version {
			fmt.Printf("version: %s\n\ndate:    %s\n\n", buildVersion, buildDate)
			os.Exit(0)
		}

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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ver string, date string) {
	buildVersion = ver
	buildDate = date
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
	// 	"config file (default is $HOME/.gnfinder.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("version", "v", false, "show version.")
	rootCmd.Flags().BoolP("bayes", "b", false, "always run Bayes algorithms.")
	rootCmd.Flags().BoolP("check-names", "c", false, "verify found name-strings.")
	rootCmd.Flags().StringP("lang", "l", "", "text's language.")
	rootCmd.Flags().IntSliceP("sources", "s", []int{1, 11, 179},
		"IDs of data sources used in verification.")
	log.SetFlags(0)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gnfinder" (without extension).
		viper.AddConfigPath(home)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
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
