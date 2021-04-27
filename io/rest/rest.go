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
	"github.com/gnames/gnfmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

			err = c.Bind(&params)

			if err == nil {
				opts := []config.Option{
					config.OptWithBayesOddsDetails(params.OddsDetails),
					config.OptLanguage(getLanguage(params.Language)),
					config.OptWithLanguageDetection(params.LanguageDetection),
					config.OptWithBayes(!params.NoBayes),
					config.OptPreferredSources(params.Sources),
					config.OptWithVerification(
						params.Verification || len(params.Sources) > 0,
					),
					config.OptTokensAround(params.WordsAround),
				}
				gnf = gnf.ChangeConfig(opts...)
				out = gnf.Find(params.Text)
			}

			if err == nil {
				err = c.JSON(http.StatusOK, out)
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
	l, err := lang.NewLanguage(s)
	if err != nil {
		l = lang.DefaultLanguage
	}
	return l
}

func getContext(c echo.Context) (ctx context.Context, cancel func()) {
	ctx = c.Request().Context()
	ctx, cancel = context.WithTimeout(ctx, 1*time.Minute)
	return ctx, cancel
}
