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
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/gnames/gndoc"
	gnfinder "github.com/gnames/gnfinder/pkg"
	"github.com/gnames/gnfinder/pkg/config"
	"github.com/gnames/gnfinder/pkg/ent/lang"
	"github.com/gnames/gnfinder/pkg/ent/nlp"
	"github.com/gnames/gnfinder/pkg/ent/verifier"
	"github.com/gnames/gnfinder/pkg/io/dict"
	"github.com/gnames/gnfinder/pkg/io/web"
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
	BayesOddsThreshold float64
	DataSources        []int
	Format             string
	IncludeInputText   bool
	InputTextOnly      bool
	Language           string
	// PreferredSources is deprecated
	PreferredSources     []int
	TikaURL              string
	TokensAround         int
	VerifierURL          string
	WithAllMatches       bool
	WithAmbiguousNames   bool
	WithBayesOddsDetails bool
	WithOddsAdjustment   bool
	WithPlainInput       bool
	WithPositionInBytes  bool
	WithUniqueNames      bool
	WithVerification     bool
	WithoutBayes         bool
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

Optional verification process returns 'the best' result for the match.
If verification needs to be limited to specific data-sources, they can be set
with '-s' '--sources' flag using IDs from
https://verifier.globalnames.org/data_sources.
If flag '-M' '--all-matches' is set, verification returns all found
verification results.
`,

	// Uncomment the following line if your bare application has an action
	// associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if versionFlag(cmd) {
			os.Exit(0)
		}

		ambiguousUninomialsFlag(cmd)
		adjustOddsFlag(cmd)
		bayesFlag(cmd)
		bytesOffsetFlag(cmd)
		formatFlag(cmd)
		inputFlag(cmd)
		inputOnlyFlag(cmd)
		langFlag(cmd)
		allMatchesFlag(cmd)
		oddsDetailsFlag(cmd)
		plainInputFlag(cmd)
		sourcesFlag(cmd)
		tikaURLFlag(cmd)
		uniqueFlag(cmd)
		verifFlag(cmd)
		verifURLFlag(cmd)
		wordsFlag(cmd)

		cfg := config.New(opts...)

		if port := portFlag(cmd); port > 0 {
			dict, err := dict.LoadDictionary()
			if err != nil {
				slog.Error("Cannot load dictionary", "error", err)
				os.Exit(1)
			}
			weights, err := nlp.BayesWeights()
			if err != nil {
				slog.Error("Cannot load Bayesian weights", "error", err)
				os.Exit(1)
			}

			gnf := gnfinder.New(cfg, dict, weights)
			err = web.Run(gnf, port)
			if err != nil {
				slog.Error("Web service stopped suddenly", "error", err)
				os.Exit(1)
			}

			os.Exit(0)
		}

		var data, input string
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
				slog.Error("Cannot read data", "error", err)
			}
			data = string(rawData)
			input = "STDIN"
		case 1:
			input = args[0]
			d := gndoc.New(cfg.TikaURL)
			if strings.HasPrefix(input, "http") {
				data, convDur, err = d.TextFromURL(input)
			} else {
				data, convDur, err = d.TextFromFile(input, cfg.WithPlainInput)
			}
			if err != nil {
				slog.Error("Cannot get input", "error", err)
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

		findNames(data, cfg, input, convDur)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolP("ambiguous-uninomials", "A", false,
		"preserve uninomials that are also common words.")
	rootCmd.Flags().BoolP("adjust-odds", "a", false,
		"adjust Bayes odds using density of found names.")
	rootCmd.Flags().BoolP("bytes-offset", "b", false,
		"names offsets in bytes, not UTF-8 chars.")
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
	rootCmd.Flags().BoolP("all-matches", "M", false,
		"verification returns all found matches")
	rootCmd.Flags().BoolP("utf8-input", "U", false,
		`affirm that the input is a plain text UTF8 file.
The direct reading of a file will be used instead of a remote
Apache Tika service.`)
	rootCmd.Flags().BoolP("unique-names", "u", false,
		"return unique names list")
	rootCmd.Flags().BoolP("verify", "v", false, "verify found name-strings.")
	rootCmd.Flags().BoolP("version", "V", false, "show version.")
	rootCmd.Flags().IntP("words-around",
		"w", 0, "show this many words surrounding name-strings.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var configDir string
	var err error
	configFile := "gnfinder"

	// Find config directory.
	configDir, err = os.UserConfigDir()
	if err != nil {
		slog.Error("Cannot find config directory", "error", err)
		os.Exit(1)
	}

	// Search config in home directory with name ".gnmatcher" (without extension).
	viper.AddConfigPath(configDir)
	viper.SetConfigName(configFile)

	// Set environment variables to override
	// config file settings
	_ = viper.BindEnv("BayesOddsThreshold", "GNF_BAYES_ODDS_THRESHOLD")
	_ = viper.BindEnv("DataSources", "GNF_DATA_SOURCES")
	_ = viper.BindEnv("Format", "GNF_FORMAT")
	_ = viper.BindEnv("InputTextOnly", "GNF_INPUT_TEXT_ONLY")
	_ = viper.BindEnv("IncludeInputText", "GNF_INCLUDE_INPUT_TEXT")
	_ = viper.BindEnv("Language", "GNF_LANGUAGE")
	_ = viper.BindEnv("TikaURL", "GNF_TIKA_URL")
	_ = viper.BindEnv("TokensAround", "GNF_TOKENS_AROUND")
	_ = viper.BindEnv("VerifierURL", "GNF_VERIFIER_URL")
	_ = viper.BindEnv("WithAmbiguousNames", "GNF_WITH_AMBIGUOUS_NAMES")
	_ = viper.BindEnv("WithAllMatches", "GNF_WITH_ALL_MATCHES")
	_ = viper.BindEnv("WithBayesOddsDetails", "GNF_WITH_BAYES_ODDS_DETAILS")
	_ = viper.BindEnv("WithOddsAdjustment", "GNF_WITH_ODDS_ADJUSTMENT")
	_ = viper.BindEnv("WithPlainInput", "GNF_WITH_PLAIN_INPUT")
	_ = viper.BindEnv("WithPositionInBytes", "GNF_WITH_POSITION_IN_BYTES")
	_ = viper.BindEnv("WithUniqueNames", "GNF_WITH_UNIQUE_NAMES")
	_ = viper.BindEnv("WithVerification", "GNF_WITH_VERIFICATION")
	_ = viper.BindEnv("WithoutBayes", "GNF_WITHOUT_BAYES")

	viper.AutomaticEnv() // read in environment variables that match
	configPath := filepath.Join(configDir, fmt.Sprintf("%s.yml", configFile))
	_ = touchConfigFile(configPath)

	// If a config file is found, read it in.
	err = viper.ReadInConfig()
	if err != nil {
		slog.Error("Cannot use config file", "config", viper.ConfigFileUsed())
		os.Exit(1)
	} else {
		getOpts()
	}
}

// getOpts imports data from the configuration file. Some of the settings can
// be overriden by command line flags.
func getOpts() {
	cfgCli := &cfgData{}
	err := viper.Unmarshal(cfgCli)
	if err != nil {
		slog.Error("Cannot deserialize config data", "error", err)
		os.Exit(1)
	}

	if cfgCli.BayesOddsThreshold > 0 {
		opts = append(opts,
			config.OptBayesOddsThreshold(cfgCli.BayesOddsThreshold))
	}

	if len(cfgCli.DataSources) > 0 || len(cfgCli.PreferredSources) > 0 {
		ds := cfgCli.DataSources
		if len(ds) == 0 {
			ds = cfgCli.PreferredSources
		}
		opts = append(opts, config.OptDataSources(ds))
	}

	if cfgCli.Format != "" {
		cfgFormat, err := gnfmt.NewFormat(cfgCli.Format)
		if err != nil {
			cfgFormat = gnfmt.CSV
		}
		opts = append(opts, config.OptFormat(cfgFormat))
	}

	if cfgCli.IncludeInputText {
		opts = append(opts, config.OptIncludeInputText(cfgCli.IncludeInputText))
	}

	if cfgCli.InputTextOnly {
		opts = append(opts, config.OptInputTextOnly(cfgCli.InputTextOnly))
	}

	if cfgCli.Language != "" {
		l, _ := lang.New(cfgCli.Language)
		opts = append(opts, config.OptLanguage(l))
	}

	if cfgCli.TikaURL != "" {
		opts = append(opts, config.OptTikaURL(cfgCli.TikaURL))
	}

	if cfgCli.TokensAround > 0 {
		opts = append(opts, config.OptTokensAround(cfgCli.TokensAround))
	}

	if cfgCli.VerifierURL != "" {
		opts = append(opts, config.OptVerifierURL(cfgCli.VerifierURL))
	}

	if cfgCli.WithAllMatches {
		opts = append(opts, config.OptWithAllMatches(true))
	}

	if cfgCli.WithAmbiguousNames {
		opts = append(opts, config.OptWithAmbiguousNames(true))
	}

	if cfgCli.WithBayesOddsDetails {
		opts = append(opts, config.OptWithBayesOddsDetails(true))
	}

	if cfgCli.WithPlainInput {
		opts = append(opts, config.OptWithPlainInput(true))
	}

	if cfgCli.WithPositionInBytes {
		opts = append(opts, config.OptWithPositonInBytes(true))
	}

	if cfgCli.WithOddsAdjustment {
		opts = append(opts, config.OptWithOddsAdjustment(true))
	}

	if cfgCli.WithUniqueNames {
		opts = append(opts, config.OptWithUniqueNames(true))
	}

	if cfgCli.WithVerification {
		opts = append(opts, config.OptWithVerification(true))
	}

	if cfgCli.WithoutBayes {
		opts = append(opts, config.OptWithBayes(false))
	}

}

func findNames(
	data string,
	cfg config.Config,
	file string,
	convDur float32,
) {
	dict, err := dict.LoadDictionary()
	if err != nil {
		slog.Error("Cannot load dictionary", "error", err)
		os.Exit(1)
	}
	weights, err := nlp.BayesWeights()
	if err != nil {
		slog.Error("Cannot load Bayesian weights", "error", err)
		os.Exit(1)
	}

	gnf := gnfinder.New(cfg, dict, weights)
	res := gnf.Find(file, data)
	res.TextExtractionSec = convDur
	if cfg.WithVerification {
		sources := cfg.DataSources
		all := cfg.WithAllMatches
		verif := verifier.New(cfg.VerifierURL, sources, all)
		verifiedNames, stats, dur := verif.Verify(res.UniqueNameStrings())
		res.MergeVerification(verifiedNames, stats, dur)
	}
	res.TotalSec = res.TextExtractionSec + res.NameFindingSec + res.NameVerifSec
	fmt.Println(res.Format(cfg.Format))
}

func langStrings() string {
	langs := lang.LangStrings()
	return strings.Join(langs, ", ")
}

func checkStdin() bool {
	stdInFile := os.Stdin
	stat, err := stdInFile.Stat()
	if err != nil {
		slog.Error("Cannot read from Stdin", "error", err)
		os.Exit(1)
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

// touchConfigFile checks if config file exists, and if not, it gets created.
func touchConfigFile(configPath string) error {
	fileExists, _ := gnsys.FileExists(configPath)
	if fileExists {
		return nil
	}

	log.Printf("Creating config file: %s.", configPath)
	return createConfig(configPath)
}

// createConfig creates config file.
func createConfig(path string) error {
	err := gnsys.MakeDir(filepath.Dir(path))
	if err != nil {
		slog.Error("Cannot create dir", "path", path, "error", err)
		return err
	}

	err = os.WriteFile(path, []byte(configText), 0644)
	if err != nil {
		slog.Error("Cannot write to file", "path", path, "error", err)
		return err
	}
	return nil
}
