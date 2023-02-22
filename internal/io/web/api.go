package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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

func infoAPI(gnf gnfinder.GNfinder) func(echo.Context) error {
	return func(c echo.Context) error {
		info := fmt.Sprintf(
			"OpenAPI for gnfinder is described at\n\n%s\n",
			gnf.GetConfig().APIDoc,
		)

		return c.String(http.StatusOK, info)
	}

}

func pingAPI() func(echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	}
}

func verAPI(gnf gnfinder.GNfinder) func(echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, gnf.GetVersion())
	}
}

func findAPI(gnf gnfinder.GNfinder) func(echo.Context) error {
	return func(c echo.Context) error {
		ctx, cancel := getContext(c)
		defer cancel()
		chErr := make(chan error)

		go func() {
			var err error
			var text, filename string
			var txtExtr float32
			var out output.Output
			var params api.FinderParams
			var format gnfmt.Format
			var opts []config.Option

			err = c.Bind(&params)

			if err == nil {
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

		}()

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
		config.OptWithAllMatches(params.WithAllMatches),
		config.OptWithAmbiguousNames(params.WithAmbiguousNames),
		config.OptWithVerification(
			params.Verification ||
				len(params.Sources) > 0 ||
				params.WithAllMatches,
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
