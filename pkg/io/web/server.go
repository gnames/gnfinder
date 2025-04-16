package web

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	gnfinder "github.com/gnames/gnfinder/pkg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const withLogs = false

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

	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	e.GET("/", home(gnf))
	e.GET("/apidoc", apidoc(gnf))
	e.POST("/find", find(gnf))

	e.GET("/api", infoApiGET(gnf))
	e.GET("/api/", infoApiGET(gnf))
	e.GET("/api/v1", infoApiGET(gnf))
	e.GET("/api/v1/", infoApiGET(gnf))
	e.GET("/api/v1/ping", pingApiGET())
	e.GET("/api/v1/version", verApiGET(gnf))
	e.GET("/api/v1/find/:text", findApiGET(gnf))
	e.POST("/api/v1/find", findApiPOST(gnf))

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
