// Copyright © 2018-2021 Dmitry Mozzherin <dmozzherin@gmail.com>
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
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gnames/gndoc"
	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/config"
	"github.com/gnames/gnfinder/ent/lang"
	"github.com/gnames/gnfinder/ent/nlp"
	"github.com/gnames/gnfinder/ent/verifier"
	"github.com/gnames/gnfinder/io/dict"
	"github.com/gnames/gnfinder/io/rest"
	"github.com/gnames/gnfmt"
	"github.com/gnames/gnsys"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed gnfinder.yml
var configText string

var opts []config.Option

// cfgData purpose is to achieve automatic import of data from the
// configuration file, if it exists.
type cfgData struct {
	BayesOddsThreshold    float64
	Format                string
	IncludeInputText      bool
	InputTextOnly         bool
	Language              string
	PreferredSources      []int
	TikaURL               string
	TokensAround          int
	VerifierURL           string
	WithBayesOddsDetails  bool
	WithLanguageDetection bool
	WithOddsAdjustment    bool
	WithPlainInput        bool
	WithUniqueNames       bool
	WithVerification      bool
	WithoutBayes          bool
}

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

			cfg := config.New()
			gnf := gnfinder.New(cfg, dict, weights)
			rest.Run(gnf, port)

			os.Exit(0)
		}

		adjustOddsFlag(cmd)
		bayesFlag(cmd)
		formatFlag(cmd)
		inputFlag(cmd)
		inputOnlyFlag(cmd)
		langFlag(cmd)
		oddsDetailsFlag(cmd)
		plainInputFlag(cmd)
		sourcesFlag(cmd)
		tikaURLFlag(cmd)
		uniqueFlag(cmd)
		verifFlag(cmd)
		verifURLFlag(cmd)
		wordsFlag(cmd)

		cfg := config.New(opts...)

		var data, file string
		var rawData []byte
		var convDur float32
		switch len(args) {
		case 0:
			if !checkStdin() {
				_ = cmd.Help()
				os.Exit(0)
			}
			rawData, err = io.ReadAll(os.Stdin)
			if err != nil {
				log.Println(err)
			}
			data = string(rawData)
			file = "STDIN"
		case 1:
			file = args[0]
			d := gndoc.New(cfg.TikaURL)
			data, convDur, err = d.TextFromFile(file, cfg.WithPlainInput)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
		default:
			_ = cmd.Help()
			os.Exit(0)
		}

		if cfg.InputTextOnly {
			fmt.Print(data)
			os.Exit(0)
		}

		findNames(data, cfg, file, convDur)
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

	rootCmd.Flags().BoolP("adjust-odds", "a", false,
		"adjust Bayes odds using density of found names.")
	rootCmd.Flags().BoolP("details-odds", "d", false,
		"show details of odds calculation.")
	rootCmd.Flags().StringP("verifier-url", "e", "",
		"custom URL for name-verification service.")
	rootCmd.Flags().StringP("format", "f", "",
		`Format of the output: "compact", "pretty", "csv".
  compact: compact JSON,
  pretty: pretty JSON,
  csv: CSV (DEFAULT)`)
	rootCmd.Flags().BoolP("input-only", "I", false,
		"return only given UTF8-encoded input without finding names.")
	rootCmd.Flags().BoolP("input", "i", false,
		"add given input to results.")
	rootCmd.Flags().StringP("lang", "l", "",
		"text's language or 'detect' for automatic detection.")
	rootCmd.Flags().BoolP("no-bayes", "n", false, "do not run Bayes algorithms.")
	rootCmd.Flags().IntP("port",
		"p", 0, "port to run the gnfinder's RESTful API service.")
	rootCmd.Flags().StringP("sources", "s", "",
		`IDs of important data-sources to verify against (ex "1,11").
If sources are set and there are matches to their data,
such matches are returned in "preferred-result" results.
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
	rootCmd.Flags().StringP("tika-url", "t", "",
		`custom URL for the Apache Tika service.
The service is used for converting files into UTF8-encoded text.`)
	rootCmd.Flags().BoolP("utf8-input", "U", false,
		`affurm that the input is a plain text UTF8 file.
The direct reading of a file will be used instead of a remote
Apache Tika service.`)
	rootCmd.Flags().BoolP("unique-names", "u", false,
		"return unique names list")
	rootCmd.Flags().BoolP("verify", "v", false, "verify found name-strings.")
	rootCmd.Flags().BoolP("version", "V", false, "show version.")
	rootCmd.Flags().IntP("words-around",
		"w", 0, "show this many words surrounding name-strings.")
	log.SetFlags(0)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var configDir string
	var err error
	configFile := "gnfinder"

	// Find config directory.
	configDir, err = os.UserConfigDir()
	if err != nil {
		log.Printf("Cannot find config directory: %s.", err)
		return
	}

	// Search config in home directory with name ".gnmatcher" (without extension).
	viper.AddConfigPath(configDir)
	viper.SetConfigName(configFile)

	// Set environment variables to override
	// config file settings
	viper.BindEnv("BayesOddsThreshold", "GNF_BAYES_ODDS_THRESHOLD")
	viper.BindEnv("Format", "GNF_FORMAT")
	viper.BindEnv("IncludeInputText", "GNF_INCLUDE_INPUT_TEXT")
	viper.BindEnv("Language", "GNF_LANGUAGE")
	viper.BindEnv("PreferredSources", "GNF_PREFERRED_SOURCES")
	viper.BindEnv("TikaURL", "GNF_TIKA_URL")
	viper.BindEnv("TokensAround", "GNF_TOKENS_AROUND")
	viper.BindEnv("VerifierURL", "GNF_VERIFIER_URL")
	viper.BindEnv("WithBayesOddsDetails", "GNF_WITH_BAYES_ODDS_DETAILS")
	viper.BindEnv("WithLanguageDetection", "GNF_WITH_LANGUAGE_DETECTION")
	viper.BindEnv("WithOddsAdjustment", "GNF_WITH_ODDS_ADJUSTMENT")
	viper.BindEnv("WithPlainInput", "GNF_WITH_PLAIN_INPUT")
	viper.BindEnv("WithUniqueNames", "GNF_WITH_UNIQUE_NAMES")
	viper.BindEnv("WithVerification", "GNF_WITH_VERIFICATION")
	viper.BindEnv("WithoutBayes", "GNF_WITHOUT_BAYES")

	viper.AutomaticEnv() // read in environment variables that match
	configPath := filepath.Join(configDir, fmt.Sprintf("%s.yml", configFile))
	_ = touchConfigFile(configPath, configFile)

	// If a config file is found, read it in.
	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Cannot use config file: %s", viper.ConfigFileUsed())
	} else {
		getOpts()
	}
}

// getOpts imports data from the configuration file. Some of the settings can
// be overriden by command line flags.
func getOpts() {
	cfg := &cfgData{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		log.Fatalf("Cannot deserialize config data: %s.", err)
	}

	if cfg.BayesOddsThreshold > 0 {
		opts = append(opts,
			config.OptBayesOddsThreshold(cfg.BayesOddsThreshold))
	}

	if cfg.Format != "" {
		cfgFormat, err := gnfmt.NewFormat(cfg.Format)
		if err != nil {
			cfgFormat = gnfmt.CSV
		}
		opts = append(opts, config.OptFormat(cfgFormat))
	}

	if cfg.IncludeInputText {
		opts = append(opts, config.OptIncludeInputText(cfg.IncludeInputText))
	}

	if cfg.InputTextOnly {
		opts = append(opts, config.OptInputTextOnly(cfg.InputTextOnly))
	}

	if cfg.Language != "" {
		l, err := lang.NewLanguage(cfg.Language)
		if err != nil {
			log.Print(err)
		}
		opts = append(opts, config.OptLanguage(l))
	}

	if len(cfg.PreferredSources) > 0 {
		opts = append(opts, config.OptPreferredSources(cfg.PreferredSources))
	}

	if cfg.TikaURL != "" {
		opts = append(opts, config.OptTikaURL(cfg.TikaURL))
	}

	if cfg.TokensAround > 0 {
		opts = append(opts, config.OptTokensAround(cfg.TokensAround))
	}

	if cfg.VerifierURL != "" {
		opts = append(opts, config.OptVerifierURL(cfg.VerifierURL))
	}

	if cfg.WithBayesOddsDetails {
		opts = append(opts, config.OptWithBayesOddsDetails(true))
	}

	if cfg.WithLanguageDetection {
		opts = append(opts, config.OptWithLanguageDetection(true))
	}

	if cfg.WithOddsAdjustment {
		opts = append(opts, config.OptWithOddsAdjustment(true))
	}

	if cfg.WithPlainInput {
		opts = append(opts, config.OptWithPlainInput(true))
	}

	if cfg.WithUniqueNames {
		opts = append(opts, config.OptWithUniqueNames(true))
	}

	if cfg.WithVerification {
		opts = append(opts, config.OptWithVerification(true))
	}

	if cfg.WithoutBayes {
		opts = append(opts, config.OptWithBayes(false))
	}

}

func findNames(
	data string,
	cfg config.Config,
	file string,
	convDur float32,
) {
	dict := dict.LoadDictionary()
	weights := nlp.BayesWeights()

	gnf := gnfinder.New(cfg, dict, weights)
	res := gnf.Find(file, data)
	res.Meta.FileConversionSec = convDur
	if gnf.GetConfig().WithVerification {
		sources := gnf.GetConfig().PreferredSources
		verif := verifier.New(cfg.VerifierURL, sources)
		verifiedNames, dur := verif.Verify(res.UniqueNameStrings())
		res.MergeVerification(verifiedNames, dur)
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

// touchConfigFile checks if config file exists, and if not, it gets created.
func touchConfigFile(configPath string, configFile string) error {
	fileExists, _ := gnsys.FileExists(configPath)
	if fileExists {
		return nil
	}

	log.Printf("Creating config file: %s.", configPath)
	return createConfig(configPath, configFile)
}

// createConfig creates config file.
func createConfig(path string, file string) error {
	err := gnsys.MakeDir(filepath.Dir(path))
	if err != nil {
		log.Printf("Cannot create dir %s: %s.", path, err)
		return err
	}

	err = os.WriteFile(path, []byte(configText), 0644)
	if err != nil {
		log.Printf("Cannot write to file %s: %s", path, err)
		return err
	}
	return nil
}
