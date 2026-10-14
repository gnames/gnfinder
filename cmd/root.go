// Copyright © 2018-2026 Dmitry Mozzherin <dmozzherin@gmail.com>
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
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/gnames/gn"
	"github.com/gnames/gndoc"
	"github.com/gnames/gnfinder/internal/iofs"
	"github.com/gnames/gnfinder/internal/iologger"
	gnfinder "github.com/gnames/gnfinder/pkg"
	"github.com/gnames/gnfinder/pkg/config"
	"github.com/gnames/gnfinder/pkg/ent/nlp"
	"github.com/gnames/gnfinder/pkg/ent/verifier"
	"github.com/gnames/gnfinder/pkg/io/dict"
	"github.com/gnames/gnfinder/pkg/io/web"
	"github.com/spf13/cobra"
)

// cfg is created during bootstrap from defaults, config.yaml and
// environment variables. CLI flags are applied to it in runRoot.
var cfg config.Config

var rootCmd = &cobra.Command{
	Version: fmt.Sprintf(
		"version: %s\nbuild:   %s", gnfinder.Version, gnfinder.Build,
	),
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

Settings can be changed in ~/.config/gnfinder/config.yaml or with GNF_
environment variables. Logs are at ~/.local/share/gnfinder/logs.
`,
	PersistentPreRunE: bootstrap,
	RunE:              runRoot,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once
// to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Command failed", "error", err)
		gn.PrintErrorMessage(err)
		os.Exit(1)
	}
}

func init() {
	// Remove the automatic "gnfinder version" prefix
	rootCmd.SetVersionTemplate("{{.Version}}\n")

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
		`Format of the output: "compact", "pretty", "csv", "tsv".
  compact: compact JSON,
  pretty: pretty JSON,
  csv: CSV (DEFAULT),
  tsv: TSV`)
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
To find IDs refer to "https://verifier.globalnames.org/data_sources".
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
	// Override version flag to use -V (consistent with gnames projects)
	rootCmd.Flags().BoolP("version", "V", false, "show version.")
	rootCmd.Flags().IntP("words-around",
		"w", 0, "show this many words surrounding name-strings.")
}

// bootstrap prepares directories, logging and configuration before
// any command runs.
func bootstrap(cmd *cobra.Command, args []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	if err = iofs.EnsureDirs(homeDir); err != nil {
		return err
	}

	// Initialize logging with hardcoded defaults ASAP so all subsequent logs
	// are captured. It is reconfigured later with user's config settings.
	defaultLog := config.New().Log
	if err = iologger.Init(config.LogDir(homeDir), defaultLog, false); err != nil {
		return err
	}

	slog.Info("Bootstrap process started",
		"version", gnfinder.Version,
		"home_dir", homeDir,
		"args", os.Args[1:],
	)

	cfgPath := config.ConfigFilePath(homeDir)
	_, statErr := os.Stat(cfgPath)
	isFirstRun := os.IsNotExist(statErr)

	if isFirstRun {
		gn.Info("Configuration file created at <em>%s</em>", cfgPath)
		gn.Info("Logs are available at <em>%s</em>", config.LogDir(homeDir))
		migrateLegacyConfig(homeDir)
	}

	if err = iofs.EnsureConfigFile(homeDir); err != nil {
		slog.Error("Failed to ensure config file", "error", err)
		return err
	}
	if isFirstRun {
		slog.Info("Config file created", "config_file", cfgPath)
	}

	opts, err := initConfig(homeDir)
	if err != nil {
		return err
	}

	cfg = config.New(opts...)
	cfg.Update(config.OptHomeDir(homeDir))

	if err = reconfigureLogging(cfg); err != nil {
		return err
	}

	slog.Info("Configuration loaded successfully",
		"config_file", cfgPath,
		"log_format", cfg.Log.Format,
		"log_level", cfg.Log.Level,
		"log_destination", cfg.Log.Destination,
	)
	return nil
}

// migrateLegacyConfig moves settings from gnfinder.yml used by gnfinder
// v1.1 and earlier to config.yaml. Failures do not stop gnfinder, because
// config.yaml with default settings is created afterwards.
func migrateLegacyConfig(homeDir string) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return
	}
	legacyPath := filepath.Join(cfgDir, iofs.LegacyConfigFile)
	if _, err = os.Stat(legacyPath); err != nil {
		return
	}

	slog.Info("Legacy config file found", "legacy_path", legacyPath)
	keys, err := iofs.MigrateLegacyConfig(legacyPath, homeDir)
	if err != nil {
		slog.Warn("Failed to migrate legacy config file",
			"legacy_path", legacyPath, "error", err,
		)
		gn.Warn("Cannot migrate settings from old <em>%s</em>", legacyPath)
		return
	}

	if len(keys) > 0 {
		slog.Info("Legacy config file migrated",
			"legacy_path", legacyPath, "keys", keys,
		)
		gn.Info("Settings <em>%s</em> migrated from <em>%s</em>",
			strings.Join(keys, ", "), legacyPath,
		)
	}
	gn.Info("Old config file <em>%s</em> is not used anymore "+
		"and can be deleted", legacyPath)
}

// reconfigureLogging reinitializes the logger with the loaded configuration.
// Appends to existing log file to preserve bootstrap logs.
func reconfigureLogging(cfg config.Config) error {
	logDir := config.LogDir(cfg.HomeDir)
	err := iologger.Init(logDir, cfg.Log, true)
	if err != nil {
		slog.Error("Failed to reconfigure logger", "error", err, "log_dir", logDir)
		return err
	}
	attrs := []any{
		"format", cfg.Log.Format,
		"level", cfg.Log.Level,
		"destination", cfg.Log.Destination,
	}
	if cfg.Log.Destination == "file" {
		attrs = append(attrs, "log_file", filepath.Join(logDir, iologger.LogFile))
	}
	slog.Info("Logger reconfigured", attrs...)
	return nil
}

func runRoot(cmd *cobra.Command, args []string) error {
	cfg.Update(flagOpts(cmd)...)
	slog.Info("Name-finding settings",
		"format", cfg.Format.String(),
		"language", cfg.Language.String(),
		"with_bayes", cfg.WithBayes,
		"with_verification", cfg.WithVerification,
		"data_sources", cfg.DataSources,
	)

	if port, _ := cmd.Flags().GetInt("port"); port > 0 {
		return runWebService(port)
	}

	var data, input string
	var convDur float32
	var err error

	switch len(args) {
	case 0:
		if !stdinHasData() {
			return cmd.Help()
		}
		input = "STDIN"
		var rawData []byte
		rawData, err = io.ReadAll(os.Stdin)
		if err != nil {
			return readInputError(input, err)
		}
		data = string(rawData)
	case 1:
		input = args[0]
		ua := "gnfinder/" + gnfinder.Version +
			" (https://github.com/gnames/gnfinder)"
		d := gndoc.New(cfg.TikaURL, gndoc.OptUserAgent(ua))
		if strings.HasPrefix(input, "http") {
			data, convDur, err = d.TextFromURL(input)
		} else {
			data, convDur, err = d.TextFromFile(input, cfg.WithPlainInput)
		}
		if err != nil {
			return readInputError(input, err)
		}
	default:
		return cmd.Help()
	}
	slog.Info("Input text received", "input", input, "length", len(data))

	if cfg.InputTextOnly {
		fmt.Print(data)
		return nil
	}

	return findNames(data, input, convDur)
}

func newGNfinder() (gnfinder.GNfinder, error) {
	dict, err := dict.LoadDictionary()
	if err != nil {
		return nil, loadDataError("dictionary", err)
	}
	weights, err := nlp.BayesWeights()
	if err != nil {
		return nil, loadDataError("Bayesian weights", err)
	}
	return gnfinder.New(cfg, dict, weights), nil
}

func runWebService(port int) error {
	gnf, err := newGNfinder()
	if err != nil {
		return err
	}

	gn.Info("Starting gnfinder web service on port <em>%d</em>", port)
	if err = web.Run(gnf, port); err != nil {
		return webServiceError(port, err)
	}
	return nil
}

func findNames(data, file string, convDur float32) error {
	gnf, err := newGNfinder()
	if err != nil {
		return err
	}

	res := gnf.Find(file, data)
	res.TextExtractionSec = convDur
	if cfg.WithVerification {
		verif := verifier.New(cfg.VerifierURL, cfg.DataSources, cfg.WithAllMatches)
		verifiedNames, stats, dur := verif.Verify(res.UniqueNameStrings())
		res.MergeVerification(verifiedNames, stats, dur)
	}
	res.TotalSec = res.TextExtractionSec + res.NameFindingSec + res.NameVerifSec
	slog.Info("Name-finding finished",
		"names", len(res.Names),
		"total_sec", res.TotalSec,
	)

	fmt.Println(res.Format(cfg.Format))
	return nil
}

// stdinHasData returns true if data is piped into STDIN.
func stdinHasData() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		slog.Warn("Cannot read from STDIN", "error", err)
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}
