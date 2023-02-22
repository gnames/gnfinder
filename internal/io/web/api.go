package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gnames/gndoc"
	"github.com/gnames/gnfinder/internal/ent/api"
	"github.com/gnames/gnfinder/internal/ent/lang"
	"github.com/gnames/gnfinder/internal/ent/output"
	"github.com/gnames/gnfinder/internal/ent/verifier"
	gnfinder "github.com/gnames/gnfinder/pkg"
	"github.com/gnames/gnfinder/pkg/config"
	"github.com/gnames/gnfmt"
	"github.com/labstack/echo/v4"
)

func infoApiGET(gnf gnfinder.GNfinder) func(echo.Context) error {
	return func(c echo.Context) error {
		info := fmt.Sprintf(
			"OpenAPI for gnfinder is described at\n\n%s\n",
			gnf.GetConfig().APIDoc,
		)

		return c.String(http.StatusOK, info)
	}

}

func pingApiGET() func(echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	}
}

func verApiGET(gnf gnfinder.GNfinder) func(echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, gnf.GetVersion())
	}
}

func findApiGET(gnf gnfinder.GNfinder) func(echo.Context) error {
	return func(c echo.Context) error {
		ctx, cancel := getContext(c)
		defer cancel()
		chErr := make(chan error)
		params := paramsFindGET(c)

		go finder(c, gnf, params, chErr)

		select {
		case <-ctx.Done():
			<-chErr
			return ctx.Err()
		case err := <-chErr:
			return err
		case <-time.After(1 * time.Minute):
			return errors.New("request took too long")
		}
	}
}

func paramsFindGET(c echo.Context) api.FinderParams {
	var textURL, text string
	text, _ = url.QueryUnescape(c.Param("text"))
	if strings.HasPrefix(text, "http") {
		textURL = text
		text = ""
	}
	var wordsAround int
	if num, err := strconv.Atoi(c.QueryParam("words_around")); err == nil {
		wordsAround = num
	}
	var sources []int
	elements := strings.Split(c.QueryParam("sources"), ",")
	for _, v := range elements {
		v = strings.TrimSpace(v)
		if num, err := strconv.Atoi(v); err == nil {
			sources = append(sources, num)
		}
	}

	params := api.FinderParams{
		URL:            textURL,
		Text:           text,
		Format:         c.QueryParam("format"),
		Language:       c.QueryParam("language"),
		BytesOffset:    c.QueryParam("bytes_offset") == "true",
		ReturnContent:  c.QueryParam("return_content") == "true",
		UniqueNames:    c.QueryParam("unique_names") == "true",
		AmbiguousNames: c.QueryParam("ambiguous_names") == "true",
		NoBayes:        c.QueryParam("no_bayes") == "true",
		OddsDetails:    c.QueryParam("odds_details") == "true",
		WordsAround:    wordsAround,
		Verification:   c.QueryParam("verification") == "true",
		Sources:        sources,
		AllMatches:     c.QueryParam("all_matches") == "true",
	}

	if len(params.Sources) > 0 || params.AllMatches {
		params.Verification = true
	}

	return params
}

func findApiPOST(gnf gnfinder.GNfinder) func(echo.Context) error {
	return func(c echo.Context) error {
		ctx, cancel := getContext(c)
		defer cancel()
		chErr := make(chan error)

		var params api.FinderParams
		err := c.Bind(&params)
		if err != nil {
			return err
		}

		go finder(c, gnf, params, chErr)

		select {
		case <-ctx.Done():
			<-chErr
			return ctx.Err()
		case err := <-chErr:
			return err
		case <-time.After(1 * time.Minute):
			return errors.New("request took too long")
		}
	}
}

func finder(c echo.Context, gnf gnfinder.GNfinder, params api.FinderParams, chErr chan<- error) {
	var err error
	var text, filename string
	var txtExtr float32
	var out output.Output
	var format gnfmt.Format
	var opts []config.Option

	text, filename, txtExtr, err = getText(
		c,
		params,
		gnf.GetConfig().TikaURL,
	)

	opts, format = getOptsAPI(params)
	gnf = gnf.ChangeConfig(opts...)
	out = gnf.Find(filename, text)
	out.TextExtractionSec = txtExtr
	cfg := gnf.GetConfig()
	if cfg.WithVerification {
		verif := verifier.New(
			cfg.VerifierURL,
			cfg.DataSources,
			cfg.WithAllMatches,
		)
		verifiedNames, stats, dur := verif.Verify(out.UniqueNameStrings())
		out.MergeVerification(verifiedNames, stats, dur)
	}

	out.TotalSec = out.TextExtractionSec + out.NameFindingSec + out.NameVerifSec

	if err == nil {
		if format == gnfmt.CompactJSON || format == gnfmt.PrettyJSON {
			err = c.JSON(http.StatusOK, out)
		} else {
			err = c.String(http.StatusOK, out.Format(format))
		}
	}

	chErr <- err
}

func getOptsAPI(params api.FinderParams) ([]config.Option, gnfmt.Format) {
	format, _ := gnfmt.NewFormat(params.Format)
	if format == gnfmt.FormatNone {
		format = gnfmt.CompactJSON
	}

	opts := []config.Option{
		config.OptWithBayesOddsDetails(params.OddsDetails),
		config.OptFormat(format),
		config.OptWithBayes(!params.NoBayes),
		config.OptWithPositonInBytes(params.BytesOffset),
		config.OptLanguage(getLanguage(params.Language)),
		config.OptDataSources(params.Sources),
		config.OptIncludeInputText(params.ReturnContent),
		config.OptWithAllMatches(params.AllMatches),
		config.OptWithAmbiguousNames(params.AmbiguousNames),
		config.OptWithVerification(
			params.Verification ||
				len(params.Sources) > 0 ||
				params.AllMatches,
		),
		config.OptTokensAround(params.WordsAround),
	}
	return opts, format
}

func getLanguage(s string) lang.Language {
	l, _ := lang.New(s)
	return l
}

func getContext(c echo.Context) (ctx context.Context, cancel func()) {
	ctx = c.Request().Context()
	ctx, cancel = context.WithTimeout(ctx, 1*time.Minute)
	return ctx, cancel
}

func getText(
	c echo.Context,
	params api.FinderParams,
	tikaURL string,
) (string, string, float32, error) {
	var err error
	var txt, filename string
	var dur float32

	if params.Text != "" {
		return params.Text, filename, dur, err
	}

	d := gndoc.New(tikaURL)
	if params.URL != "" {
		txt, dur, err = d.TextFromURL(params.URL)
		return txt, filename, dur, err
	}

	return textFromFile(c, tikaURL)
}
