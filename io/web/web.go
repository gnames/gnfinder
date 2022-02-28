package web

import (
	"net/http"
	"time"

	"github.com/gnames/gndoc"
	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/config"
	"github.com/gnames/gnfinder/ent/api"
	"github.com/gnames/gnfinder/ent/output"
	"github.com/gnames/gnfinder/ent/verifier"
	"github.com/gnames/gnfmt"
	"github.com/labstack/echo/v4"
)

type Duration struct {
	TextExtraction, NameFinding, Verification, Total float32
}

// Data contains information needed to render web-pages.
type Data struct {
	Input       string
	Output      output.Output
	Duration    Duration
	Page        string
	Format      string
	UniqueNames bool
	Version     string
}

func home(gnf gnfinder.GNfinder) func(echo.Context) error {
	return func(c echo.Context) error {
		data := Data{Page: "home", Version: gnf.GetVersion().Version}
		return c.Render(http.StatusOK, "layout", data)
	}
}

func apidoc(gnf gnfinder.GNfinder) func(echo.Context) error {
	return func(c echo.Context) error {
		data := Data{Page: "apidoc", Version: gnf.GetVersion().Version}
		return c.Render(http.StatusOK, "layout", data)
	}
}

func find(gnf gnfinder.GNfinder) func(echo.Context) error {

	return func(c echo.Context) error {
		var err error
		var params api.FinderParams
		var txt, filename string
		var out output.Output
		var dur Duration

		err = c.Bind(&params)
		if err != nil {
			return err
		}

		tikaURL := gnf.GetConfig().TikaURL

		opts := getOpts(params)
		gnf = gnf.ChangeConfig(opts...)

		if t := params.Text; t != "" {
			txt = t
		} else if url := params.URL; url != "" {
			doc := gndoc.New(tikaURL)
			txt, dur.TextExtraction, err = doc.TextFromURL(url)
		} else {
			txt, filename, dur.TextExtraction, err = textFromFile(c, tikaURL)
		}
		if err != nil {
			return err
		}

		out, dur.NameFinding = findNames(gnf, filename, txt)
		dur.Verification = out.NameVerifSec
		dur.Total = dur.NameFinding + dur.TextExtraction + dur.Verification
		data := Data{
			Format:      params.Format,
			Output:      out,
			Duration:    dur,
			Page:        "find",
			UniqueNames: false,
			Version:     gnf.GetVersion().Version,
		}

		data.Output.TextExtractionSec = dur.TextExtraction
		data.Output.NameFindingSec = dur.NameFinding
		data.Output.TotalSec = dur.Total

		switch data.Format {
		case "json":
			return c.JSON(http.StatusOK, data.Output)
		case "csv":
			res := out.Format(gnfmt.CSV)
			return c.String(http.StatusOK, res)
		case "tsv":
			res := out.Format(gnfmt.TSV)
			return c.String(http.StatusOK, res)
		default:
			return c.Render(http.StatusOK, "layout", data)
		}
	}
}

func getOpts(params api.FinderParams) []config.Option {
	if len(params.Sources) > 0 {
		params.Verification = true
	}
	return []config.Option{
		config.OptWithUniqueNames(params.UniqueNames),
		config.OptIncludeInputText(params.ReturnContent),
		config.OptWithVerification(params.Verification),
		config.OptDataSources(params.Sources),
		config.OptWithAllMatches(params.WithAllMatches),
		config.OptWithAmbiguousNames(params.WithAmbiguousNames),
		config.OptWithBayesOddsDetails(params.OddsDetails),
	}
}

// findNames finds names and returns duration of name-finding.
func findNames(
	gnf gnfinder.GNfinder,
	file, txt string,
) (output.Output, float32) {
	start := time.Now()
	cfg := gnf.GetConfig()
	res := gnf.Find(file, txt)
	dur := float32(time.Since(start)) / float32(time.Second)
	if cfg.WithVerification {
		sources := cfg.DataSources
		all := cfg.WithAllMatches
		verif := verifier.New(cfg.VerifierURL, sources, all)
		verifiedNames, stats, durVerif := verif.Verify(res.UniqueNameStrings())
		res.MergeVerification(verifiedNames, stats, durVerif)
	}
	if cfg.IncludeInputText {
		res.InputText = txt
	}
	return res, float32(dur)
}

// textFromFile converts uploaded file into text, returns UTF-8 encoded
// text of the file content, duration of conversion, and a potential error.
func textFromFile(
	c echo.Context,
	tikaURL string,
) (string, string, float32, error) {
	start := time.Now()
	doc := gndoc.New(tikaURL)
	file, err := c.FormFile("file")
	if err != nil {
		return "", "", 0, err
	}

	filename := file.Filename

	f, err := file.Open()
	if err != nil {
		return "", filename, 0, err
	}

	txt, err := doc.GetText(f)
	if err != nil {
		return "", filename, 0, err
	}

	dur := float32(time.Since(start)) / float32(time.Second)
	return txt, filename, dur, nil
}
