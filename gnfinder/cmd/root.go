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
 JSON format to the standard output.`,

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

		findNames(data, lang, bayes, verify)
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
	rootCmd.Flags().BoolP("version", "v", false, "show version")
	rootCmd.Flags().BoolP("bayes", "b", false, "always run Bayes algorithms")
	rootCmd.Flags().BoolP("check-names", "c", false, "verify found name-strings")
	rootCmd.Flags().StringP("lang", "l", "", "text's language")
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
		viper.SetConfigName(".gnfinder")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func findNames(data []byte, langString string, bayes bool,
	verify bool) {

	dictionary := dict.LoadDictionary()
	var output []byte
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
		util.WithResolverVerification(verify),
	)
	output = gnfinder.FindNamesJSON(data, &dictionary, opts...)
	fmt.Println(string(output))
}

func checkStdin() bool {
	stdInFile := os.Stdin
	stat, err := stdInFile.Stat()
	if err != nil {
		log.Panic(err)
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}
