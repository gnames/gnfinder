package web

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gnames/gnfinder/internal/iodoc"
	gnfinder "github.com/gnames/gnfinder/pkg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed static
var static embed.FS

// Run starts GNfinder service for its webiste and RESTful API.
func Run(gnf gnfinder.GNfinder, port int) error {
	slog.Info("Starting the HTTP API server", "port", port)
	e := echo.New()

	var err error
	e.Renderer, err = NewTemplate()
	if err != nil {
		return err
	}

	// doc converts URLs and uploaded files to texts, it is safe for
	// concurrent use, so it is shared by all requests.
	doc := iodoc.New(gnf.GetConfig(), gnfinder.Version)

	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	e.GET("/", home(gnf))
	e.GET("/apidoc", apidoc(gnf))
	e.POST("/find", find(gnf, doc))

	e.GET("/api", infoApiGET(gnf))
	e.GET("/api/", infoApiGET(gnf))
	e.GET("/api/v1", infoApiGET(gnf))
	e.GET("/api/v1/", infoApiGET(gnf))
	e.GET("/api/v1/ping", pingApiGET())
	e.GET("/api/v1/version", verApiGET(gnf))
	e.GET("/api/v1/find/:text", findApiGET(gnf, doc))
	e.POST("/api/v1/find", findApiPOST(gnf, doc))

	fs := http.FileServer(http.FS(static))
	e.GET("/static/*", echo.WrapHandler(fs))

	addr := fmt.Sprintf(":%d", port)
	s := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
	err = e.StartServer(s)
	if err != nil {
		return err
	}

	return nil
}
