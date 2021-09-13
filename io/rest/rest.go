package rest

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/config"
	"github.com/gnames/gnfinder/ent/api"
	"github.com/gnames/gnfinder/ent/lang"
	"github.com/gnames/gnfinder/ent/output"
	"github.com/gnames/gnfinder/ent/verifier"
	"github.com/gnames/gnfmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const withLogs = true

// Run starts HTTP/1 service for scientific names verification.
func Run(gnf gnfinder.GNfinder, port int) {
	log.Printf("Starting the HTTP API server on port %d.", port)
	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())
	if withLogs {
		e.Use(middleware.Logger())
	}

	e.GET("/api/v1/ping", ping())
	e.GET("/api/v1/version", ver(gnf))
	e.POST("/api/v1/find", find(gnf))

	addr := fmt.Sprintf(":%d", port)
	s := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
}

func ping() func(echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	}
}

func ver(gnf gnfinder.GNfinder) func(echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, gnf.GetVersion())
	}
}

func find(gnf gnfinder.GNfinder) func(echo.Context) error {
	return func(c echo.Context) error {
		ctx, cancel := getContext(c)
		defer cancel()
		chErr := make(chan error)

		go func() {
			var err error
			var out output.Output
			var params api.FinderParams
			var format gnfmt.Format

			err = c.Bind(&params)

			if err == nil {
				format, _ = gnfmt.NewFormat(params.Format)
				if format == gnfmt.FormatNone {
					format = gnfmt.CompactJSON
				}
				opts := []config.Option{
					config.OptWithBayesOddsDetails(params.OddsDetails),
					config.OptFormat(format),
					config.OptWithBayes(!params.NoBayes),
					config.OptWithBytesOffset(params.BytesOffset),
					config.OptLanguage(getLanguage(params.Language)),
					config.OptPreferredSources(params.Sources),
					config.OptWithVerification(
						params.Verification || len(params.Sources) > 0,
					),
					config.OptTokensAround(params.WordsAround),
				}

				gnf = gnf.ChangeConfig(opts...)
				out = gnf.Find("", params.Text)
				cfg := gnf.GetConfig()
				if cfg.WithVerification {
					verif := verifier.New(cfg.VerifierURL, cfg.PreferredSources)
					verifiedNames, dur := verif.Verify(out.UniqueNameStrings())
					out.MergeVerification(verifiedNames, dur)
				}
			}

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

func getFormat(s string) gnfmt.Format {
	format, _ := gnfmt.NewFormat(s)
	if format == gnfmt.FormatNone {
		format = gnfmt.CSV
	}
	return format
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
